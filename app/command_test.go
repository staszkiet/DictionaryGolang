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

func TestUpdateSentenceCommand_Execute(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	SetClientInstance(mockClient)

	cmd := UpdateSentenceCommand{request: graphql.NewRequest(
		`mutation UpdateSentence($polish: String!, $english: String!, $sentence: String! ,$newSentence: String!) 
			{updateSentence(polish: $polish, english: $english, sentence: $sentence ,newSentence: $newSentence)}`)}

	input := []string{"kot", "cat", "old cat", "new cat"}

	mockClient.On("Request", mock.Anything, mock.Anything).Return(nil)

	err := cmd.Execute(input)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestUpdateTranslationCommand_Execute(t *testing.T) {
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

func TestUpdateWordCommand_Execute(t *testing.T) {
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

func TestDeleteWordCommand_Execute(t *testing.T) {
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

func TestDeleteTranslationCommand_Execute(t *testing.T) {
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

func TestDeleteSentenceCommand_Execute(t *testing.T) {
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

func TestAddSentenceCommand_Execute(t *testing.T) {
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

func TestAddTranslationCommand_Execute(t *testing.T) {
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

func TestAddWordCommand_Execute(t *testing.T) {
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
