package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"
	"gorm.io/gorm"
)

type IRepository interface {
	Add(tx *gorm.DB, entity interface{}) error
	GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error
	GetSentence(tx *gorm.DB, polish string, english string, sentence string, s *dbmodels.Sentence) error
	DeleteSentence(tx *gorm.DB, s dbmodels.Sentence) error
	GetTranslation(tx *gorm.DB, polish string, english string, translation *dbmodels.Translation) error
	DeleteTranslation(tx *gorm.DB, translation *dbmodels.Translation) error
	DeleteWord(tx *gorm.DB, polish string) error
	Update(tx *gorm.DB, entity interface{}, newEntityString string, updateType string) error
	WithTransaction(fn func(tx *gorm.DB) error) (bool, error)
}

type dictionaryRepository struct {
	db *gorm.DB
}

func (d *dictionaryRepository) GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error {
	err := tx.Model(&dbmodels.Word{}).Preload("Translations.Sentences").Where("polish = ?", polish).First(word).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.WordNotExistsError{Word: polish}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) Add(tx *gorm.DB, entity interface{}) error {

	existsErr := customerrors.GetEntityExistsError(entity)

	if err := tx.Create(entity).Error; err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return existsErr
		}
		return err
	}
	return nil

}

func (d *dictionaryRepository) GetSentence(tx *gorm.DB, polish string, english string, sentence string, s *dbmodels.Sentence) error {

	err := tx.Joins("JOIN translations ON sentences.translation_id = translations.id").
		Joins("JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ? AND sentences.sentence = ?", polish, english, sentence).
		First(s).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.SentenceNotExistsError{Word: polish, Translation: english, Sentence: sentence}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) DeleteSentence(tx *gorm.DB, s dbmodels.Sentence) error {
	if err := tx.Delete(s).Error; err != nil {
		return err
	}
	return nil
}

func (d *dictionaryRepository) GetTranslation(tx *gorm.DB, polish string, english string, translation *dbmodels.Translation) error {

	err := tx.Joins("RIGHT JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ?", polish, english).
		First(translation).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.TranslationNotExistsError{Word: polish, Translation: english}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) DeleteTranslation(tx *gorm.DB, translation *dbmodels.Translation) error {

	var count int64

	if err := tx.Model(translation).Delete(&translation).Error; err != nil {
		return err
	}

	if err := tx.Model(&dbmodels.Translation{}).Where("word_id = ?", translation.WordID).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}

	if count == 0 {
		if err := tx.Where("ID = ?", translation.WordID).Delete(&dbmodels.Word{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (d *dictionaryRepository) DeleteWord(tx *gorm.DB, polish string) error {

	if err := tx.Where("polish = ?", polish).Delete(&dbmodels.Word{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *dictionaryRepository) Update(tx *gorm.DB, entity interface{}, newEntity string, updateType string) error {

	existsErr := customerrors.GetUpdatedEntityExistsError(entity, newEntity)

	err := tx.Model(entity).Update(updateType, newEntity).Error
	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return existsErr
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) WithTransaction(fn func(tx *gorm.DB) error) (bool, error) {

	tx := d.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}
