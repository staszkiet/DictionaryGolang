package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/staszkiet/DictionaryGolang/graph/model"
)

type ICommand interface {
	Execute(polish string) error
}

type AddSentenceCommand struct {
	request *graphql.Request
}

type DeleteSentenceCommand struct {
	request *graphql.Request
}

type AddTranslationCommand struct {
	request *graphql.Request
}

type DeleteTranslationCommand struct {
	request *graphql.Request
}

type AddWordCommand struct {
	request *graphql.Request
}

type DeleteWordCommand struct {
	request *graphql.Request
}

type SelectWordCommand struct {
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

			"DELETE TRANSLATION": &DeleteTranslationCommand{request: graphql.NewRequest(`
				mutation deleteTranslation($polish: String!, $english: String!) 
				{deleteTranslation(polish: $polish, english: $english)}`)},

			"ADD SENTENCE": &AddSentenceCommand{request: graphql.NewRequest(`
			mutation createSentence($polish: String!, $english: String!, $sentence: String!) {
			createSentence(polish: $polish, english: $english, sentence: $sentence)}`)},

			"DELETE SENTENCE": &DeleteSentenceCommand{request: graphql.NewRequest(`
			mutation deleteSentence($polish: String!, $english: String!, $sentence: String!) {
			deleteSentence(polish: $polish, english: $english, sentence: $sentence)}`)},

			"ADD": &AddWordCommand{request: graphql.NewRequest(`
			mutation CreateWord($polish: String!, $translation: NewTranslation!)
			 {createWord(polish: $polish, translation: $translation)}`)},

			"DELETE": &DeleteWordCommand{request: graphql.NewRequest(
				`mutation CreateWord($polish: String!, $translation: NewTranslation!) 
			{createWord(polish: $polish, translation: $translation)}`)},

			"SELECT": &SelectWordCommand{request: graphql.NewRequest(`query selectWord($polish: String!) 
			{selectWord(polish: $polish){translations{english sentences{sentence}}}}`)},
		},
	}
}

func (f *CommandFactory) GetCommand(action string) (ICommand, bool) {
	strategy, exists := f.commands[action]
	return strategy, exists
}

func (s SelectWordCommand) Execute(polish string) error {
	graphqlClient := GetClientInstance()

	s.request.Var("polish", polish)

	var graphqlResponse SelectResponse

	if err := graphqlClient.client.Run(context.Background(), s.request, &graphqlResponse); err != nil {
		panic(err)
	}

	PrintSelectOutput(graphqlResponse, polish)

	return nil
}

func (d DeleteWordCommand) Execute(polish string) error {
	graphqlClient := GetClientInstance()

	d.request.Var("polish", polish)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(d.request, &graphqlResponse); err != nil {
		return err
	}

	fmt.Println(graphqlResponse)
	return nil
}

func (a AddWordCommand) Execute(polish string) error {
	var translation, sentence string
	reader := GetReaderInstance()
	sentences := []string{}

	graphqlClient := GetClientInstance()

	fmt.Println("translation:")

	translation = reader.Read()
	fmt.Println("example sentences:")
	for {
		sentence = reader.Read()
		if sentence == "" {
			break
		}
		sentences = append(sentences, sentence)
	}

	a.request.Var("polish", polish)
	newTran := model.NewTranslation{English: translation, Sentences: sentences}
	a.request.Var("translation", newTran)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(a.request, &graphqlResponse); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Pomyślnie dodano tłumaczenie do słownika!")

	return nil
}

func (d DeleteTranslationCommand) Execute(polish string) error {
	reader := GetReaderInstance()
	fmt.Println("translation:")
	translation := reader.Read()
	graphqlClient := GetClientInstance()

	d.request.Var("polish", polish)
	d.request.Var("english", translation)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(d.request, &graphqlResponse); err != nil {
		return err
	}

	return nil
}

func (a AddTranslationCommand) Execute(polish string) error {
	var translation, sentence string
	reader := GetReaderInstance()
	sentences := []string{}

	graphqlClient := GetClientInstance()

	fmt.Println("translation:")

	translation = reader.Read()
	fmt.Println("example sentences:")
	for {
		sentence = reader.Read()
		if sentence == "" {
			break
		}
		sentences = append(sentences, sentence)
	}

	a.request.Var("polish", polish)
	newTran := model.NewTranslation{English: translation, Sentences: sentences}
	a.request.Var("translation", newTran)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(a.request, &graphqlResponse); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Pomyślnie dodano tłumaczenie do słownika!")

	return nil
}

func (d DeleteSentenceCommand) Execute(polish string) error {
	reader := GetReaderInstance()
	fmt.Println("translation:")
	translation := reader.Read()
	fmt.Println("sentence:")
	sentence := reader.Read()
	graphqlClient := GetClientInstance()

	d.request.Var("polish", polish)
	d.request.Var("english", translation)
	d.request.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(d.request, &graphqlResponse); err != nil {
		panic(err)
	}

	return nil
}

func (a AddSentenceCommand) Execute(polish string) error {
	reader := GetReaderInstance()
	fmt.Println("translation:")
	translation := reader.Read()
	fmt.Println("sentence:")
	sentence := reader.Read()
	graphqlClient := GetClientInstance()

	a.request.Var("polish", polish)
	a.request.Var("english", translation)
	a.request.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(a.request, &graphqlResponse); err != nil {
		panic(err)
	}

	return nil
}
