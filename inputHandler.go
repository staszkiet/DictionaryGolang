package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/staszkiet/DictionaryGolang/graph/model"
)

type IHandler interface {
	PerformAction(string, string) bool
}

type AddHandler struct {
	next IHandler
	req  *graphql.Request
}

type DeleteHandler struct {
	next IHandler
	req  *graphql.Request
}

type SelectHandler struct {
	next IHandler
	req  *graphql.Request
}

type UpdateHandler struct {
	next IHandler
	req  *graphql.Request
}

type SelectResponse struct {
	SelectWord struct {
		Translations []struct {
			English   string `json:"english"`
			Sentences []struct {
				Sentence string `json:"sentence"`
			} `json:"sentences"`
		} `json:"translations"`
	} `json:"selectWord"`
}

func (a AddHandler) PerformAction(word string, action string) bool {
	if action == "ADD" {

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

		graphqlRequest.Var("polish", word)
		newTran := model.NewTranslation{English: translation, Sentences: sentences}
		graphqlRequest.Var("translation", newTran)

		var graphqlResponse interface{}

		if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
			fmt.Println(err.Error())
			return true
		}

		fmt.Println("Pomyślnie dodano tłumaczenie do słownika!")

		return true
	}
	if a.next == nil {
		return false
	} else {
		return a.next.PerformAction(word, action)
	}

}

func (c SelectHandler) PerformAction(word string, action string) bool {
	if action == "SELECT" {

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

		graphqlRequest.Var("polish", word)

		var graphqlResponse SelectResponse

		if err := graphqlClient.client.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
			panic(err)
		}

		PrintSelectOutput(graphqlResponse, word)

		return true
	}
	if c.next == nil {
		return false
	} else {
		return c.next.PerformAction(word, action)
	}
}

func PrintSelectOutput(response SelectResponse, polish string) {
	fmt.Printf("\n\nTłumaczenia dla słowa %s\n\n", polish)
	for _, t := range response.SelectWord.Translations {
		fmt.Printf("%s\n\n", t.English)
		fmt.Printf("Przykładowe zdania:\n\n")
		for _, s := range t.Sentences {
			fmt.Printf("%s\n", s.Sentence)
		}
	}
	fmt.Printf("\n\n")
}

func (d DeleteHandler) PerformAction(word string, action string) bool {
	if action == "DELETE" {

		graphqlClient := GetClientInstance()
		graphqlRequest := graphql.NewRequest(`
					mutation deleteWord($polish: String!) {
				deleteWord(polish: $polish)
			}
		`)

		graphqlRequest.Var("polish", word)

		var graphqlResponse interface{}

		if err := graphqlClient.Request(graphqlRequest, &graphqlResponse); err != nil {
			panic(err)
		}

		fmt.Println(graphqlResponse)
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

		command := AddSentenceCommand{polish: word}
		command.Execute()

		return true
	}
	if u.next == nil {
		return false
	} else {
		return u.next.PerformAction(word, action)
	}
}

type IUpdateCommand interface {
	Execute() error
}

type AddSentenceCommand struct {
	polish string
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

func ListenForInput() {
	var action string
	var word string
	add := AddHandler{next: SelectHandler{next: DeleteHandler{next: UpdateHandler{}}}}
	reader := GetReaderInstance()
	for {

		fmt.Println("choose action:")
		fmt.Scanln(&action)
		if action == "exit" {
			break
		}
		word = reader.Read()
		add.PerformAction(word, action)
	}
}
