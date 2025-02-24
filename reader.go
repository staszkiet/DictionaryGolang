package main

import (
	"bufio"
	"os"
	"strings"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader() *Reader {
	return &Reader{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (r *Reader) Read() string {
	input, _ := r.reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	return input
}
