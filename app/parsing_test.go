package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputParsing(t *testing.T) {

	tests := []struct {
		input    string
		expected []string
	}{
		{"ADD word translation (example sentence) (example sentence 2)", []string{"ADD", "word", "translation", "example sentence", "example sentence 2"}},
		{"delete word", []string{"delete", "word"}},
		{"word", []string{"word"}},
		{"(sentence sentence)", []string{"sentence sentence"}},
	}

	for _, test := range tests {
		result := ParseInput(test.input)
		assert.Equal(t, test.expected, result, "they should be equal")
	}

}
