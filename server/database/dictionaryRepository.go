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
