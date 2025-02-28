package main

import (
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/staszkiet/DictionaryGolang/graph/model"
)

type ICommand interface {
	Execute() error
}

type AddSentenceCommand struct {
	polish string
}

type DeleteSentenceCommand struct {
	polish string
}

type AddTranslationCommand struct {
	polish string
}

type DeleteTranslationCommand struct {
	polish string
}

func (d DeleteTranslationCommand) Execute() error {
	reader := GetReaderInstance()
	fmt.Println("translation:")
	translation := reader.Read()
	graphqlClient := GetClientInstance()
	graphqlRequest := graphql.NewRequest(`
				mutation deleteTranslation($polish: String!, $english: String!) {
			deleteTranslation(polish: $polish, english: $english)
		}
	`)

	graphqlRequest.Var("polish", d.polish)
	graphqlRequest.Var("english", translation)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	return nil
}

func (a AddTranslationCommand) Execute() error {
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

	graphqlRequest.Var("polish", a.polish)
	newTran := model.NewTranslation{English: translation, Sentences: sentences}
	graphqlRequest.Var("translation", newTran)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Pomyślnie dodano tłumaczenie do słownika!")

	return nil
}

func (d DeleteSentenceCommand) Execute() error {
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

	graphqlRequest.Var("polish", d.polish)
	graphqlRequest.Var("english", translation)
	graphqlRequest.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	return nil
}

func (a AddSentenceCommand) Execute() error {
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

	graphqlRequest.Var("polish", a.polish)
	graphqlRequest.Var("english", translation)
	graphqlRequest.Var("sentence", sentence)

	var graphqlResponse interface{}

	if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
		panic(err)
	}

	return nil
}
