package database

import (
	"testing"

	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"

	"github.com/staszkiet/DictionaryGolang/server/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWordOrAddTranslationOrSentence_WhenDataIsValid_ShouldReturnSuccess(t *testing.T) {

	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

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

	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(customerrors.WordNotExistsError{Word: polish})

	mockRepo.On("WithTransaction", mock.Anything).Return(true)

	mockRepo.On("AddWord", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(0).(*dbmodels.Word)
			assert.Equal(t, expectedWord.Polish, wordArg.Polish)
			assert.Equal(t, expectedWord.Translations[0].English, wordArg.Translations[0].English)
			assert.ElementsMatch(t, expectedWord.Translations[0].Sentences, wordArg.Translations[0].Sentences)
		})
	success, err := dbService.CreateWordOrAddTranslationOrSentence(polish, translation)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestCreateWordOrAddTranslationOrSentence_WordExistsButTranslationDoesnt_ShouldReturnSuccess(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "pisać"
	translation := model.NewTranslation{
		English:   "write",
		Sentences: []string{"I write poems", "I write books"},
	}

	dbWord := &dbmodels.Word{
		ID:     2,
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

	expectedTranslation := &dbmodels.Translation{
		WordID:  2,
		English: "write",
		Sentences: []dbmodels.Sentence{
			{Sentence: "I write poems"},
			{Sentence: "I write books"},
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true)
	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(0).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything).Return(customerrors.TranslationNotExistsError{Word: polish, Translation: translation.English})

	mockRepo.On("AddTranslation", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(0).(*dbmodels.Translation)
			assert.Equal(t, expectedTranslation, wordArg)
		})

	success, err := dbService.CreateWordOrAddTranslationOrSentence(polish, translation)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestCreateWordOrAddTranslationOrSentence_WordAndTranslationAlreadyExistsButSomeSentencesDont_ShouldReturnSuccess(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "book"
	sentence := "I have never read a book"

	dbTranslation := &dbmodels.Translation{
		ID:      2,
		English: "book",
		Sentences: []dbmodels.Sentence{
			{Sentence: "I read a good book"},
		},
	}

	dbWord := &dbmodels.Word{
		ID:     2,
		Polish: "pisać",
		Translations: []dbmodels.Translation{
			*dbTranslation,
		},
	}

	expectedSentence := []dbmodels.Sentence{dbmodels.Sentence{
		TranslationID: 2,
		Sentence:      "I have never read a book",
	}}

	mockRepo.On("WithTransaction", mock.Anything).Return(true)

	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(0).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})
	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(2).(*dbmodels.Translation)
		*(wordArg) = *(dbTranslation)
	})

	mockRepo.On("AddSentences", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(0).([]dbmodels.Sentence)
			assert.Equal(t, expectedSentence, wordArg)
		})

	success, err := dbService.CreateWordOrAddTranslationOrSentence(polish,
		model.NewTranslation{English: English, Sentences: []string{sentence}})

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestCreateWord_WordTranslationAndAllSentencesAlreadyExist(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "book"
	sentence := "I read a good book"

	dbTranslation := &dbmodels.Translation{
		ID:      2,
		English: "book",
		Sentences: []dbmodels.Sentence{
			{Sentence: "I read a good book"},
		},
	}

	dbWord := &dbmodels.Word{
		ID:     2,
		Polish: "pisać",
		Translations: []dbmodels.Translation{
			*dbTranslation,
		},
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)
	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(0).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})
	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(2).(*dbmodels.Translation)
		*(wordArg) = *(dbTranslation)
	})

	success, err := dbService.CreateWordOrAddTranslationOrSentence(polish,
		model.NewTranslation{English: English, Sentences: []string{sentence}})

	assert.Nil(t, err)
	assert.False(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteSentence_SentenceExists_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "book"
	sentence := "I have never read a book"

	dbSentence := dbmodels.Sentence{Sentence: "I have never read a book"}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)
	mockRepo.On("GetSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(3).(*dbmodels.Sentence)
		*(wordArg) = dbSentence
	})

	mockRepo.On("DeleteSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(0).(dbmodels.Sentence)
			assert.Equal(t, wordArg, dbSentence)
		})

	success, err := dbService.DeleteSentence(polish, English, sentence)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteSentence_SentenceDoesntExist_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "book"
	sentence := "I have never read a book"

	expectedError := customerrors.SentenceNotExistsError{Word: polish, Translation: English, Sentence: sentence}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)
	mockRepo.On("GetSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.DeleteSentence(polish, English, sentence)

	assert.Error(t, err)
	assert.False(t, success)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}
func TestDeleteTranslation_TranslationExists_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

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
	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(2).(*dbmodels.Translation)
		*(wordArg) = *(dbTranslation)
	})

	mockRepo.On("DeleteTranslation", mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			wordArg := args.Get(0).(*dbmodels.Translation)
			assert.Equal(t, wordArg.English, dbTranslation.English)
			assert.ElementsMatch(t, wordArg.Sentences, dbTranslation.Sentences)
		})

	success, err := dbService.DeleteTranslation(polish, English)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteTranslation_TranslationOrWordDoesntExist_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "book"

	expectedError := customerrors.TranslationNotExistsError{Word: polish, Translation: English}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)
	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.DeleteTranslation(polish, English)

	assert.Nil(t, err)
	assert.False(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteWord_WordExists_Success(t *testing.T) {

	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"

	mockRepo.On("DeleteWord", mock.Anything, mock.Anything).Return(nil)

	success, err := dbService.DeleteWord(polish)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestDeleteWord_WordDoesntExist_Success(t *testing.T) {

	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"

	mockRepo.On("DeleteWord", mock.Anything, mock.Anything).Return(nil)

	success, err := dbService.DeleteWord(polish)

	assert.Nil(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestUpdateWord_WordToUpdateExists_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

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

	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(0).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	mockRepo.On("UpdateWord", mock.Anything, mock.Anything).Return(nil)

	success, err := dbService.UpdateWord(polish, newPolish)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestUpdateWord_WordToUpdateDoesntExist_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "dm"
	newPolish := "dom"
	expectedError := customerrors.WordNotExistsError{Word: polish}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)

	mockRepo.On("GetWord", mock.Anything, mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.UpdateWord(polish, newPolish)

	assert.Error(t, err)
	assert.False(t, success)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateWord_UpdatedWordExists_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

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

	expectedError := customerrors.WordExistsError{Word: polish}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)

	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(0).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	mockRepo.On("UpdateWord", mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.UpdateWord(polish, newPolish)

	assert.Error(t, err)
	assert.False(t, success)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}
func TestUpdateSentence_SentenceToUpdateExists_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "bok"
	sentence := "I have never red a book"
	newSentence := "I have never read a book"

	dbSentence := &dbmodels.Sentence{
		Sentence: "I have never read a book",
	}

	mockRepo.On("WithTransaction", mock.Anything).Return(true, nil)

	mockRepo.On("GetSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(3).(*dbmodels.Sentence)
		*(wordArg) = *(dbSentence)
	})

	mockRepo.On("UpdateSentence", mock.Anything, mock.Anything).Return(nil)

	success, err := dbService.UpdateSentence(polish, English, sentence, newSentence)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestUpdateSentence_SentenceToUpdateDoesntExist_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "bok"
	sentence := "I have never red a book"
	newSentence := "I have never read a book"
	expectedError := customerrors.SentenceNotExistsError{Word: polish, Translation: English, Sentence: sentence}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)

	mockRepo.On("GetSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.UpdateSentence(polish, English, sentence, newSentence)

	assert.Error(t, err)
	assert.False(t, success)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateSentence_UpdatedSentenceExists_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "bok"
	sentence := "I have never red a book"
	newSentence := "I have never read a book"

	dbSentence := &dbmodels.Sentence{
		Sentence: "I have never read a book",
	}

	expectedError := customerrors.SentenceExistsError{Sentence: newSentence}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)

	mockRepo.On("GetSentence", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(3).(*dbmodels.Sentence)
		*(wordArg) = *(dbSentence)
	})

	mockRepo.On("UpdateSentence", mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.UpdateSentence(polish, English, sentence, newSentence)

	assert.Error(t, err)
	assert.False(t, success)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTranslation_TranslationToUpdateExists_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

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

	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(2).(*dbmodels.Translation)
		*(wordArg) = *(dbTranslation)
	})

	mockRepo.On("UpdateTranslation", mock.Anything, mock.Anything).Return(nil)

	success, err := dbService.UpdateTranslation(polish, English, newEnglish)

	assert.NoError(t, err)
	assert.True(t, success)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTranslation_TranslationToUpdateDoesntExist_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "książka"
	English := "bok"
	newEnglish := "book"

	expectedError := customerrors.TranslationNotExistsError{Word: polish, Translation: English}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)

	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.UpdateTranslation(polish, English, newEnglish)

	assert.Error(t, err)
	assert.False(t, success)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTranslation_UpdatedTranslationExists_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

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

	expectedError := customerrors.SentenceExistsError{Sentence: English}

	mockRepo.On("WithTransaction", mock.Anything).Return(false)

	mockRepo.On("GetTranslation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(2).(*dbmodels.Translation)
		*(wordArg) = *(dbTranslation)
	})

	mockRepo.On("UpdateTranslation", mock.Anything, mock.Anything).Return(expectedError)

	success, err := dbService.UpdateTranslation(polish, English, newEnglish)

	assert.Error(t, err)
	assert.False(t, success)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestSelectWord_WordExists_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

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

	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wordArg := args.Get(0).(*dbmodels.Word)
		*(wordArg) = *(dbWord)
	})

	retWord, err := dbService.SelectWord(polish)

	assert.NoError(t, err)
	assert.Equal(t, retWord, expectedWord)

	mockRepo.AssertExpectations(t)
}

func TestSelectWord_WordDoesntExist_ShouldReturnError(t *testing.T) {
	mockRepo := new(MockRepository)
	dbService := &DictionaryService{repository: mockRepo}

	polish := "dom"

	expectedError := customerrors.WordNotExistsError{Word: polish}

	mockRepo.On("GetWord", mock.Anything, mock.Anything).Return(expectedError)

	retWord, err := dbService.SelectWord(polish)

	assert.Error(t, err)
	assert.Nil(t, retWord)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}
