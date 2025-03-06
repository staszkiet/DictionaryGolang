package database

import (
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Add(tx *gorm.DB, entity interface{}) error {
	args := m.Called(tx, entity)
	return args.Error(0)
}

func (m *MockRepository) GetWord(tx *gorm.DB, polish string, word *dbmodels.Word) error {
	args := m.Called(tx, word)
	return args.Error(0)
}

func (m *MockRepository) GetSentence(tx *gorm.DB, polish string, english string, sentence string, s *dbmodels.Sentence) error {

	args := m.Called(tx, polish, english, sentence, s)
	return args.Error(0)
}

func (m *MockRepository) DeleteSentence(tx *gorm.DB, s dbmodels.Sentence) error {
	args := m.Called(tx, s)
	return args.Error(0)
}

func (m *MockRepository) GetTranslation(tx *gorm.DB, polish string, english string, translation *dbmodels.Translation) error {

	args := m.Called(tx, polish, english, translation)
	return args.Error(0)
}

func (m *MockRepository) DeleteTranslation(tx *gorm.DB, translation *dbmodels.Translation) error {

	args := m.Called(tx, translation)
	return args.Error(0)
}

func (m *MockRepository) DeleteWord(tx *gorm.DB, polish string) error {

	args := m.Called(tx, polish)
	return args.Error(0)
}

func (m *MockRepository) Update(tx *gorm.DB, sentence interface{}, newSentence string, updateType string) error {

	args := m.Called(tx, sentence, newSentence, updateType)
	return args.Error(0)
}

func (m *MockRepository) WithTransaction(fn func(tx *gorm.DB) error) (bool, error) {

	args := m.Called(fn)
	var err error

	mockTx := &gorm.DB{}
	if fn != nil {
		err = fn(mockTx)
	}

	return args.Bool(0), err
}
