package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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
	reader := Reader{bufio.NewReader(os.Stdin)}
	commands := NewCommandFactory()
	fmt.Println("wybierz operację:\nADD - dodaj nowe słowo i jego tłumaczenie\nDELETE - usuń słowo\nSELECT - otrzymaj informacje o tłumaczeniu\n\nPolecenia modyfikujące istniejące tłumaczenia:\nADD TRANSLATION - dodaj tłumaczenie do słowa ze słownika\nDELETE TRANSLATION - usuń tłumaczenie\nADD SENTENCE - dodaj przykładowe zdanie do tłumaczenia\nDELETE SENTENCE - usuń przykładowe zdanie z danego tłumaczenia\nUPDATE - modyfikuje polską część\nUPDATE TRANSLATION - modyfikuje angielską częśc\nUPDATE SENTENCE - modyfikuje dane zdanie przykładowe")
	for {
		action = reader.Read()
		if action == "exit" {
			break
		}
		parsed := ParseInput(action)
		command, exists := commands.GetCommand(parsed[0])
		if exists {
			if err := command.Execute(parsed[1:]); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Podane działanie nie istnieje")
		}
	}
}

func ParseInput(input string) []string {
	pattern := `\([^\(\)]+\)|\S+`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllString(input, -1)

	for i, match := range matches {
		if strings.HasPrefix(match, "(") && strings.HasSuffix(match, ")") {
			matches[i] = match[1 : len(match)-1]
		}
	}

	return matches
}
