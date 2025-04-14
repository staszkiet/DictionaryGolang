package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"

	"github.com/staszkiet/DictionaryGolang/server/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DictionaryService struct {
	repository IRepository
}

// Creates new database service to handle operations on repository
func NewDatabaseService() *DictionaryService {

	var db *gorm.DB
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")
	sslmode := os.Getenv("POSTGRES_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s",
		host, user, password, dbname, sslmode)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&dbmodels.Word{}, &dbmodels.Translation{}, &dbmodels.Sentence{})
	if err != nil {
		log.Fatal("Failed to migrate")
	}
	repo := &dictionaryRepository{db: db}
	return &DictionaryService{repository: repo}

}

// Adds a translation to the dictionary (whether given polish word exists or not). I translation to given word already exists then adds
// sum of sentences to the translation
func (r *DictionaryService) CreateWordOrAddTranslationOrSentence(polish string, translation model.NewTranslation) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {

		var dbword dbmodels.Word
		var dbtranslation dbmodels.Translation

		var err error
		if err = txRepo.GetWord(polish, &dbword); err == nil {
			if err = txRepo.GetTranslation(polish, translation.English, &dbtranslation); err == nil {
				existingSentencesMap := make(map[string]bool)
				newSentences := make([]dbmodels.Sentence, 0)

				for _, s := range dbtranslation.Sentences {
					existingSentencesMap[s.Sentence] = true
				}

				for _, s := range translation.Sentences {
					if !existingSentencesMap[s] {
						newSentences = append(newSentences, dbmodels.Sentence{Sentence: s, TranslationID: dbtranslation.ID})
					}
				}

				if len(newSentences) > 0 {
					txRepo.AddSentences(newSentences)
				}

				return nil
			} else {
				sentences := make([]dbmodels.Sentence, 0)

				for _, s := range translation.Sentences {
					sentences = append(sentences, dbmodels.Sentence{Sentence: s})
				}
				newTranslation := &dbmodels.Translation{
					WordID:    dbword.ID,
					English:   translation.English,
					Sentences: sentences,
				}

				if err = txRepo.AddTranslation(newTranslation); err != nil {
					return err
				}

				return nil
			}
		}

		if errors.Is(err, customerrors.WordNotExistsError{Word: polish}) {
			sentences := make([]dbmodels.Sentence, 0)

			for _, s := range translation.Sentences {
				sentences = append(sentences, dbmodels.Sentence{Sentence: s})
			}

			var convertedTranslations []dbmodels.Translation

			convertedTranslations = append(convertedTranslations, dbmodels.Translation{
				English:   translation.English,
				Sentences: sentences,
			})

			word := &dbmodels.Word{
				Polish:       polish,
				Translations: convertedTranslations,
			}

			if err := txRepo.AddWord(word); err != nil {
				return err
			}
			return nil
		}
		return err

	}, true, true)
}

// Deletes an example sentence from given translation
func (r *DictionaryService) DeleteSentence(polish string, english string, sentence string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var s dbmodels.Sentence
		err := txRepo.GetSentence(polish, english, sentence, &s)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		if err := txRepo.DeleteSentence(s); err != nil {
			return err
		}
		return nil
	}, false, false)

}

// Deletes an english part of translation
// (If it was the last translation attached to the polish part, the polish part also gets deleted)
func (r *DictionaryService) DeleteTranslation(polish string, english string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var translation dbmodels.Translation
		err := txRepo.GetTranslation(polish, english, &translation)
		if err != nil {
			if errors.Is(err, customerrors.TranslationNotExistsError{Word: polish, Translation: english}) {
				return nil
			}
			return err
		}

		if err := txRepo.DeleteTranslation(&translation); err != nil {
			return err
		}
		return nil
	}, false, false)

}

// Deletes whole translation (polish part, english counterparts and its sentences)
func (r *DictionaryService) DeleteWord(polish string) (bool, error) {

	if err := r.repository.DeleteWord(polish); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	return true, nil
}

// Updates polish part of the translation
func (r *DictionaryService) UpdateWord(polish string, newPolish string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var word dbmodels.Word
		err := txRepo.GetWord(polish, &word)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return customerrors.WordNotExistsError{Word: polish}
			}
			return err
		}

		if err := txRepo.UpdateWord(&word, newPolish); err != nil {
			return err
		}
		return nil
	}, false, false)

}

// Updates english part of the translation
func (r *DictionaryService) UpdateTranslation(polish string, english string, newEnglish string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var translation dbmodels.Translation

		err := txRepo.GetTranslation(polish, english, &translation)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return customerrors.TranslationNotExistsError{Word: polish, Translation: english}
			}
			return err
		}

		err = txRepo.UpdateTranslation(&translation, newEnglish)
		if err != nil {
			return err
		}
		return nil
	}, false, false)

}

// Updates an example sentence of given translation
func (r *DictionaryService) UpdateSentence(polish string, english string, sentence string, newSentence string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {

		var s dbmodels.Sentence
		err := txRepo.GetSentence(polish, english, sentence, &s)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return customerrors.SentenceNotExistsError{Word: polish, Translation: english, Sentence: sentence}
			}
			return err
		}

		err = txRepo.UpdateSentence(&s, newSentence)
		if err != nil {
			return err
		}
		return nil
	}, false, false)

}

// Fetches data regarding given polish word
func (r *DictionaryService) SelectWord(polish string) (*model.Word, error) {
	var word dbmodels.Word
	var err error

	if err = r.repository.GetWord(polish, &word); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerrors.WordNotExistsError{Word: polish}
		}
		return nil, err
	}

	return dbmodels.DBWordToGQLWord(&word), nil
}
