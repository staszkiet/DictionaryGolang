package main

import (
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGraphQLClient struct {
	mock.Mock
}

func (m *MockGraphQLClient) Request(req *graphql.Request, response interface{}) error {
	args := m.Called(req, response)
	return args.Error(0)
}

func TestUpdateTranslationCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := UpdateTranslationCommand{request: graphql.NewRequest(
		`mutation UpdateTranslation($polish: String!, $english: String!, $newEnglish: String!) 
	{updateTranslation(polish: $polish, english: $english, newEnglish: $newEnglish)}`)}

	input := []string{"kot", "cst", "cat"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestUpdateTranslationCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := UpdateTranslationCommand{request: graphql.NewRequest(
		`mutation UpdateTranslation($polish: String!, $english: String!, $newEnglish: String!) 
	{updateTranslation(polish: $polish, english: $english, newEnglish: $newEnglish)}`)}

	input := []string{"kot", "cst"}

	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")
}

func TestUpdateWordCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := UpdateWordCommand{request: graphql.NewRequest(`mutation UpdateWord($polish: String!, $newPolish: String!) 
	{updateWord(polish: $polish, newPolish: $newPolish)}`)}

	input := []string{"kst", "kot"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestUpdateWordCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := UpdateWordCommand{request: graphql.NewRequest(`mutation UpdateWord($polish: String!, $newPolish: String!) 
	{updateWord(polish: $polish, newPolish: $newPolish)}`)}

	input := []string{"kst", "kot", "(kot zdanie)"}

	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")
}

func TestDeleteWordCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := DeleteWordCommand{request: graphql.NewRequest(
		`mutation DeleteWord($polish: String!) 
	{deleteWord(polish: $polish)}`)}

	input := []string{"kot"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestDeleteWordCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := DeleteWordCommand{request: graphql.NewRequest(
		`mutation DeleteWord($polish: String!) 
	{deleteWord(polish: $polish)}`)}

	input := []string{"kot", "cat"}
	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")

}

func TestDeleteTranslationCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := DeleteTranslationCommand{request: graphql.NewRequest(`
	mutation deleteTranslation($polish: String!, $english: String!) 
	{deleteTranslation(polish: $polish, english: $english)}`)}

	input := []string{"kot", "cat"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestDeleteTranslationCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := DeleteTranslationCommand{request: graphql.NewRequest(`
	mutation deleteTranslation($polish: String!, $english: String!) 
	{deleteTranslation(polish: $polish, english: $english)}`)}

	input := []string{"kot", "cat", "(zdanie zdanie)"}

	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")
}

func TestDeleteSentenceCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := DeleteSentenceCommand{request: graphql.NewRequest(`
	mutation deleteSentence($polish: String!, $english: String!, $sentence: String!) {
	deleteSentence(polish: $polish, english: $english, sentence: $sentence)}`)}

	input := []string{"kot", "cat", "I hate my cat"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestDeleteSentenceCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := DeleteSentenceCommand{request: graphql.NewRequest(`
	mutation deleteSentence($polish: String!, $english: String!, $sentence: String!) {
	deleteSentence(polish: $polish, english: $english, sentence: $sentence)}`)}

	input := []string{"kot", "cat"}

	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")
}

func TestAddSentenceCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := AddSentenceCommand{request: graphql.NewRequest(`
	mutation createSentence($polish: String!, $english: String!, $sentence: String!) {
	createSentence(polish: $polish, english: $english, sentence: $sentence)}`)}

	input := []string{"kot", "cat", "I hate my cat"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAddSentenceCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := AddSentenceCommand{request: graphql.NewRequest(`
	mutation createSentence($polish: String!, $english: String!, $sentence: String!) {
	createSentence(polish: $polish, english: $english, sentence: $sentence)}`)}

	input := []string{"kot", "cat", "I hate my cat", "I like my cat"}

	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")
}

func TestAddTranslationCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := AddTranslationCommand{request: graphql.NewRequest(`
	mutation CreateTranslation($polish: String!, $translation: NewTranslation!) {
	createTranslation(polish: $polish, translation: $translation)}`)}

	input := []string{"kot", "cat", "I hate my cat", "I love my cat"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAddTranslationCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := AddTranslationCommand{request: graphql.NewRequest(`
	mutation CreateTranslation($polish: String!, $translation: NewTranslation!) {
	createTranslation(polish: $polish, translation: $translation)}`)}

	input := []string{"kot"}

	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")
}

func TestAddWordCommand_Execute_ValidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := AddWordCommand{request: graphql.NewRequest(`
	mutation CreateTranslation($polish: String!, $translation: NewTranslation!) {
	createTranslation(polish: $polish, translation: $translation)}`)}

	input := []string{"kot", "cat", "I hate my cat", "I love my cat"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestAddWordCommand_Execute_InvalidInput(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := AddWordCommand{request: graphql.NewRequest(`
	mutation CreateTranslation($polish: String!, $translation: NewTranslation!) {
	createTranslation(polish: $polish, translation: $translation)}`)}

	input := []string{"kot"}

	err := cmd.Execute(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "niepoprawna liczba argumentów")
}
