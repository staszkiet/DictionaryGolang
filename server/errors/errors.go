package customerrors

import (
	"fmt"

	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
)

//errors for adding values to DB

type WordExistsError struct {
	Word string
}

func (e WordExistsError) Error() string {
	return fmt.Sprintf("słowo %s znajduje się już w słowniku", e.Word)
}

type SentenceExistsError struct {
	Sentence string
}

func (e SentenceExistsError) Error() string {
	return fmt.Sprintf("zdanie %s już jest dodane do tłumaczenia danego tłumaczenia", e.Sentence)
}

type TranslationExistsError struct {
	Translation string
}

func (e TranslationExistsError) Error() string {
	return fmt.Sprintf("tłumaczenie %s już jest dodane do danego słowa", e.Translation)
}

func GetEntityExistsError(entity interface{}) error {
	switch entity := entity.(type) {
	case *dbmodels.Word:
		{
			return WordExistsError{Word: entity.Polish}
		}
	case *dbmodels.Translation:
		{
			return TranslationExistsError{Translation: entity.English}
		}
	case *dbmodels.Sentence:
		{
			return SentenceExistsError{Sentence: entity.Sentence}
		}
	default:
		{
			break
		}
	}
	return fmt.Errorf("zły argument funkcji")
}

//errors for retrieving value from DB

type WordNotExistsError struct {
	Word string
}

func (e WordNotExistsError) Error() string {
	return fmt.Sprintf("słowa %s nie ma w słowniku", e.Word)
}

type SentenceNotExistsError struct {
	Word        string
	Translation string
	Sentence    string
}

func (e SentenceNotExistsError) Error() string {
	return fmt.Sprintf("zdanie %s prezentujące tłumaczenie %s słowa %s nie istnieje w słowniku", e.Sentence, e.Translation, e.Word)
}

type TranslationNotExistsError struct {
	Word        string
	Translation string
}

func (e TranslationNotExistsError) Error() string {
	return fmt.Sprintf("tłumaczenie %s słowa %s nie istnieje w słowniku", e.Translation, e.Word)
}

//errors for deleting values form DB

type CantDeleteWordError struct {
	Word string
}

func (e CantDeleteWordError) Error() string {
	return fmt.Sprintf("słowa %s nie ma w słowniku", e.Word)
}

type CantDeleteSentenceError struct {
	Sentence string
}

func (e CantDeleteSentenceError) Error() string {
	return fmt.Sprintf("zdanie %s dla podanego słowa i tłumaczenia nie istnieje. Sprawdź czy słowo i tłumaczenie są w słowniku", e.Sentence)
}

type CantDeleteTranslationError struct {
	Translation string
}

func (e CantDeleteTranslationError) Error() string {
	return fmt.Sprintf("tłumaczenie %s podanego słowa nie istnieje. Sprawdź czy słowo znajduje się w słowniku ", e.Translation)
}
