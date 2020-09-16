package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument error")
		return
	}
	tok := NewToken(os.Args[1])
	parser := NewParser(tok)
	c := Compile(parser.Program())
	c.output()
}
