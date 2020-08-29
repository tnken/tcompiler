package token

import "strings"

// Tokenizer has source code and read position
type Tokenizer struct {
	input string
	pos   int
}

func newTokenizer(input string) *Tokenizer {
	return &Tokenizer{input, 0}
}

func (t *Tokenizer) recognizeMany(f func(byte) bool) {
	for t.pos < len(t.input) && f(t.input[t.pos]) {
		t.pos++
	}
}

func isChar(b byte) bool {
	return 'a' <= b && b <= 'z'
}

func isDigit(b byte) bool {
	return (strings.IndexByte("0123456789", b) > -1)
}

func isAlnum(b byte) bool {
	return isChar(b) || isDigit(b)
}

func (t *Tokenizer) lexNumber() Token {
	start := t.pos
	t.recognizeMany(isDigit)
	return Token{Num, t.input[start:t.pos]}
}

func (t *Tokenizer) lexIdent() Token {
	start := t.pos
	t.recognizeMany(isAlnum)
	return Token{Identifier, t.input[start:t.pos]}
}

func (t *Tokenizer) lexSpaces() {
	t.recognizeMany(func(b byte) bool { return (strings.IndexByte(" \n\t", b) > -1) })
}

func (t *Tokenizer) next() Token {
	// TODO: Refactor from LL:50~LL:62
	if t.pos >= len(t.input) {
		return t.newToken(EOF, "")
	}
	ch := t.input[t.pos]
	if ch == ' ' || ch == '\t' || ch == '\n' {
		t.lexSpaces()
	}

	if t.pos >= len(t.input) {
		return t.newToken(EOF, "")
	}
	ch = t.input[t.pos]

	switch {
	case ch == '+':
		return t.newToken(Plus, string(ch))
	case ch == '-':
		return t.newToken(Minus, string(ch))
	case ch == '*':
		return t.newToken(Asterisk, string(ch))
	case ch == '/':
		return t.newToken(Slash, string(ch))
	case ch == '[':
		return t.newToken(Lbracket, string(ch))
	case ch == ']':
		return t.newToken(Rbracket, string(ch))
	case ch == '(':
		return t.newToken(LParen, string(ch))
	case ch == ')':
		return t.newToken(RParen, string(ch))
	case ch == ',':
		return t.newToken(Comma, string(ch))
	case ch == '=':
		return t.newToken(Assign, string(ch))
	case ch == '{':
		return t.newToken(Lbrace, string(ch))
	case ch == '}':
		return t.newToken(Rbrace, string(ch))
	case t.isReserved():
		for _, v := range reserved {
			if t.input[t.pos:t.pos+len(v)] == v {
				return t.newToken(reservedToKind[t.input[t.pos:t.pos+len(v)]], t.input[t.pos:t.pos+len(v)])
			}
		}
	case isDigit(ch):
		return t.lexNumber()
	case isChar((ch)):
		return t.lexIdent()
	}
	return t.newToken(EOF, "")
}

// Kind express kind of the token as enum
type Kind int

// Define Kind as enum
const (
	Num      Kind = iota // 0 - 9
	Plus                 // +
	Minus                // -
	Asterisk             // *
	Slash                // /
	Lbracket             // [
	Rbracket             // ]
	LParen               // (
	RParen               // )
	Assign               // =
	Comma                // ,
	Lbrace               // {
	Rbrace               // }
	Identifier
	EOF
	KeyDo
	KeyEnd
	KeyLoop
)

var reserved = []string{
	"do",
	"end",
	"loop",
}

var reservedToKind = map[string]Kind{
	"do":   KeyDo,
	"end":  KeyEnd,
	"loop": KeyLoop,
}

func (t Tokenizer) isReserved() bool {
	for _, v := range reserved {
		if len(t.input)-t.pos <= len(v) {
			continue
		}
		if t.input[t.pos:t.pos+len(v)] == v {
			return true
		}
	}
	return false
}

// Token consits of its kind and literal
type Token struct {
	Kind    Kind
	Literal string
}

func (t *Tokenizer) newToken(tk Kind, lit string) Token {
	for i := 0; i < len(lit); i++ {
		t.pos++
	}
	return Token{tk, lit}
}
