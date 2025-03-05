package database

import (
	"context"
	"errors"

	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"
	"github.com/staszkiet/DictionaryGolang/server/graph/model"
	"gorm.io/gorm"
)

type IDatabase interface {
	CreateWord(tx *gorm.DB, word *dbmodels.Word) (bool, error)
	CreateSentence(ctx context.Context, polish string, english string, sentence string) (bool, error)
	CreateTranslation(ctx context.Context, polish string, translation model.NewTranslation) (bool, error)
	DeleteSentence(ctx context.Context, polish string, english string, sentence string) (bool, error)
	DeleteTranslation(ctx context.Context, polish string, english string) (bool, error)
	DeleteWord(ctx context.Context, polish string) (bool, error)
	UpdateTranslation(ctx context.Context, polish string, english string, newEnglish string) (bool, error)
	UpdateWord(ctx context.Context, polish string, newPolish string) (bool, error)
	SelectWord(ctx context.Context, polish string) (*model.Word, error)
}

type dictionaryRepository struct {
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

// TODO no record error handling
func (d *dictionaryRepository) GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error {
	err := tx.Model(&dbmodels.Word{}).Preload("Translations.Sentences").Where("polish = ?", polish).First(word).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *dictionaryRepository) AddSentence(tx *gorm.DB, word *dbmodels.Word, translation string, sentence string) error {
	if err := tx.Save(word).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customerrors.SentenceExistsError{Word: word.Polish, Translation: translation, Sentence: sentence}
		}
		return err
	}
	return nil
}

func (d *dictionaryRepository) AddTranslation(tx *gorm.DB, word *dbmodels.Word, translation string) error {
	if err := tx.Save(word).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customerrors.TranslationExistsError{Word: word.Polish, Translation: translation}
		}
		return err
	}
	return nil
}
