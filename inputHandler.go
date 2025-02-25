package main

import (
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/staszkiet/DictionaryGolang/graph/model"
)

type IHandler interface {
	PerformAction(string, string) bool
}

type AddHandler struct {
	next IHandler
}

type DeleteHandler struct {
	next IHandler
}

type CreateHandler struct {
	next IHandler
}

type UpdateHandler struct {
	next IHandler
}

func (a AddHandler) PerformAction(word string, action string) bool {
	if action == "ADD" {
		var translation, sentence string
		reader := GetReaderInstance()
		sentences := []string{}

		graphqlClient := GetClientInstance()
		graphqlRequest := graphql.NewRequest(`
					mutation CreateWord($polish: String!, $translation: NewTranslation!) {
				createWord(polish: $polish, translation: $translation) {
					polish
				}
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

		graphqlRequest.Var("polish", word)
		newTran := model.NewTranslation{English: translation, Sentences: sentences}
		graphqlRequest.Var("translation", newTran)

		var graphqlResponse interface{}

		if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
			panic(err)
		}

		fmt.Println(graphqlResponse)
		return true
	}
	if a.next == nil {
		return false
	} else {
		return a.next.PerformAction(word, action)
	}

}

func (c CreateHandler) PerformAction(word string, action string) bool {
	if action == "CREATE" {

		//perform action
		fmt.Println(word)
		return true
	}
	if c.next == nil {
		return false
	} else {
		return c.next.PerformAction(word, action)
	}
}

func (d DeleteHandler) PerformAction(word string, action string) bool {
	if action == "DELETE" {

		//perform action
		fmt.Println(word)
		return true
	}
	if d.next == nil {
		return false
	} else {
		return d.next.PerformAction(word, action)
	}
}

func (u UpdateHandler) PerformAction(word string, action string) bool {
	if action == "UPDATE" {

		//perform action
		fmt.Println(word)
		return true
	}
	if u.next == nil {
		return false
	} else {
		return u.next.PerformAction(word, action)
	}
}

func ListenForInput() {
	var action string
	var word string
	add := AddHandler{next: CreateHandler{next: DeleteHandler{next: UpdateHandler{}}}}
	reader := GetReaderInstance()
	for {

		fmt.Println("choose action:")
		fmt.Scanln(&action)
		word = reader.Read()
		add.PerformAction(word, action)
	}
}
