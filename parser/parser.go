package parser

import (
	"strconv"

	"github.com/takeru56/t/token"
)

// Parser has the information of curToken and peekToken
type Parser struct {
	tokenizer *token.Tokenizer
	curToken  token.Token
	peekToken token.Token
}

// New initialize a Parser and returns its pointer
func New(t *token.Tokenizer) *Parser {
	p := &Parser{tokenizer: t}
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances forward curToken in the Parser
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.tokenizer.Next()
}

func (p *Parser) consume(s string) bool {
	if p.curToken.Literal == s {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) Program() []Node {
	program := []Node{}
	for p.curToken.Kind != token.EOF {
		program = append(program, p.stmt())
	}
	return program
}

func (p *Parser) stmt() Node {
	return p.assign()
}

func (p *Parser) assign() Node {
	node := p.expr()
	switch node.(type) {
	case IdentExpr:
		if p.consume("=") {
			return AssignStmt{node, p.expr()}
		}
	}
	return node
}

func (p *Parser) expr() Node {
	node := p.add()
	return node
}

func (p *Parser) add() Node {
	node := p.mul()
	tok := p.curToken
	for {
		if p.consume("+") {
			node = InfixExpr{tok, Add, node, p.mul()}
		} else if p.consume("-") {
			node = InfixExpr{tok, Sub, node, p.mul()}
		} else {
			return node
		}
	}
}

func (p *Parser) mul() Node {
	node := p.prim()
	tok := p.curToken
	for {
		if p.consume("*") {
			node = InfixExpr{tok, Mul, node, p.prim()}
		} else if p.consume("/") {
			node = InfixExpr{tok, Div, node, p.prim()}
		} else {
			return node
		}
	}
}

// prim ::= atom |
func (p *Parser) prim() Node {
	return p.atom()
}

// atom ::= IntegerLiteral | Identifier
func (p *Parser) atom() Node {
	switch p.curToken.Kind {
	case token.Num:
		return p.newIntegerLiteral()
	}
	return p.newIdentifier()
}

func (p *Parser) newIntegerLiteral() Node {
	val, _ := strconv.Atoi(p.curToken.Literal)
	node := IntegerLiteral{p.curToken, val}
	p.nextToken()
	return node
}

func (p *Parser) newIdentifier() Node {
	node := IdentExpr{variable, p.curToken.Literal}
	p.nextToken()
	return node
}
