package database

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"
	"gorm.io/gorm"
)

type IRepository interface {
	CreateWord(tx *gorm.DB, word *dbmodels.Word) error
	GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error
	AddSentence(tx *gorm.DB, polish string, translation string, sentence *dbmodels.Sentence) error
	AddTranslation(tx *gorm.DB, polish string, translation *dbmodels.Translation) error
	GetSentence(tx *gorm.DB, polish string, english string, sentence string, s *dbmodels.Sentence) error
	DeleteSentence(tx *gorm.DB, s dbmodels.Sentence) error
	GetTranslation(tx *gorm.DB, polish string, english string, translation *dbmodels.Translation) error
	DeleteTranslation(tx *gorm.DB, translation *dbmodels.Translation) error
	DeleteWord(tx *gorm.DB, polish string) error
	Update(tx *gorm.DB, sentence interface{}, newSentence string, updateType string) error
	WithTransaction(fn func(tx *gorm.DB) error) (bool, error)
}

type dictionaryRepository struct {
	db *gorm.DB
}

func (d *dictionaryRepository) CreateWord(tx *gorm.DB, word *dbmodels.Word) error {
	if err := tx.Create(word).Error; err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return customerrors.WordExistsError{Word: word.Polish}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error {
	err := tx.Model(&dbmodels.Word{}).Preload("Translations.Sentences").Where("polish = ?", polish).First(word).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.WordNotExistsError{Word: polish}
		}
		return err
	}
	fmt.Println(err)
	return nil
}

func (d *dictionaryRepository) AddSentence(tx *gorm.DB, polish string, translation string, sentence *dbmodels.Sentence) error {
	if err := tx.Create(sentence).Error; err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return customerrors.SentenceExistsError{Word: polish, Translation: translation, Sentence: sentence.Sentence}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) AddTranslation(tx *gorm.DB, polish string, translation *dbmodels.Translation) error {
	if err := tx.Create(translation).Error; err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return customerrors.TranslationExistsError{Word: polish, Translation: translation.English}
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
		tx.Rollback()
		return err
	} else if tx.RowsAffected < 1 {
		return customerrors.CantDeleteSentenceError{Sentence: s.Sentence}
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
		tx.Rollback()
		return err
	} else if tx.RowsAffected < 1 {
		return customerrors.CantDeleteTranslationError{Translation: translation.English}
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
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.WordNotExistsError{Word: polish}
		}
		return err
	} else if tx.RowsAffected < 1 {
		return customerrors.CantDeleteWordError{Word: polish}
	}
	return nil
}

func (d *dictionaryRepository) Update(tx *gorm.DB, sentence interface{}, newSentence string, updateType string) error {

	err := tx.Model(sentence).Update(updateType, newSentence).Error
	if err != nil {
		tx.Rollback()
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
