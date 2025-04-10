package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"

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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5430 sslmode=%s",
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

// Adds a translation to the dictionary (assumes that polish part wasn't in the dictionary at the time of calling)
func (r *DictionaryService) CreateWord(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {

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
	})
}

// Adds and examplse sentence to the existing translation
func (r *DictionaryService) CreateSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {
	var translation dbmodels.Translation

	return r.repository.WithTransaction(func(txRepo IRepository) error {

		err := txRepo.GetTranslation(polish, english, &translation)
		if err != nil {
			return err
		}

		newSentence := &dbmodels.Sentence{TranslationID: translation.ID, Sentence: sentence}

		if err := txRepo.AddSentence(newSentence); err != nil {
			return err
		}
		return nil
	})

}

// Adds a translation to the dictionary (assumes that the polish part already exists in the dictionary with another translation of it)
func (r *DictionaryService) CreateTranslation(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {

		var word dbmodels.Word
		sentences := make([]dbmodels.Sentence, 0)

		for _, s := range translation.Sentences {
			sentences = append(sentences, dbmodels.Sentence{Sentence: s})
		}
		err := txRepo.GetWord(polish, &word)
		if err != nil {
			return err
		}

		newTranslation := &dbmodels.Translation{
			WordID:    word.ID,
			English:   translation.English,
			Sentences: sentences,
		}

		if err = txRepo.AddTranslation(newTranslation); err != nil {
			return err
		}
		return nil
	})

}

// Deletes an example sentence from given translation
func (r *DictionaryService) DeleteSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var s dbmodels.Sentence
		err := txRepo.GetSentence(polish, english, sentence, &s)
		if err != nil {
			return err
		}

		if err := txRepo.DeleteSentence(s); err != nil {
			return err
		}
		return nil
	})

}

// Deletes an english part of translation
// (If it was the last translation attached to the polish part, the polish part also gets deleted)
func (r *DictionaryService) DeleteTranslation(ctx context.Context, polish string, english string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var translation dbmodels.Translation
		err := txRepo.GetTranslation(polish, english, &translation)
		if err != nil {
			return err
		}

		if err := txRepo.DeleteTranslation(&translation); err != nil {
			return err
		}
		return nil
	})

}

// Deletes whole translation (polish part, english counterparts and its sentences)
func (r *DictionaryService) DeleteWord(ctx context.Context, polish string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		if err := txRepo.DeleteWord(polish); err != nil {
			return err
		}
		return nil
	})
}

// Updates polish part of the translation
func (r *DictionaryService) UpdateWord(ctx context.Context, polish string, newPolish string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var word dbmodels.Word
		err := txRepo.GetWord(polish, &word)
		if err != nil {
			return err
		}

		if err := txRepo.UpdateWord(&word, newPolish); err != nil {
			return err
		}
		return nil
	})

}

// Updates english part of the translation
func (r *DictionaryService) UpdateTranslation(ctx context.Context, polish string, english string, newEnglish string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {
		var translation dbmodels.Translation

		err := txRepo.GetTranslation(polish, english, &translation)

		if err != nil {
			return err
		}

		err = txRepo.UpdateTranslation(&translation, newEnglish)
		if err != nil {
			return err
		}
		return nil
	})

}

// Updates an example sentence of given translation
func (r *DictionaryService) UpdateSentence(ctx context.Context, polish string, english string, sentence string, newSentence string) (bool, error) {

	return r.repository.WithTransaction(func(txRepo IRepository) error {

		var s dbmodels.Sentence
		err := txRepo.GetSentence(polish, english, sentence, &s)

		if err != nil {
			return err
		}

		err = txRepo.UpdateSentence(&s, newSentence)
		if err != nil {
			return err
		}
		return nil
	})

}

// Fetches data regarding given polish word
func (r *DictionaryService) SelectWord(ctx context.Context, polish string) (*model.Word, error) {
	var word dbmodels.Word
	var err error

	_, retErr := r.repository.WithTransaction(func(txRepo IRepository) error {
		if err = txRepo.GetWord(polish, &word); err != nil {
			return err
		}
		return nil
	})
	if retErr == nil {
		return dbmodels.DBWordToGQLWord(&word), nil
	} else {
		return nil, retErr
	}
}
