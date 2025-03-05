package database

import (
	"context"
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

type IDictionary interface {
	CreateWord(ctx context.Context, polish string, translation model.NewTranslation) (bool, error)
	CreateSentence(ctx context.Context, polish string, english string, sentence string) (bool, error)
	CreateTranslation(ctx context.Context, polish string, translation model.NewTranslation) (bool, error)
	DeleteSentence(ctx context.Context, polish string, english string, sentence string) (bool, error)
	DeleteTranslation(ctx context.Context, polish string, english string) (bool, error)
	DeleteWord(ctx context.Context, polish string) (bool, error)
	UpdateTranslation(ctx context.Context, polish string, english string, newEnglish string) (bool, error)
	UpdateWord(ctx context.Context, polish string, newPolish string) (bool, error)
	SelectWord(ctx context.Context, polish string) (*model.Word, error)
}

type DatabaseService struct {
	queryHandler *dictionaryRepository
	DB           *gorm.DB
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

	return &DatabaseService{DB: db}

}

func (r *DatabaseService) CreateWord(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

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
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) CreateSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {
	var word dbmodels.Word

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	err := r.queryHandler.GetWord(tx, polish, &word)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	for i, t := range word.Translations {
		if t.English == english {
			word.Translations[i].Sentences = append(word.Translations[i].Sentences, dbmodels.Sentence{Sentence: sentence})
		}
	}

	if err := r.queryHandler.AddSentence(tx, &word, english, sentence); err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) CreateTranslation(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {
	var word dbmodels.Word
	sentences := make([]dbmodels.Sentence, 0)

	for _, s := range translation.Sentences {
		sentences = append(sentences, dbmodels.Sentence{Sentence: s})
	}

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	err := r.queryHandler.GetWord(tx, polish, &word)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	word.Translations = append(word.Translations, dbmodels.Translation{
		English:   translation.English,
		Sentences: sentences,
	})

	if err = r.queryHandler.AddTranslation(tx, &word, translation.English); err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) DeleteSentence(ctx context.Context, polish string, english string, sentence string) (bool, error) {

	var s dbmodels.Sentence

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	err := tx.Joins("JOIN translations ON sentences.translation_id = translations.id").
		Joins("JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ? AND sentences.sentence = ?", polish, english, sentence).
		First(&s).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, r.CheckWhichDoesntExits(ErrorOptions{Polish: polish, English: english, Sentence: sentence}, r.DB)
		}
		return false, err
	}

	if err := tx.Delete(s).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) DeleteTranslation(ctx context.Context, polish string, english string) (bool, error) {
	var translation dbmodels.Translation
	var count int64

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	err := tx.Joins("RIGHT JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ?", polish, english).
		First(&translation).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, r.CheckWhichDoesntExits(ErrorOptions{Polish: polish, English: english}, r.DB)
		}
		return false, err
	}

	if err := tx.Model(&translation).Delete(&translation).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Model(&dbmodels.Translation{}).Where("word_id = ?", translation.WordID).Count(&count).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if count == 0 {
		if err := tx.Where("ID = ?", translation.WordID).Delete(&dbmodels.Word{}).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) DeleteWord(ctx context.Context, polish string) (bool, error) {
	var word dbmodels.Word

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	if err := tx.Where("polish = ?", polish).First(&word).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, customerrors.WordNotExistsError{Word: polish}
		}
		return false, err
	}

	if err := tx.Delete(&word).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) UpdateWord(ctx context.Context, polish string, newPolish string) (bool, error) {
	var word dbmodels.Word

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	err := tx.Model(&dbmodels.Word{}).Where("polish = ?", polish).First(&word).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, customerrors.WordNotExistsError{Word: polish}
		}
		return false, err
	}

	if err := tx.Model(&word).Update("polish", newPolish).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) UpdateTranslation(ctx context.Context, polish string, english string, newEnglish string) (bool, error) {
	var translation dbmodels.Translation

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	err := tx.Joins("RIGHT JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ?", polish, english).
		First(&translation).Error

	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, r.CheckWhichDoesntExits(ErrorOptions{Polish: polish, English: english}, r.DB)
		}
		return false, err
	}

	err = tx.Model(&translation).Update("english", newEnglish).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) UpdateSentence(ctx context.Context, polish string, english string, sentence string, newSentence string) (bool, error) {
	var s dbmodels.Sentence

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	err := tx.Joins("JOIN translations ON sentences.translation_id = translations.id").
		Joins("JOIN words ON words.id = translations.word_id").
		Where("words.polish = ? AND translations.english = ? AND sentences.sentence = ?", polish, english, sentence).
		First(&s).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, r.CheckWhichDoesntExits(ErrorOptions{Polish: polish, English: english, Sentence: sentence}, r.DB)
		}
		return false, err
	}

	err = tx.Model(&s).Update("sentence", newSentence).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit().Error
}

func (r *DatabaseService) SelectWord(ctx context.Context, polish string) (*model.Word, error) {
	var word dbmodels.Word

	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Preload("Translations.Sentences").Where("polish = ?", polish).First(&word).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerrors.WordNotExistsError{Word: polish}
		}
		return nil, err
	}

	return dbmodels.DBWordToGQLWord(&word), tx.Commit().Error
}

type ErrorOptions struct {
	Polish   string
	English  string
	Sentence string
}

func (r *DatabaseService) CheckWhichDoesntExits(eo ErrorOptions, tx *gorm.DB) error {
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
