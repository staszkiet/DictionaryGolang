package customerrors

import "fmt"

type WordExistsError struct {
	Word string
}

func (e *WordExistsError) Error() string {
	return fmt.Sprintf("word '%s' already exists in the database", e.Word)
}
