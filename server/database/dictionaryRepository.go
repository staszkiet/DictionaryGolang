package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"
	"gorm.io/gorm"
)

type IRepository interface {
	Add(entity interface{}) error
	GetWord(polish string, word *dbmodels.Word) error
	GetSentence(polish string, english string, sentence string, s *dbmodels.Sentence) error
	DeleteSentence(s dbmodels.Sentence) error
	GetTranslation(polish string, english string, translation *dbmodels.Translation) error
	DeleteTranslation(translation *dbmodels.Translation) error
	DeleteWord(polish string) error
	Update(entity interface{}, newEntityString string, updateType string) error
	WithTransaction(fn func(tx *gorm.DB) error) (bool, error)
}

type dictionaryRepository struct {
	db *gorm.DB
}

func (d *dictionaryRepository) GetWord(polish string, word *dbmodels.Word) error {
	err := d.db.Model(&dbmodels.Word{}).Preload("Translations.Sentences").Where("polish = ?", polish).First(word).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.WordNotExistsError{Word: polish}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) Add(entity interface{}) error {

	existsErr := customerrors.GetEntityExistsError(entity)

	if err := d.db.Create(entity).Error; err != nil {
		d.db.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return existsErr
		}
		return err
	}
	return nil

}

func (d *dictionaryRepository) GetSentence(polish string, english string, sentence string, s *dbmodels.Sentence) error {

	err := d.db.Joins("JOIN translations ON sentences.translation_id = translations.id").
		Joins("JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ? AND sentences.sentence = ?", polish, english, sentence).
		First(s).Error
	if err != nil {
		d.db.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.SentenceNotExistsError{Word: polish, Translation: english, Sentence: sentence}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) DeleteSentence(s dbmodels.Sentence) error {
	if err := d.db.Delete(s).Error; err != nil {
		return err
	}
	return nil
}

func (d *dictionaryRepository) GetTranslation(polish string, english string, translation *dbmodels.Translation) error {

	err := d.db.Joins("RIGHT JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ?", polish, english).
		First(translation).Error
	if err != nil {
		d.db.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.TranslationNotExistsError{Word: polish, Translation: english}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) DeleteTranslation(translation *dbmodels.Translation) error {

	var count int64

	if err := d.db.Model(translation).Delete(&translation).Error; err != nil {
		return err
	}

	if err := d.db.Model(&dbmodels.Translation{}).Where("word_id = ?", translation.WordID).Count(&count).Error; err != nil {
		d.db.Rollback()
		return err
	}

	if count == 0 {
		if err := d.db.Where("ID = ?", translation.WordID).Delete(&dbmodels.Word{}).Error; err != nil {
			d.db.Rollback()
			return err
		}
	}
	return nil
}

func (d *dictionaryRepository) DeleteWord(polish string) error {

	if err := d.db.Where("polish = ?", polish).Delete(&dbmodels.Word{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *dictionaryRepository) Update(entity interface{}, newEntity string, updateType string) error {

	existsErr := customerrors.GetUpdatedEntityExistsError(entity, newEntity)

	err := d.db.Model(entity).Update(updateType, newEntity).Error
	if err != nil {
		d.db.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return existsErr
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) WithTransaction(fn func(tx *gorm.DB) error) (bool, error) {

	if err := d.db.Transaction(fn); err != nil {
		return false, err
	}

	return true, nil
}
