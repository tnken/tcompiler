package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Tokenizer struct {
	input string
	pos   int
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input, 0}
}

func (t *Tokenizer) recognizeMany(f func(byte) bool) {
	for t.pos < len(t.input) && f(t.input[t.pos]) {
		t.pos++
	}
}

func isDigit(b byte) bool {
	return (strings.IndexByte("0123456789", b) > -1)
}

func (t *Tokenizer) lexNumber() Token {
	start := t.pos
	t.recognizeMany(isDigit)
	return Token{Num, t.input[start:t.pos]}
}

func (t *Tokenizer) skipSpaces() {
	t.recognizeMany(func(b byte) bool { return (strings.IndexByte(" \n\t", b) > -1) })
}

func (t *Tokenizer) next() Token {
	// TODO: more simple
	if t.pos >= len(t.input) {
		return t.newToken(Eof, "")
	}
	ch := t.input[t.pos]
	if ch == ' ' || ch == '\t' || ch == '\n' {
		t.skipSpaces()
	}

	if t.pos >= len(t.input) {
		return t.newToken(Eof, "")
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
	case isDigit(ch):
		return t.lexNumber()
	}
	return t.newToken(Eof, "")
}

type TokenKind int

const (
	Num TokenKind = iota
	Plus
	Minus
	Asterisk
	Slash
	Eof
)

const (
	Lowest = iota
	Sum    // + -
	Mult   // * /
)

var precedences = map[TokenKind]int{
	Plus:     Sum,
	Minus:    Sum,
	Asterisk: Mult,
	Slash:    Mult,
}

type Token struct {
	Kind    TokenKind
	Literal string
}

func (tok Token) precedence() int {
	if precedence, ok := precedences[tok.Kind]; ok {
		return precedence
	}
	return Lowest
}

func (t *Tokenizer) newToken(k TokenKind, lit string) Token {
	t.pos++
	return Token{k, lit}
}

type Parser struct {
	tokenizer *Tokenizer
	curToken  Token
	peekToken Token
}

func NewParser(tok *Tokenizer) *Parser {
	p := &Parser{tokenizer: tok}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.tokenizer.next()
}

type Expr interface {
	string() string
	literal() string
}

type InfixExpr struct {
	tok   Token
	op    TokenKind
	left  Expr
	right Expr
}

func (ie InfixExpr) string() string {
	return "(" + ie.left.string() + " " + ie.tok.Literal + " " + ie.right.string() + ")"
}

func (ie InfixExpr) literal() string {
	return ie.tok.Literal
}

type NumberLiteral struct {
	tok Token
	val string
}

func (nl NumberLiteral) string() string {
	return nl.val
}

func (nl NumberLiteral) literal() string {
	return nl.val
}

// Pratt Parsing
// https://github.sfpgmr.net/tdop.github.io/
// https://dev.to/jrop/pratt-parsing
func (p *Parser) expr(precedence int) Expr {
	var lhd Expr
	// Prefix
	switch p.curToken.Kind {
	case Num:
		lhd = p.numberLiteral()
	default:
		return nil
	}

	for precedence < p.peekToken.precedence() {
		// Infix
		switch p.peekToken.Kind {
		case Plus, Minus, Asterisk, Slash:
			p.nextToken()
			lhd = p.infixExpr(lhd)
		default:
			return lhd
		}
	}

	return lhd
}

func (p *Parser) infixExpr(left Expr) InfixExpr {
	exp := InfixExpr{tok: p.curToken, op: p.curToken.Kind, left: left}

	precedence := p.curToken.precedence()
	p.nextToken()
	exp.right = p.expr(precedence)

	return exp
}

func (p *Parser) numberLiteral() NumberLiteral {
	return NumberLiteral{p.curToken, p.curToken.Literal}
}

type Object interface {
	stringVal() string
}

type Integer struct {
	value int
}

func (i Integer) stringVal() string { return strconv.Itoa(i.value) }

// Tree Walk
func eval(expr Expr) Object {
	switch v := expr.(type) {
	case InfixExpr:
		l := eval(v.left).(Integer).value
		r := eval(v.right).(Integer).value
		switch v.op {
		case Plus:
			return Integer{value: l + r}
		case Minus:
			return Integer{value: l - r}
		case Asterisk:
			return Integer{value: l * r}
		case Slash:
			return Integer{value: int(l / r)}
		}
	case NumberLiteral:
		i, _ := strconv.Atoi(v.val)
		return Integer{value: i}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		panic("error: argument missing")
	}

	tokenizer := NewTokenizer(os.Args[1])
	p := NewParser(tokenizer)
	exp := p.expr(Lowest)
	fmt.Println(exp.string())
	fmt.Println(eval(exp).stringVal())
}
