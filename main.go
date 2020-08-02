package main

import (
	"fmt"
	"os"
)

type Tokenizer struct {
	input string
	pos   int
	ch    byte
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input, 0, 0}
}

//func (t *Tokenizer) eatSpaces() {
//}

//func (t *Tokenizer) readNumber() {
//}

func (t *Tokenizer) next() Token {
	if t.pos >= len(t.input) {
		return Token{Eof, ""}
	}

	t.ch = t.input[t.pos]
	t.pos++
	switch {
	case t.ch == '+':
		return t.newToken(Plus, string(t.ch))
	case t.ch == '-':
		return t.newToken(Minus, string(t.ch))
	default:
		return t.newToken(Num, string(t.ch))
	}
}

type TokenKind int

const (
	Num TokenKind = iota
	Plus
	Minus
	Eof
)

const (
	Lowest = iota
	Sum    // + -
)

var precedences = map[TokenKind]int{
	Plus:  Sum,
	Minus: Sum,
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

type NumberLiteral struct {
	tok Token
	val string
}

func (nl NumberLiteral) string() string {
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
		case Plus:
			p.nextToken()
			lhd = p.infixExpr(lhd)
		case Minus:
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

func main() {
	if len(os.Args) < 2 {
		panic("error: argument missing")
	}

	tokenizer := NewTokenizer(os.Args[1])
	p := NewParser(tokenizer)
	exp := p.expr(Lowest)
	fmt.Println(exp.string())
}
