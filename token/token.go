package token

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Tokenizer has source code and read position
type Tokenizer struct {
	Input string
	Pos   int
}

type TokenizeErr struct {
	Err error
	L   Loc
	t   *Tokenizer
}

// custom error
var (
	ErrSyntax   = errors.New("Syntax error, undefined token")
	ErrConstant = errors.New("constant not support")
)

func (te *TokenizeErr) Error() string {
	st, l := lineNum(te.t.Input, te.L.Start)
	line := displayLine(te.t.Input, te.L.Start)

	switch te.Err {
	case ErrSyntax:
		return fmt.Sprintf("%d:%d: %v\n%v", l, te.L.Start-st+1, te.Err, line)
	case ErrConstant:
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
	line += "^"
	return line
}

// New initialize a Tokenizer and returns its pointer
func New(input string) *Tokenizer {
	return &Tokenizer{input, 0}
}

func (t *Tokenizer) recognizeMany(f func(byte) bool) {
	for t.Pos < len(t.Input) && f(t.Input[t.Pos]) {
		t.Pos++
	}
}

func isChar(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

func isDigit(b byte) bool {
	return (strings.IndexByte("0123456789", b) > -1)
}

func isAlnum(b byte) bool {
	return isChar(b) || isDigit(b)
}

func (t *Tokenizer) lexNumber() Token {
	start := t.Pos
	t.recognizeMany(isDigit)
	return Token{Num, t.Input[start:t.Pos], Loc{start, t.Pos}}
}

func (t *Tokenizer) lexIdent() Token {
	start := t.Pos
	t.recognizeMany(isAlnum)
	return Token{Identifier, t.Input[start:t.Pos], Loc{start, t.Pos}}
}

func (t *Tokenizer) lexSpaces() {
	t.recognizeMany(func(b byte) bool { return (strings.IndexByte(" \n\t", b) > -1) })
}

func (t *Tokenizer) skipLine() {
	t.recognizeMany(func(b byte) bool { return (b != '\n') })
}

// Next returns a Token and move forward current position
func (t *Tokenizer) Next() (Token, error) {
	// TODO: Refactoring from LL:51 to LL:62
	if t.Pos >= len(t.Input) {
		return t.newToken(EOF, ""), nil
	}
	ch := t.Input[t.Pos]

	if ch == ' ' || ch == '\t' || ch == '\n' {
		t.lexSpaces()
	}
	if t.Pos >= len(t.Input) {
		return t.newToken(EOF, ""), nil
	}
	ch = t.Input[t.Pos]

	if ch == '#' {
		t.skipLine()
	}

	if t.Pos >= len(t.Input) {
		return t.newToken(EOF, ""), nil
	}
	ch = t.Input[t.Pos]

	if ch == ' ' || ch == '\t' || ch == '\n' {
		t.lexSpaces()
	}
	if t.Pos >= len(t.Input) {
		return t.newToken(EOF, ""), nil
	}
	ch = t.Input[t.Pos]

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
	case ch == '.':
		return t.newToken(Dot, string(ch)), nil
	case ch == '=':
		if t.Input[t.Pos+1] == '=' {
			return t.newToken(Eq, "=="), nil
		}
		return t.newToken(Assign, string(ch)), nil
	case ch == '{':
		return t.newToken(Lbrace, string(ch)), nil
	case ch == '}':
		return t.newToken(Rbrace, string(ch)), nil
	case ch == '!':
		if t.Input[t.Pos+1] == '=' {
			return t.newToken(NEq, "!="), nil
		}
	case ch == '<':
		return t.newToken(LessThan, string(ch)), nil
	case ch == '>':
		return t.newToken(GreaterThan, string(ch)), nil
	case isDigit(ch):
		head := t.Pos
		tk := t.lexNumber()
		val, err := strconv.Atoi(tk.Literal)
		if err != nil || val < 0 || val >= 65536 {
			return Token{}, &TokenizeErr{ErrConstant, Loc{head, t.Pos}, t}
		}
		return tk, nil
	case t.isReserved():
		for _, v := range reserved {
			blank := len(t.Input) - t.Pos
			if blank < len(v) {
				continue
			}
			if t.Input[t.Pos:t.Pos+len(v)] == v {
				return t.newToken(reservedToKind[v], t.Input[t.Pos:t.Pos+len(v)]), nil
			}
		}
	case isChar(ch):
		return t.lexIdent(), nil
	}
	return Token{}, &TokenizeErr{ErrSyntax, Loc{t.Pos, t.Pos}, t}
}

// Kind express the token kind as enum
type Kind int

// Define Token Kind as enum
const (
	Num         Kind = iota // 0: 0 - 9
	Plus                    // 1: +
	Minus                   // 2: -
	Asterisk                // 3: *
	Slash                   // 4: /
	Lbracket                // 5: [
	Rbracket                // 6: ]
	LParen                  // 7: (
	RParen                  // 8: )
	Assign                  // 9: =
	Comma                   // 10: ,
	Lbrace                  // 11: {
	Rbrace                  // 12: }
	Eq                      // 13: ==
	NEq                     // 14: !=
	LessThan                // 15: <
	GreaterThan             // 16: >
	Identifier              // 17:
	EOF                     // 18:
	KeyIf                   // 19:
	KeyDo                   // 20:
	KeyThen                 // 21:
	KeyEnd                  // 22:
	KeyLoop                 // 23:
	KeyWhile                // 24:
	KeyDef                  // 25:
	KeyReturn               // 26:
	KeyClass                // 27:
	Dot                     // 28: .
	KeySelf                 // 29:
	Number                  // 30: #
)

var reserved = []string{
	"loop",
	"if",
	"do",
	"then",
	"end",
	"while",
	"def",
	"return",
	"class",
	"self",
}

var reservedToKind = map[string]Kind{
	"loop":   KeyLoop,
	"if":     KeyIf,
	"do":     KeyDo,
	"then":   KeyThen,
	"end":    KeyEnd,
	"while":  KeyWhile,
	"def":    KeyDef,
	"return": KeyReturn,
	"class":  KeyClass,
	"self":   KeySelf,
}

func (t Tokenizer) isReserved() bool {
	for _, v := range reserved {
		blank := len(t.Input) - t.Pos
		if blank < len(v) {
			continue
		}
		if t.Input[t.Pos:t.Pos+len(v)] == v {
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
	start := t.Pos
	t.Pos += len(lit)
	return Token{tk, lit, Loc{start, t.Pos}}
}
