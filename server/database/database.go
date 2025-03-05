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
	queryHandler *dictionaryRepository
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
	qh := &dictionaryRepository{DB: db}
	return &DatabaseService{queryHandler: qh}

}

func (r *DatabaseService) CreateWord(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
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

		if err := r.queryHandler.CreateWord(tx, word); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

func (r *DatabaseService) CreateSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {
	var word dbmodels.Word

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		err := r.queryHandler.GetWord(tx, polish, &word)
		if err != nil {
			tx.Rollback()
			return err
		}

		for i, t := range word.Translations {
			if t.English == english {
				word.Translations[i].Sentences = append(word.Translations[i].Sentences, dbmodels.Sentence{Sentence: sentence})
			}
		}

		if err := r.queryHandler.AddSentence(tx, &word, english, sentence); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

}

func (r *DatabaseService) CreateTranslation(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		var word dbmodels.Word
		sentences := make([]dbmodels.Sentence, 0)

		for _, s := range translation.Sentences {
			sentences = append(sentences, dbmodels.Sentence{Sentence: s})
		}
		err := r.queryHandler.GetWord(tx, polish, &word)
		if err != nil {
			tx.Rollback()
			return err
		}

		word.Translations = append(word.Translations, dbmodels.Translation{
			English:   translation.English,
			Sentences: sentences,
		})

		if err = r.queryHandler.AddTranslation(tx, &word, translation.English); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

}

func (r *DatabaseService) DeleteSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		var s dbmodels.Sentence
		err := r.queryHandler.GetSentence(tx, polish, english, sentence, &s)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := r.queryHandler.DeleteSentence(tx, s); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

}

func (r *DatabaseService) DeleteTranslation(ctx context.Context, polish string, english string) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		var translation dbmodels.Translation
		err := r.queryHandler.GetTranslation(tx, polish, english, &translation)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := r.queryHandler.DeleteTranslation(tx, &translation); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

}

func (r *DatabaseService) DeleteWord(ctx context.Context, polish string) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		if err := r.queryHandler.DeleteWord(tx, polish); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

func (r *DatabaseService) UpdateWord(ctx context.Context, polish string, newPolish string) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		var word dbmodels.Word
		err := r.queryHandler.GetWord(tx, polish, &word)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := r.queryHandler.UpdateWord(tx, &word, newPolish); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

}

func (r *DatabaseService) UpdateTranslation(ctx context.Context, polish string, english string, newEnglish string) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		var translation dbmodels.Translation

		err := r.queryHandler.GetTranslation(tx, polish, english, &translation)

		if err != nil {
			tx.Rollback()
			return err
		}

		err = r.queryHandler.UpdateTranslation(tx, &translation, newEnglish)
		if err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

}

func (r *DatabaseService) UpdateSentence(ctx context.Context, polish string, english string, sentence string, newSentence string) (bool, error) {

	return r.queryHandler.WithTransaction(func(tx *gorm.DB) error {

		var s dbmodels.Sentence
		err := r.queryHandler.GetSentence(tx, polish, english, sentence, &s)

		if err != nil {
			tx.Rollback()
			return err
		}

		err = r.queryHandler.UpdateSentence(tx, &s, newSentence)
		if err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

}

func (r *DatabaseService) SelectWord(ctx context.Context, polish string) (*model.Word, error) {
	var word dbmodels.Word

	r.queryHandler.WithTransaction(func(tx *gorm.DB) error {
		if err := r.queryHandler.GetWord(tx, polish, &word); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	return dbmodels.DBWordToGQLWord(&word), nil
}
