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

		//perform action
		fmt.Println(word)
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
	for {

		fmt.Println("choose action:")
		fmt.Scanln(&action)
		fmt.Scanln(&word)
		add.PerformAction(word, action)
	}
}
