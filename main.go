package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Tokenizer
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
	case ch == '(':
		return t.newToken(LParen, string(ch))
	case ch == ')':
		return t.newToken(RParen, string(ch))
	case ch == ',':
		return t.newToken(Comma, string(ch))
	case ch == '=':
		return t.newToken(Assign, string(ch))
	case isDigit(ch):
		return t.lexNumber()
	case isChar((ch)):
		return t.lexIdent()
	}
	return t.newToken(Eof, "")
}

// TokenKind express the kind of the token
type TokenKind int

const (
	Num      TokenKind = iota //0-9
	Plus                      // +
	Minus                     // -
	Asterisk                  // *
	Slash                     // /
	Lbracket                  // [
	Rbracket                  // ]
	LParen                    // (
	RParen                    // )
	Assign                    // =
	Comma                     // ,
	Identifier
	Eof
)

const (
	lowest = iota
	assign
	sum
	mult
	index
)

var precedences = map[TokenKind]int{
	Plus:     sum,
	Minus:    sum,
	Asterisk: mult,
	Slash:    mult,
	Lbracket: index,
	Assign:   assign,
}

type Token struct {
	Kind    TokenKind
	Literal string
}

func (tok Token) precedence() int {
	if precedence, ok := precedences[tok.Kind]; ok {
		return precedence
	}
	return lowest
}

func (t *Tokenizer) newToken(k TokenKind, lit string) Token {
	t.pos++
	return Token{k, lit}
}

// parse
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

// ast
// TODO: fix methods of interface, stmt, expr
type Node interface {
	string() string
}
type Expr interface {
	Node
	nodeExpr()
}

type Stmt interface {
	Node
	nodeStmt()
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

func (ie InfixExpr) nodeExpr() {}

type NumberLiteral struct {
	tok Token
	val string
}

func (nl NumberLiteral) string() string {
	return nl.val
}

func (nl NumberLiteral) nodeExpr() {}

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

func (ai ArrayInit) nodeExpr() {}

type IdentKind int

const (
	variable IdentKind = iota
	fn
)

type Ident struct {
	kind IdentKind
	name string
}

func (i Ident) string() string {
	return i.name
}

func (i Ident) nodeExpr() {}

type FnCallExpr struct {
	ident Ident
	args  []Expr
}

func (fc FnCallExpr) string() string {
	args := ""
	for i, arg := range fc.args {
		if i == 0 {
			args += arg.string()
		} else {
			args += " " + arg.string()
		}
	}
	return fc.ident.name + "(" + args + ")"
}

func (fc FnCallExpr) nodeExpr() {}

type VarDecl struct {
	left  Ident
	right Expr
}

func (vd VarDecl) string() string {
	return vd.left.string() + " = " + vd.right.string()
}

func (vd VarDecl) nodeStmt() {}

type ExprStmt struct {
	val Expr
}

func (es ExprStmt) string() string {
	return es.val.string()
}

func (es ExprStmt) nodeStmt() {}

// 最初に単体のstmtを返せるようにして，あとからファイルを導入して配列で返せるようにする
func (p *Parser) stmt() Stmt {
	lhd := p.expr(lowest)
	switch v := lhd.(type) {
	case Ident:
		if p.peekToken.Kind == Assign {
			p.nextToken()
			return p.varDecl(v)
		} else {
			return ExprStmt{val: v}
		}
	default:
		return ExprStmt{val: v}
	}
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
	case Identifier:
		if p.peekToken.Kind == LParen {
			ident := Ident{fn, p.curToken.Literal}
			p.nextToken()
			p.check(LParen)
			return p.FnCallExpr(ident)
		} else {
			lhd = Ident{variable, p.curToken.Literal}
		}
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

func (p *Parser) varDecl(lhd Ident) VarDecl {
	p.nextToken()
	return VarDecl{lhd, p.expr(lowest)}
}

func (p *Parser) FnCallExpr(ident Ident) FnCallExpr {
	args := []Expr{}
	i := 0
	for p.curToken.Kind != RParen {
		args = append(args, p.expr(lowest))
		i++
		p.nextToken()
		if i > 6 {
			panic("error: too many argument")
		}
	}
	p.check(RParen)
	return FnCallExpr{args: args, ident: ident}
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

// eval
type Var struct {
	name string
	obj  Object
}

func (v Var) stringVal() string { return v.obj.stringVal() }

type Eval struct {
	vars map[string]Var
}

func newEval() Eval {
	return Eval{vars: map[string]Var{}}
}

func (e Eval) eval(node Node) Object {
	return e.stmt(node.(Stmt))
}

func (e Eval) stmt(stmt Stmt) Object {
	switch s := stmt.(type) {
	case VarDecl:
		name := s.left.name
		v := Var{name: name, obj: e.expr(s.right)}
		e.vars[name] = v
		return v
	case ExprStmt:
		return e.expr(s.val)
	}
	panic("error")
}

// Tree Walk
func (e Eval) expr(expr Expr) Object {
	switch v := expr.(type) {
	case InfixExpr:
		l := e.expr(v.left).(Integer).value
		r := e.expr(v.right).(Integer).value
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
			val = append(val, e.expr(ele))
		}
		return Array{val: val}
	case Ident:
		return e.vars[v.name].obj
	}
	return nil
}

func repl() {
	stdin := bufio.NewScanner(os.Stdin)
	e := newEval()
	fmt.Print(">> ")
	for stdin.Scan() {
		text := stdin.Text()
		tokenizer := newTokenizer(text)
		p := NewParser(tokenizer)
		stmt := p.stmt()
		fmt.Println(e.eval(stmt).stringVal())
		fmt.Print(">> ")
	}
}

func main() {
	repl()
}
