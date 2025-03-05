package main

import (
	"bufio"
	"strings"
)

type Reader struct {
	reader *bufio.Reader
}

func (r Reader) Read() string {
	input, _ := r.reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	return input
}
