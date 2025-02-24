package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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

		reader := bufio.NewReader(os.Stdin)

		graphqlClient := graphql.NewClient("http://localhost:8080/query")
		graphqlRequest := graphql.NewRequest(`
					mutation CreateWord($polish: String!, $translation: NewTranslation!) {
				createWord(polish: $polish, translation: $translation) {
					polish
				}
			}
		`)
		fmt.Println("translation:")
		translation, _ := reader.ReadString('\n')
		translation = strings.Replace(translation, "\n", "", -1)
		sentences := []string{}
		fmt.Println("example sentences:")
		for {
			sentence, _ := reader.ReadString('\n')
			sentence = strings.Replace(sentence, "\n", "", -1)
			if sentence == "" {
				break
			}

			sentences = append(sentences, sentence)
		}

		graphqlRequest.Var("polish", word)
		newTran := model.NewTranslation{English: translation, Sentences: sentences}
		graphqlRequest.Var("translation", newTran)

		var graphqlResponse interface{}
		if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
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
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Println("choose action:")
		fmt.Scanln(&action)
		word, _ = reader.ReadString('\n')
		word = strings.Replace(word, "\n", "", -1)
		add.PerformAction(word, action)
	}
}
