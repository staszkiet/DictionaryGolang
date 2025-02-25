package main

import (
	"bufio"
	"os"
	"strings"
)

type Reader struct {
	reader *bufio.Reader
}

var instance *Reader

func GetReaderInstance() *Reader {
	if instance == nil {
		instance = &Reader{
			reader: bufio.NewReader(os.Stdin)}
	}
	return instance
}

func (r *Reader) Read() string {
	input, _ := r.reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	return input
}
