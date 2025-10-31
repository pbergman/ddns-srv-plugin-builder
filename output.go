package main

import (
	"bufio"
	"bytes"
	"log"
)

type output struct {
	out    *log.Logger
	prefix string
}

func (o *output) Write(p []byte) (n int, err error) {

	var scanner = bufio.NewScanner(bytes.NewReader(p))

	for scanner.Scan() {
		o.out.Printf("%s%s\n", o.prefix, scanner.Text())
	}

	return len(p), nil
}
