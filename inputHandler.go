package main

import (
	"fmt"
)

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

func ListenForInput() {
	var action string
	reader := GetReaderInstance()
	commands := NewCommandFactory()
	for {

		fmt.Println("choose action:")
		fmt.Scanln(&action)
		if action == "exit" {
			break
		}
		command, exists := commands.GetCommand(action)
		if exists {
			polish := reader.Read()
			command.Execute(polish)
		} else {
			fmt.Println("Podane działanie nie istnieje")
		}
	}
}
