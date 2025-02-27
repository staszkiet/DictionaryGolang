package dbmodels

import (
	"github.com/staszkiet/DictionaryGolang/graph/model"
	"gorm.io/gorm"
)

type Word struct {
	gorm.Model
	Polish       string        `json:"polish"`
	Translations []Translation `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}

type Translation struct {
	gorm.Model
	WordID    uint       `json:"wordId"`
	English   string     `json:"english"`
	Sentences []Sentence `gorm:"foreignKey:TranslationID;constraint:OnDelete:CASCADE;"`
}

type Sentence struct {
	gorm.Model
	TranslationID uint   `json:"translationId"`
	Sentence      string `json:"sentence"`
}

func DBSentenceToGQLSentence(s *Sentence) *model.Sentence {
	return &model.Sentence{Sentence: s.Sentence}
}

func DBTranslationToGQLTranslation(t *Translation) *model.Translation {

	sentences := []*model.Sentence{}

	for _, s := range t.Sentences {
		sentences = append(sentences, DBSentenceToGQLSentence(&s))
	}

	return &model.Translation{English: t.English, Sentences: sentences}
}

func DBWordToGQLWord(w *Word) *model.Word {
	translations := []*model.Translation{}

	for _, t := range w.Translations {
		translations = append(translations, DBTranslationToGQLTranslation(&t))
	}

	return &model.Word{Polish: w.Polish, Translations: translations}
}
