package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/server/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {

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

	return db

}

type ErrorOptions struct {
	Polish   string
	English  string
	Sentence string
}

func CheckWhichDoesntExits(eo ErrorOptions, tx *gorm.DB) error {
	var word dbmodels.Word
	var translation dbmodels.Translation
	var sentence dbmodels.Sentence
	err := tx.Model(&dbmodels.Word{}).Where("polish = ?", eo.Polish).First(&word).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.WordNotExistsError{Word: eo.Polish}
		}
		return err
	}
	err = tx.Model(&dbmodels.Translation{}).Where("english = ? AND word_id = ?", eo.English, word.ID).First(&translation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerrors.TranslationNotExistsError{Word: eo.Polish, Translation: eo.English}
		}
		return err
	}
	if eo.Sentence != "" {
		err = tx.Model(&dbmodels.Sentence{}).Where("sentence = ? AND translation_id = ?", eo.Sentence, translation.ID).First(&sentence).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return customerrors.SentenceNotExistsError{Word: eo.Polish, Translation: eo.English, Sentence: eo.Sentence}
			}
			return err
		}
	}
	return fmt.Errorf("nieznany błąd")
}
