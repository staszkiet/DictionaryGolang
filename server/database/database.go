package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	"github.com/staszkiet/DictionaryGolang/server/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseService struct {
	repository IRepository
}

func NewDatabaseService() *DatabaseService {

	var db *gorm.DB
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&dbmodels.Word{}, &dbmodels.Translation{}, &dbmodels.Sentence{})
	if err != nil {
		log.Fatal("Failed to migrate")
	}
	repo := &dictionaryRepository{db: db}
	return &DatabaseService{repository: repo}

}

func NewMockedDatabaseService() *DatabaseService {

	return &DatabaseService{repository: &MockRepository{}}
}

func (r *DatabaseService) CreateWord(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
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

		if err := r.repository.CreateWord(tx, word); err != nil {
			return err
		}
		return nil
	})
}

func (r *DatabaseService) CreateSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {
	var word dbmodels.Word

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
		err := r.repository.GetWord(tx, polish, &word)
		if err != nil {
			return err
		}

		for i, t := range word.Translations {
			if t.English == english {
				word.Translations[i].Sentences = append(word.Translations[i].Sentences, dbmodels.Sentence{Sentence: sentence})
			}
		}

		if err := r.repository.AddSentence(tx, &word, english, sentence); err != nil {
			return err
		}
		return nil
	})

}

func (r *DatabaseService) CreateTranslation(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
		var word dbmodels.Word
		sentences := make([]dbmodels.Sentence, 0)

		for _, s := range translation.Sentences {
			sentences = append(sentences, dbmodels.Sentence{Sentence: s})
		}
		err := r.repository.GetWord(tx, polish, &word)
		if err != nil {
			return err
		}

		word.Translations = append(word.Translations, dbmodels.Translation{
			English:   translation.English,
			Sentences: sentences,
		})

		if err = r.repository.AddTranslation(tx, &word, translation.English); err != nil {
			return err
		}
		return nil
	})

}

func (r *DatabaseService) DeleteSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
		var s dbmodels.Sentence
		err := r.repository.GetSentence(tx, polish, english, sentence, &s)
		if err != nil {
			return err
		}

		if err := r.repository.DeleteSentence(tx, s); err != nil {
			return err
		}
		return nil
	})

}

func (r *DatabaseService) DeleteTranslation(ctx context.Context, polish string, english string) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
		var translation dbmodels.Translation
		err := r.repository.GetTranslation(tx, polish, english, &translation)
		if err != nil {
			return err
		}

		if err := r.repository.DeleteTranslation(tx, &translation); err != nil {
			return err
		}
		return nil
	})

}

func (r *DatabaseService) DeleteWord(ctx context.Context, polish string) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
		if err := r.repository.DeleteWord(tx, polish); err != nil {
			return err
		}
		return nil
	})
}

func (r *DatabaseService) UpdateWord(ctx context.Context, polish string, newPolish string) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
		var word dbmodels.Word
		err := r.repository.GetWord(tx, polish, &word)
		if err != nil {
			return err
		}

		if err := r.repository.Update(tx, &word, newPolish, Word); err != nil {
			return err
		}
		return nil
	})

}

func (r *DatabaseService) UpdateTranslation(ctx context.Context, polish string, english string, newEnglish string) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {
		var translation dbmodels.Translation

		err := r.repository.GetTranslation(tx, polish, english, &translation)

		if err != nil {
			return err
		}

		err = r.repository.Update(tx, &translation, newEnglish, Translation)
		if err != nil {
			return err
		}
		return nil
	})

}

func (r *DatabaseService) UpdateSentence(ctx context.Context, polish string, english string, sentence string, newSentence string) (bool, error) {

	return r.repository.WithTransaction(func(tx *gorm.DB) error {

		var s dbmodels.Sentence
		err := r.repository.GetSentence(tx, polish, english, sentence, &s)

		if err != nil {
			return err
		}

		err = r.repository.Update(tx, &s, newSentence, Sentence)
		if err != nil {
			return err
		}
		return nil
	})

}

func (r *DatabaseService) SelectWord(ctx context.Context, polish string) (*model.Word, error) {
	var word dbmodels.Word

	r.repository.WithTransaction(func(tx *gorm.DB) error {
		if err := r.repository.GetWord(tx, polish, &word); err != nil {
			return err
		}
		return nil
	})
	return dbmodels.DBWordToGQLWord(&word), nil
}
