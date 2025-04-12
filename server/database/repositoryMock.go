package database

import (
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddWord(word *dbmodels.Word) error {
	args := m.Called(word)
	return args.Error(0)
}

func (m *MockRepository) AddTranslation(translation *dbmodels.Translation) error {
	args := m.Called(translation)
	return args.Error(0)
}

func (m *MockRepository) AddSentences(sentences []dbmodels.Sentence) error {
	args := m.Called(sentences)
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

func (m *MockRepository) UpdateWord(word *dbmodels.Word, newPolish string) error {

	args := m.Called(word, newPolish)
	return args.Error(0)
}

func (m *MockRepository) UpdateTranslation(translation *dbmodels.Translation, newTranslation string) error {

	args := m.Called(translation, newTranslation)
	return args.Error(0)
}

func (m *MockRepository) UpdateSentence(sentence *dbmodels.Sentence, newSentence string) error {

	args := m.Called(sentence, newSentence)
	return args.Error(0)
}

func (m *MockRepository) WithTransaction(fn func(repo IRepository) error) (bool, error) {

	args := m.Called(fn)
	var err error

	if fn != nil {
		err = fn(m)
	}

	return args.Bool(0), err
}

func (r *MockRepository) withTx(tx *gorm.DB) IRepository {
	return &MockRepository{}
}
