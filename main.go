package main

import (
	"bufio"
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
	case ch == '[':
		return t.newToken(Lbracket, string(ch))
	case ch == ']':
		return t.newToken(Rbracket, string(ch))
	case ch == ',':
		return t.newToken(Comma, string(ch))
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
	Lbracket // [
	Rbracket // ]
	Comma    // ,
)

const (
	Lowest = iota
	Sum    // + -
	Mult   // * /
	Index
)

var precedences = map[TokenKind]int{
	Plus:     Sum,
	Minus:    Sum,
	Asterisk: Mult,
	Slash:    Mult,
	Lbracket: Index,
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

type ArrayInit struct {
	exprs []Expr
}

func (ai ArrayInit) string() string {
	s := "["
	for i, v := range ai.exprs {
		if i == 0 {
			s += v.string()
		} else {
			s += " " + v.string()
		}
	}
	return s + "]"
}

// Pratt Parsing
// https://github.sfpgmr.net/tdop.github.io/
// https://dev.to/jrop/pratt-parsing
func (p *Parser) expr(precedence int) Expr {
	var lhd Expr
	// Prefix
	switch p.curToken.Kind {
	case Num:
		lhd = NumberLiteral{p.curToken, p.curToken.Literal}
	case Lbracket:
		lhd = p.arrayInit()
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

func (p *Parser) check(expected TokenKind) {
	if p.curToken.Kind != expected {
		panic("error: unexpected Token")
	}
	p.nextToken()
}

func (p *Parser) infixExpr(left Expr) InfixExpr {
	exp := InfixExpr{tok: p.curToken, op: p.curToken.Kind, left: left}

	precedence := p.curToken.precedence()
	p.nextToken()
	exp.right = p.expr(precedence)

	return exp
}

func (p *Parser) arrayInit() ArrayInit {
	p.nextToken()
	exprs := []Expr{}
	for p.curToken.Kind != Rbracket {
		exprs = append(exprs, p.expr(p.curToken.precedence()))
		p.nextToken()
		if p.curToken.Kind != Rbracket {
			p.check(Comma)
		}
	}
	return ArrayInit{exprs: exprs}
}

type Object interface {
	stringVal() string
}

type Integer struct {
	value int
}

func (i Integer) stringVal() string { return strconv.Itoa(i.value) }

type Array struct {
	val []Object
}

func (arr Array) stringVal() string {
	s := "["
	for i, v := range arr.val {
		switch ele := v.(type) {
		//TODO: remove redundancy
		case Integer:
			if i == 0 {
				s += strconv.Itoa(ele.value)
			} else {
				s += " " + strconv.Itoa(ele.value)
			}
		case Array:
			if i == 0 {
				s += ele.stringVal()
			} else {
				s += " " + ele.stringVal()
			}
		}
	}
	return s + "]"
}

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
	case ArrayInit:
		val := []Object{}
		for _, ele := range v.exprs {
			val = append(val, eval(ele))
		}
		return Array{val: val}
	}
	return nil
}

func repl() {
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	for stdin.Scan() {
		text := stdin.Text()
		tokenizer := NewTokenizer(text)
		p := NewParser(tokenizer)
		exp := p.expr(Lowest)
		fmt.Println(eval(exp).stringVal())
		fmt.Print(">> ")
	}
}

func main() {
	repl()
}
