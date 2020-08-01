package main

import (
	"fmt"
	//"github.com/tarm/serial"
	"os"
	"strings"
)

type Tokenizer struct {
	input string
	pos   int
	ch    byte
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input, 0, 0}
}

func (t *Tokenizer) Advance() Token {
	if t.pos >= len(t.input) {
		return Token{Eof, ""}
	}

	t.ch = t.input[t.pos]
	t.pos++
	switch {
	case strings.IndexByte("+-", t.ch) > -1:
		return Token{Reserved, string(t.ch)}
	default:
		return Token{Num, string(t.ch)}
	}
}

func (t *Tokenizer) Run() []Token {
	var tokens []Token
	token := t.Advance()

	for token.Kind != Eof {
		tokens = append(tokens, token)
		token = t.Advance()
	}
	tokens = append(tokens, token)
	return tokens
}

type TokenKind int

const (
	Num TokenKind = iota
	Reserved
	Eof
)

type Token struct {
	Kind    TokenKind
	Literal string
}

func main() {
	if len(os.Args) < 2 {
		panic("error: argument missing")
	}

	tokenizer := NewTokenizer(os.Args[1])
	fmt.Println(tokenizer.Run())
}
