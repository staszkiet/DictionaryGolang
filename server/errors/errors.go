package customerrors

import "fmt"

type WordExistsError struct {
	Word string
}

func (e WordExistsError) Error() string {
	return fmt.Sprintf("słowo %s znajduje się już w słowniku. Aby dodać tłumaczenie użyj polecenia ADD_TRANSLATION", e.Word)
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
