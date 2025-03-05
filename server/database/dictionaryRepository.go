package database

import (
	"errors"
	"fmt"

	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"
	"gorm.io/gorm"
)

type IRepository interface {
	CreateWord(tx *gorm.DB, word *dbmodels.Word) error
	GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error
	AddSentence(tx *gorm.DB, word *dbmodels.Word, translation string, sentence string) error
	AddTranslation(tx *gorm.DB, word *dbmodels.Word, translation string) error
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customerrors.WordExistsError{Word: word.Polish}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error {
	err := tx.Model(&dbmodels.Word{}).Preload("Translations.Sentences").Where("polish = ?", polish).First(word).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *dictionaryRepository) AddSentence(tx *gorm.DB, word *dbmodels.Word, translation string, sentence string) error {
	if err := tx.Save(word).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customerrors.SentenceExistsError{Word: word.Polish, Translation: translation, Sentence: sentence}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) AddTranslation(tx *gorm.DB, word *dbmodels.Word, translation string) error {
	if err := tx.Save(word).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customerrors.TranslationExistsError{Word: word.Polish, Translation: translation}
		}
		return err
	}
	return nil
}

// func (d *dictionaryRepository) CheckWhichDoesntExits(eo ErrorOptions, tx *gorm.DB) error {
// 	var word dbmodels.Word
// 	var translation dbmodels.Translation
// 	var sentence dbmodels.Sentence
// 	err := tx.Model(&dbmodels.Word{}).Where("polish = ?", eo.Polish).First(&word).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return customerrors.WordNotExistsError{Word: eo.Polish}
// 		}
// 		return err
// 	}
// 	err = tx.Model(&dbmodels.Translation{}).Where("english = ? AND word_id = ?", eo.English, word.ID).First(&translation).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return customerrors.TranslationNotExistsError{Word: eo.Polish, Translation: eo.English}
// 		}
// 		return err
// 	}
// 	if eo.Sentence != "" {
// 		err = tx.Model(&dbmodels.Sentence{}).Where("sentence = ? AND translation_id = ?", eo.Sentence, translation.ID).First(&sentence).Error
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return customerrors.SentenceNotExistsError{Word: eo.Polish, Translation: eo.English, Sentence: eo.Sentence}
// 			}
// 			return err
// 		}
// 	}
// 	return fmt.Errorf("nieznany błąd")
// }

func (d *dictionaryRepository) GetSentence(tx *gorm.DB, polish string, english string, sentence string, s *dbmodels.Sentence) error {

	err := tx.Joins("JOIN translations ON sentences.translation_id = translations.id").
		Joins("JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ? AND sentences.sentence = ?", polish, english, sentence).
		First(s).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (d *dictionaryRepository) DeleteSentence(tx *gorm.DB, s dbmodels.Sentence) error {
	if err := tx.Delete(s).Error; err != nil {
		tx.Rollback()
		return err
	} else if tx.RowsAffected < 1 {
		return fmt.Errorf("rekord nie może zostać usunięty, gdyż nie ma go w słowniku")
	}
	return nil
}

func (d *dictionaryRepository) GetTranslation(tx *gorm.DB, polish string, english string, translation *dbmodels.Translation) error {

	err := tx.Joins("RIGHT JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ?", polish, english).
		First(translation).Error
	if err != nil {
		tx.Rollback()
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
		return fmt.Errorf("rekord nie może zostać usunięty, gdyż nie ma go w słowniku")
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
		return fmt.Errorf("rekord nie może zostać usunięty, gdyż nie ma go w słowniku")
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

type ErrorOptions struct {
	Polish   string
	English  string
	Sentence string
}

const (
	Word        = "polish"
	Translation = "english"
	Sentence    = "sentence"
)
