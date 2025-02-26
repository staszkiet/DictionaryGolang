package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	dbmodels "github.com/staszkiet/DictionaryGolang/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&dbmodels.Word{}, &dbmodels.Translation{}, &dbmodels.Sentence{})
	if err != nil {
		log.Fatal("Failed to migrate")
	}

	DB = db

	// 	DB.Exec(`
	// 	ALTER TABLE translations
	// 	ADD CONSTRAINT fk_word
	// 	FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE;
	// `)

	//	DB.Exec(`
	//	ALTER TABLE sentences
	//	ADD CONSTRAINT fk_translation
	//	FOREIGN KEY (translation_id) REFERENCES translations(id) ON DELETE CASCADE;
	//
	// `)
}
