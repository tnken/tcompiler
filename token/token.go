package token

import (
	"errors"
	"fmt"
	"strings"
)

// Tokenizer has source code and read position
type Tokenizer struct {
	input string
	pos   int
}

type TokenizeErr struct {
	Err error
	L   Loc
	t   *Tokenizer
}

// custom error
var (
	ErrSyntax = errors.New("Syntax error, undefined token")
)

func (te *TokenizeErr) Error() string {
	switch te.Err {
	case ErrSyntax:
		st, l := lineNum(te.t.input, te.L.Start)
		line := displayLine(te.t.input, te.L.Start)
		return fmt.Sprintf("%d:%d: %v\n%v", l, te.L.Start-st+1, te.Err, line)
	}
	return te.Err.Error()
}

// posは，s中のstart番目から始まるline行目の文字
func lineNum(s string, pos int) (int, int) {
	line := 1
	start := 0
	for i := 0; i < len(s); i++ {
		// CF: 0x0A, CD: 0x0D
		if s[i] == 10 || s[i] == 13 {
			start = (i + 1)
			line++
		}

		if i == pos {
			break
		}
	}
	return start, line
}

func displayLine(s string, pos int) string {
	st, _ := lineNum(s, pos)
	line := ""
	i := st
	for i < len(s) {
		// CF: 0x0A, CD: 0x0D
		if s[i] == 10 || s[i] == 13 {
			break
		}
		line += s[i : i+1]
		i++
	}
	line += "\n"
	for j := 0; j < pos-st; j++ {
		line += " "
	}
	line += "^\n"
	return line
}

// New initialize a Tokenizer and returns its pointer
func New(input string) *Tokenizer {
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
	return Token{Num, t.input[start:t.pos], Loc{start, t.pos}}
}

func (t *Tokenizer) lexIdent() Token {
	start := t.pos
	t.recognizeMany(isAlnum)
	return Token{Identifier, t.input[start:t.pos], Loc{start, t.pos}}
}

func (t *Tokenizer) lexSpaces() {
	t.recognizeMany(func(b byte) bool { return (strings.IndexByte(" \n\t", b) > -1) })
}

// Next returns a Token and move forward current position
func (t *Tokenizer) Next() (Token, error) {
	// TODO: Refactoring from LL:51 to LL:62
	if t.pos >= len(t.input) {
		return t.newToken(EOF, ""), nil
	}
	ch := t.input[t.pos]

	if ch == ' ' || ch == '\t' || ch == '\n' {
		t.lexSpaces()
	}

	if t.pos >= len(t.input) {
		return t.newToken(EOF, ""), nil
	}
	ch = t.input[t.pos]

	switch {
	case ch == '+':
		return t.newToken(Plus, string(ch)), nil
	case ch == '-':
		return t.newToken(Minus, string(ch)), nil
	case ch == '*':
		return t.newToken(Asterisk, string(ch)), nil
	case ch == '/':
		return t.newToken(Slash, string(ch)), nil
	case ch == '[':
		return t.newToken(Lbracket, string(ch)), nil
	case ch == ']':
		return t.newToken(Rbracket, string(ch)), nil
	case ch == '(':
		return t.newToken(LParen, string(ch)), nil
	case ch == ')':
		return t.newToken(RParen, string(ch)), nil
	case ch == ',':
		return t.newToken(Comma, string(ch)), nil
	case ch == '=':
		if t.input[t.pos+1] == '=' {
			return t.newToken(Eq, "=="), nil
		}
		return t.newToken(Assign, string(ch)), nil
	case ch == '{':
		return t.newToken(Lbrace, string(ch)), nil
	case ch == '}':
		return t.newToken(Rbrace, string(ch)), nil
	case ch == '!':
		if t.input[t.pos+1] == '=' {
			return t.newToken(NEq, "!="), nil
		}
	case ch == '<':
		return t.newToken(LessThan, string(ch)), nil
	case ch == '>':
		return t.newToken(GreaterThan, string(ch)), nil
	case t.isReserved():
		for _, v := range reserved {
			if t.input[t.pos:t.pos+len(v)] == v {
				return t.newToken(reservedToKind[t.input[t.pos:t.pos+len(v)]], t.input[t.pos:t.pos+len(v)]), nil
			}
		}
	case isDigit(ch):
		return t.lexNumber(), nil
	case isChar(ch):
		return t.lexIdent(), nil
	}
	return Token{}, &TokenizeErr{ErrSyntax, Loc{t.pos, t.pos}, t}
}

// Kind express the token kind as enum
type Kind int

// Define Token Kind as enum
const (
	Num         Kind = iota // 0 - 9
	Plus                    // +
	Minus                   // -
	Asterisk                // *
	Slash                   // /
	Lbracket                // [
	Rbracket                // ]
	LParen                  // (
	RParen                  // )
	Assign                  // =
	Comma                   // ,
	Lbrace                  // {
	Rbrace                  // }
	Eq                      // ==
	NEq                     // !=
	LessThan                // <
	GreaterThan             // >
	Identifier
	EOF
	KeyIf
	KeyDo
	KeyThen
	KeyEnd
	KeyLoop
	KeyWhile
)

var reserved = []string{
	"loop",
	"if",
	"do",
	"then",
	"end",
	"while",
}

var reservedToKind = map[string]Kind{
	"loop":  KeyLoop,
	"if":    KeyIf,
	"do":    KeyDo,
	"then":  KeyThen,
	"end":   KeyEnd,
	"while": KeyWhile,
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
	Loc     Loc
}

// Loc express Line of code for token
type Loc struct {
	Start int
	End   int
}

func (t *Tokenizer) newToken(tk Kind, lit string) Token {
	start := t.pos
	t.pos += len(lit)
	return Token{tk, lit, Loc{start, t.pos}}
}
