package dbmodels

import "gorm.io/gorm"

type Word struct {
	gorm.Model
	Polish       string        `json:"polish"`
	Translations []Translation `gorm:"foreignKey:WordID"`
}

type Translation struct {
	gorm.Model
	WordID    uint       `json:"wordId"`
	English   string     `json:"english"`
	Sentences []Sentence `gorm:"foreignKey:TranslationID"`
}

type Sentence struct {
	gorm.Model
	TranslationID uint   `json:"translationId"`
	Sentence      string `json:"sentence"`
}
