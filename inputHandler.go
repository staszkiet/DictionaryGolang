package main

import (
	"fmt"
)

type IHandler interface {
	PerformAction(string, string) bool
}

type AddHandler struct {
	next IHandler
}

func (a AddHandler) PerformAction(word string, action string) bool {
	if action == "ADD" {
		fmt.Println(word)
		return true
	}
	if a.next == nil {
		return false
	} else {
		return a.next.PerformAction(word, action)
	}
}

func ListenForInput() {
	var action string
	var word string
	add := AddHandler{next: nil}
	for {

		fmt.Println("choose action:")
		fmt.Scanln(&action)
		fmt.Scanln(&word)
		add.PerformAction(word, action)
	}
}
