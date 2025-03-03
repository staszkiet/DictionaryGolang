package main

import (
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/staszkiet/DictionaryGolang/server/graph/model"
)

type ICommand interface {
	Execute(input []string) error
}

type AddSentenceCommand struct {
	request *graphql.Request
}

type DeleteSentenceCommand struct {
	request *graphql.Request
}

type AddWordCommand struct {
	request *graphql.Request
}

type AddTranslationCommand struct {
	request *graphql.Request
}

type DeleteTranslationCommand struct {
	request *graphql.Request
}

type DeleteWordCommand struct {
	request *graphql.Request
}

type SelectWordCommand struct {
	request *graphql.Request
}

type UpdateWordCommand struct {
	request *graphql.Request
}

type UpdateTranslationCommand struct {
	request *graphql.Request
}

type UpdateSentenceCommand struct {
	request *graphql.Request
}

type CommandFactory struct {
	commands map[string]ICommand
}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{
		commands: map[string]ICommand{
			"ADD TRANSLATION": &AddTranslationCommand{request: graphql.NewRequest(`
				mutation CreateTranslation($polish: String!, $translation: NewTranslation!) {
			createTranslation(polish: $polish, translation: $translation)}`)},

			"ADD": &AddWordCommand{request: graphql.NewRequest(`
			mutation CreateWord($polish: String!, $translation: NewTranslation!) {
		createWord(polish: $polish, translation: $translation)}`)},

			"DELETE TRANSLATION": &DeleteTranslationCommand{request: graphql.NewRequest(`
				mutation deleteTranslation($polish: String!, $english: String!) 
				{deleteTranslation(polish: $polish, english: $english)}`)},

			"ADD SENTENCE": &AddSentenceCommand{request: graphql.NewRequest(`
			mutation createSentence($polish: String!, $english: String!, $sentence: String!) {
			createSentence(polish: $polish, english: $english, sentence: $sentence)}`)},

			"DELETE SENTENCE": &DeleteSentenceCommand{request: graphql.NewRequest(`
			mutation deleteSentence($polish: String!, $english: String!, $sentence: String!) {
			deleteSentence(polish: $polish, english: $english, sentence: $sentence)}`)},

			"DELETE": &DeleteWordCommand{request: graphql.NewRequest(
				`mutation DeleteWord($polish: String!) 
			{deleteWord(polish: $polish)}`)},

			"SELECT": &SelectWordCommand{request: graphql.NewRequest(`query selectWord($polish: String!) 
			{selectWord(polish: $polish){translations{english sentences{sentence}}}}`)},

			"UPDATE": &UpdateWordCommand{request: graphql.NewRequest(`mutation UpdateWord($polish: String!, $newPolish: String!) 
			{updateWord(polish: $polish, newPolish: $newPolish)}`)},
			"UPDATE TRANSLATION": &UpdateTranslationCommand{request: graphql.NewRequest(
				`mutation UpdateTranslation($polish: String!, $english: String!, $newEnglish: String!) 
			{updateTranslation(polish: $polish, english: $english, newEnglish: $newEnglish)}`)},
			"UPDATE SENTENCE": &UpdateSentenceCommand{request: graphql.NewRequest(
				`mutation UpdateSentence($polish: String!, $english: String!, $sentence: String! ,$newSentence: String!) 
			{updateSentence(polish: $polish, english: $english, sentence: $sentence ,newSentence: $newSentence)}`)},
		},
	}
}

func (f *CommandFactory) GetCommand(action string) (ICommand, bool) {
	command, exists := f.commands[action]
	return command, exists
}

func (s SelectWordCommand) Execute(input []string) error {

	polish := input[0]

	graphqlClient := GetClientInstance()

	s.request.Var("polish", polish)

	var graphqlResponse SelectResponse

	if err := graphqlClient.Request(s.request, &graphqlResponse); err != nil {
		return err
	}

	PrintSelectOutput(graphqlResponse, polish)

	return nil
}

func (u UpdateSentenceCommand) Execute(input []string) error {

	polish := input[0]
	English := input[1]
	Sentence := input[2]
	newSentence := input[3]

	graphqlClient := GetClientInstance()
	u.request.Var("polish", polish)
	u.request.Var("english", English)
	u.request.Var("sentence", Sentence)
	u.request.Var("newSentence", newSentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(u.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (u UpdateTranslationCommand) Execute(input []string) error {
	polish := input[0]
	English := input[1]
	newEnglish := input[2]
	graphqlClient := GetClientInstance()
	u.request.Var("polish", polish)
	u.request.Var("english", English)
	u.request.Var("newEnglish", newEnglish)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(u.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (u UpdateWordCommand) Execute(input []string) error {

	graphqlClient := GetClientInstance()
	polish := input[0]
	newPolish := input[1]
	u.request.Var("polish", polish)
	u.request.Var("newPolish", newPolish)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(u.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (d DeleteWordCommand) Execute(input []string) error {

	graphqlClient := GetClientInstance()
	polish := input[0]
	d.request.Var("polish", polish)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(d.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (a AddWordCommand) Execute(input []string) error {

	sentences := []string{}

	graphqlClient := GetClientInstance()
	polish := input[0]
	translation := input[1]
	for i := 2; i < len(input); i++ {
		sentences = append(sentences, input[i])
	}

	a.request.Var("polish", polish)
	newTran := model.NewTranslation{English: translation, Sentences: sentences}
	a.request.Var("translation", newTran)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(a.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (d DeleteTranslationCommand) Execute(input []string) error {
	polish := input[0]
	translation := input[1]
	graphqlClient := GetClientInstance()

	d.request.Var("polish", polish)
	d.request.Var("english", translation)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(d.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (a AddTranslationCommand) Execute(input []string) error {

	sentences := []string{}
	polish := input[0]
	graphqlClient := GetClientInstance()

	fmt.Println("translation:")

	translation := input[1]
	fmt.Println("example sentences:")
	for i := 2; i < len(input); i++ {
		sentences = append(sentences, input[i])
	}

	a.request.Var("polish", polish)
	newTran := model.NewTranslation{English: translation, Sentences: sentences}
	a.request.Var("translation", newTran)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(a.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (d DeleteSentenceCommand) Execute(input []string) error {

	polish := input[0]
	translation := input[1]
	sentence := input[2]
	graphqlClient := GetClientInstance()

	d.request.Var("polish", polish)
	d.request.Var("english", translation)
	d.request.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(d.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (a AddSentenceCommand) Execute(input []string) error {

	polish := input[0]
	translation := input[1]
	sentence := input[2]
	graphqlClient := GetClientInstance()

	a.request.Var("polish", polish)
	a.request.Var("english", translation)
	a.request.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(a.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}
