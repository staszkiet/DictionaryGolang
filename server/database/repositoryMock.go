package database

import (
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Add(entity interface{}) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockRepository) GetWord(polish string, word *dbmodels.Word) error {
	args := m.Called(word)
	return args.Error(0)
}

func (m *MockRepository) GetSentence(polish string, english string, sentence string, s *dbmodels.Sentence) error {

	args := m.Called(polish, english, sentence, s)
	return args.Error(0)
}

func (m *MockRepository) DeleteSentence(s dbmodels.Sentence) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m *MockRepository) GetTranslation(polish string, english string, translation *dbmodels.Translation) error {

	args := m.Called(polish, english, translation)
	return args.Error(0)
}

func (m *MockRepository) DeleteTranslation(translation *dbmodels.Translation) error {

	args := m.Called(translation)
	return args.Error(0)
}

func (m *MockRepository) DeleteWord(polish string) error {

	args := m.Called(polish)
	return args.Error(0)
}

func (m *MockRepository) Update(sentence interface{}, newSentence string, updateType string) error {

	args := m.Called(sentence, newSentence, updateType)
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
