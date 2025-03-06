package dbmodels

import (
	"github.com/staszkiet/DictionaryGolang/server/graph/model"
)

type Entity interface {
	Word | Translation | Sentence
}

type Word struct {
	ID           uint          `gorm:"primarykey"`
	Polish       string        `json:"polish" gorm:"index;unique"`
	Translations []Translation `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}

type Translation struct {
	ID        uint       `gorm:"primarykey"`
	WordID    uint       `json:"wordId" gorm:"index:,uniqueIndex:translation"`
	English   string     `json:"english" gorm:"index;uniqueIndex:translation"`
	Sentences []Sentence `gorm:"foreignKey:TranslationID;constraint:OnDelete:CASCADE"`
}

type Sentence struct {
	ID            uint   `gorm:"primarykey"`
	TranslationID uint   `json:"translationId" gorm:"uniqueIndex:sentence"`
	Sentence      string `json:"sentence" gorm:"index;uniqueIndex:sentence"`
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
