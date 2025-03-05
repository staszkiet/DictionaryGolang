package database

import (
	"context"
	"testing"

	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	"github.com/staszkiet/DictionaryGolang/server/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWord_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "dom"
	translation := model.NewTranslation{
		English:   "house",
		Sentences: []string{"This is my house", "I bought a new house"},
	}

	expectedWord := &dbmodels.Word{
		Polish: polish,
		Translations: []dbmodels.Translation{
			{
				English: translation.English,
				Sentences: []dbmodels.Sentence{
					{Sentence: "This is my house"},
					{Sentence: "I bought a new house"},
				},
			},
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)

	mockRepo.On("CreateWord", mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {

			wordArg := args.Get(1).(*dbmodels.Word)
			assert.Equal(t, expectedWord.Polish, wordArg.Polish)
			assert.Equal(t, expectedWord.Translations[0].English, wordArg.Translations[0].English)
			assert.ElementsMatch(t, expectedWord.Translations[0].Sentences, wordArg.Translations[0].Sentences)
		})

	success, err := dbService.CreateWord(context.Background(), polish, translation)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestCreateTranslation_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "pisać"
	translation := model.NewTranslation{
		English:   "write",
		Sentences: []string{"I write poems", "I write books"},
	}

	dbWord := &dbmodels.Word{
		Polish: "pisać",
		Translations: []dbmodels.Translation{
			{
				English: "type",
				Sentences: []dbmodels.Sentence{
					{Sentence: "I often type"},
					{Sentence: "What should I type?"},
				},
			},
		},
	}

	expectedWord := &dbmodels.Word{
		Polish: "pisać",
		Translations: []dbmodels.Translation{
			{
				English: "type",
				Sentences: []dbmodels.Sentence{
					{Sentence: "I often type"},
					{Sentence: "What should I type?"},
				},
			},
			{
				English: "write",
				Sentences: []dbmodels.Sentence{
					{Sentence: "I write poems"},
					{Sentence: "I write books"},
				},
			},
		},
	}
	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)
	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(1).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	mockRepo.On("AddTranslation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(1).(*dbmodels.Word)
			assert.Equal(t, expectedWord.Polish, wordArg.Polish)
			assert.ElementsMatch(t, expectedWord.Translations, wordArg.Translations)
		})

	success, err := dbService.CreateTranslation(context.Background(), polish, translation)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestCreateSentence_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "książka"
	English := "book"
	sentence := "I have never read a book"

	dbWord := &dbmodels.Word{
		Polish: "książka",
		Translations: []dbmodels.Translation{
			{
				English: "book",
				Sentences: []dbmodels.Sentence{
					{Sentence: "I read a good book"},
				},
			},
		},
	}

	expectedWord := &dbmodels.Word{
		Polish: "książka",
		Translations: []dbmodels.Translation{
			{
				English: "book",
				Sentences: []dbmodels.Sentence{
					{Sentence: "I read a good book"},
					{Sentence: "I have never read a book"},
				},
			},
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)
	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(1).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	mockRepo.On("AddSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(1).(*dbmodels.Word)
			assert.Equal(t, expectedWord.Polish, wordArg.Polish)
			assert.ElementsMatch(t, expectedWord.Translations, wordArg.Translations)
		})

	success, err := dbService.CreateSentence(context.Background(), polish, English, sentence)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteSentence_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "książka"
	English := "book"
	sentence := "I have never read a book"

	dbSentence := dbmodels.Sentence{Sentence: "I have never read a book"}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)
	mockRepo.On("GetSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(4).(*dbmodels.Sentence)
		*(wordArg) = dbSentence
	})

	mockRepo.On("DeleteSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(1).(dbmodels.Sentence)
			assert.Equal(t, wordArg, dbSentence)
		})

	success, err := dbService.DeleteSentence(context.Background(), polish, English, sentence)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteTranslation_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "książka"
	English := "book"

	dbTranslation := &dbmodels.Translation{
		English: "book",
		Sentences: []dbmodels.Sentence{
			{
				Sentence: "I have never read a book",
			},
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)
	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(3).(*dbmodels.Translation)
		*(wordArg) = *(dbTranslation)
	})

	mockRepo.On("DeleteTranslation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(1).(*dbmodels.Translation)
			assert.Equal(t, wordArg.English, dbTranslation.English)
			assert.ElementsMatch(t, wordArg.Sentences, dbTranslation.Sentences)
		})

	success, err := dbService.DeleteTranslation(context.Background(), polish, English)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteWord_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "książka"

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)

	mockRepo.On("DeleteWord", mock.Anything, mock.Anything).Return(nil)

	success, err := dbService.DeleteWord(context.Background(), polish)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestUpdateWord_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "dm"
	newPolish := "dom"

	dbWord := &dbmodels.Word{
		Polish: polish,
		Translations: []dbmodels.Translation{
			{
				English: "house",
				Sentences: []dbmodels.Sentence{
					{Sentence: "This is my house"},
					{Sentence: "I bought a new house"},
				},
			},
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)

	mockRepo.On("GetWord", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(1).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.Anything, Word).Return(nil)

	success, err := dbService.UpdateWord(context.Background(), polish, newPolish)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestUpdateSentence_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "książka"
	English := "bok"
	sentence := "I have never red a book"
	newSentence := "I have never read a book"

	dbSentence := &dbmodels.Sentence{
		Sentence: "I have never read a book",
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)

	mockRepo.On("GetSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(4).(*dbmodels.Sentence)
		*(wordArg) = *(dbSentence)
	})

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.Anything, Sentence).Return(nil)

	success, err := dbService.UpdateSentence(context.Background(), polish, English, sentence, newSentence)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTranslation_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "książka"
	English := "bok"
	newEnglish := "book"

	dbTranslation := &dbmodels.Translation{
		English: "bok",
		Sentences: []dbmodels.Sentence{
			{
				Sentence: "I have never read a book",
			},
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)

	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(3).(*dbmodels.Translation)
		*(wordArg) = *(dbTranslation)
	})

	mockRepo.On("Update", mock.Anything, mock.Anything, mock.Anything, Translation).Return(nil)

	success, err := dbService.UpdateTranslation(context.Background(), polish, English, newEnglish)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestSelectWord_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	dbService := &DatabaseService{repository: mockRepo}

	polish := "dom"

	dbWord := &dbmodels.Word{
		Polish: "dom",
		Translations: []dbmodels.Translation{
			{
				English: "house",
				Sentences: []dbmodels.Sentence{
					{Sentence: "This is my house"},
					{Sentence: "I bought a new house"},
				},
			},
		},
	}

	expectedWord := &model.Word{
		Polish: "dom",
		Translations: []*model.Translation{
			{
				English: "house",
				Sentences: []*model.Sentence{
					{Sentence: "This is my house"},
					{Sentence: "I bought a new house"},
				},
			},
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)

	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(1).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	retWord, err := dbService.SelectWord(context.Background(), polish)

	assert.NoError(t, err)
	assert.Equal(t, retWord, expectedWord)

	mockRepo.AssertExpectations(t)
}
