package customerrors

import "fmt"

//errors for adding values to DB

type WordExistsError struct {
	Word string
}

func (e WordExistsError) Error() string {
	return fmt.Sprintf("słowo %s znajduje się już w słowniku", e.Word)
}

type SentenceExistsError struct {
	Word        string
	Translation string
	Sentence    string
}

func (e SentenceExistsError) Error() string {
	return fmt.Sprintf("zdanie %s już jest dodane do tłumaczenia %s słowa %s", e.Sentence, e.Translation, e.Word)
}

type TranslationExistsError struct {
	Word        string
	Translation string
}

func (e TranslationExistsError) Error() string {
	return fmt.Sprintf("tłumaczenie %s już jest dodane do słowa %s", e.Translation, e.Word)
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
