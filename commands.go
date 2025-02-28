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
}

type DeleteSentenceCommand struct {
}

type AddTranslationCommand struct {
}

type DeleteTranslationCommand struct {
}

type AddWordCommand struct {
}

type DeleteWordCommand struct {
}

type SelectWordCommand struct {
}

type CommandFactory struct {
	commands map[string]ICommand
}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{
		commands: map[string]ICommand{
			"ADD TRANSLATION":    &AddTranslationCommand{},
			"DELETE TRANSLATION": &DeleteTranslationCommand{},
			"ADD SENTENCE":       &AddSentenceCommand{},
			"DELETE SENTENCE":    &DeleteSentenceCommand{},
			"ADD":                &AddWordCommand{},
			"DELETE":             &DeleteWordCommand{},
			"SELECT":             &SelectWordCommand{},
		},
	}
}

func (f *CommandFactory) GetCommand(action string) (ICommand, bool) {
	strategy, exists := f.commands[action]
	return strategy, exists
}

func (s SelectWordCommand) Execute(polish string) error {
	graphqlClient := GetClientInstance()
	graphqlRequest := graphql.NewRequest(`
				query selectWord($polish: String!) {
			selectWord(polish: $polish){
	  translations{
		english
		sentences{
		  sentence
		}
	  }
	
			}
		}
	`)

	graphqlRequest.Var("polish", polish)

	var graphqlResponse SelectResponse

	if err := graphqlClient.client.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	PrintSelectOutput(graphqlResponse, polish)

	return nil
}

func (d DeleteWordCommand) Execute(polish string) error {
	graphqlClient := GetClientInstance()
	graphqlRequest := graphql.NewRequest(`
				mutation deleteWord($polish: String!) {
			deleteWord(polish: $polish)
		}
	`)

	graphqlRequest.Var("polish", polish)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
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
	graphqlRequest := graphql.NewRequest(`
				mutation CreateWord($polish: String!, $translation: NewTranslation!) {
			createWord(polish: $polish, translation: $translation)
		}
	`)
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

	graphqlRequest.Var("polish", polish)
	newTran := model.NewTranslation{English: translation, Sentences: sentences}
	graphqlRequest.Var("translation", newTran)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
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
	graphqlRequest := graphql.NewRequest(`
				mutation deleteTranslation($polish: String!, $english: String!) {
			deleteTranslation(polish: $polish, english: $english)
		}
	`)

	graphqlRequest.Var("polish", polish)
	graphqlRequest.Var("english", translation)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	return nil
}

func (a AddTranslationCommand) Execute(polish string) error {
	var translation, sentence string
	reader := GetReaderInstance()
	sentences := []string{}

	graphqlClient := GetClientInstance()
	graphqlRequest := graphql.NewRequest(`
				mutation CreateTranslation($polish: String!, $translation: NewTranslation!) {
			createTranslation(polish: $polish, translation: $translation)
		}
	`)
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

	graphqlRequest.Var("polish", polish)
	newTran := model.NewTranslation{English: translation, Sentences: sentences}
	graphqlRequest.Var("translation", newTran)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
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
	graphqlRequest := graphql.NewRequest(`
				mutation deleteSentence($polish: String!, $english: String!, $sentence: String!) {
			deleteSentence(polish: $polish, english: $english, sentence: $sentence)
		}
	`)

	graphqlRequest.Var("polish", polish)
	graphqlRequest.Var("english", translation)
	graphqlRequest.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
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
	graphqlRequest := graphql.NewRequest(`
				mutation createSentence($polish: String!, $english: String!, $sentence: String!) {
			createSentence(polish: $polish, english: $english, sentence: $sentence)
		}
	`)

	graphqlRequest.Var("polish", polish)
	graphqlRequest.Var("english", translation)
	graphqlRequest.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	return nil
}
