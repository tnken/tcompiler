package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/takeru56/tcompiler/token"
)

// Parser has the information of curToken and peekToken
type Parser struct {
	tokenizer *token.Tokenizer
	curToken  token.Token
	peekToken token.Token
}

type ParseErr struct {
	Err error
	L   token.Loc
	p   *Parser
}

// custom error
var (
	ErrSyntax = errors.New("Syntax error")
)

func (pe *ParseErr) Error() string {
	st, l := lineNum(pe.p.tokenizer.Input, pe.L.Start)
	line := displayLine(pe.p.tokenizer.Input, pe.L.Start)

	switch pe.Err {
	case ErrSyntax:
		return fmt.Sprintf("%d:%d: %v\n%v", l, pe.L.Start-st+1, pe.Err, line)
	}
	return pe.Err.Error()
}

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

// New initialize a Parser and returns its pointer
func New(t *token.Tokenizer) (*Parser, error) {
	p := &Parser{tokenizer: t}
	err := p.nextToken()
	if err != nil {
		return p, err
	}
	err = p.nextToken()
	if err != nil {
		return p, err
	}
	return p, nil
}

// nextToken advances forward curToken in the Parser
func (p *Parser) nextToken() error {
	p.curToken = p.peekToken
	t, err := p.tokenizer.Next()
	if err != nil {
		return err
	}
	p.peekToken = t
	return nil
}

func (p *Parser) consume(s string) (bool, error) {
	if p.curToken.Literal == s {
		err := p.nextToken()
		if err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

// 以下LL(1)parser
// TODO: BNFで可視化

func (p *Parser) Program() ([]Node, error) {
	program := []Node{}
	for p.curToken.Kind != token.EOF {
		n, err := p.class()
		if err != nil {
			return nil, err
		}
		program = append(program, n)
	}
	return program, nil
}

func (p *Parser) class() (Node, error) {
	// parse classDef
	f, err := p.consume("class")
	if err != nil {
		return ClassDef{}, err
	}

	if f {
		ident, ok := p.newFnIdentifier().(IdentExpr)
		if !ok {
			return ClassDef{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
		}

		methods := []FunctionDef{}
		for {
			f, err = p.consume("end")
			if err != nil {
				return ClassDef{}, err
			}
			if f {
				break
			}
			if p.curToken.Kind == token.EOF {
				return ClassDef{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
			}

			node, err := p.function()
			if err != nil {
				return FunctionDef{}, err
			}
			method, ok := node.(FunctionDef)
			if !ok {
				return nil, ErrSyntax
			}
			method.FlagMethod = true
			methods = append(methods, method)
		}

		return ClassDef{ident, methods}, nil
	}

	node, err := p.function()
	if err != nil {
		return ClassDef{}, err
	}
	return node, nil
}

func (p *Parser) function() (Node, error) {
	// parse functionDef
	f, err := p.consume("def")
	if err != nil {
		return FunctionDef{}, err
	}
	if f {
		// ident
		ident, ok := p.newFnIdentifier().(IdentExpr)
		if !ok {
			return FunctionDef{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
		}
		// params
		_, err := p.consume("(")
		if err != nil {
			return FunctionDef{}, err
		}

		args := []IdentExpr{}
		for {
			if p.curToken.Kind == token.EOF {
				return FunctionDef{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
			}
			f, err = p.consume(")")
			if err != nil {
				return FunctionDef{}, err
			}
			if f {
				break
			}
			if len(args) > 0 {
				f, err = p.consume(",")
				if err != nil {
					return FunctionDef{}, err
				}
			}
			arg, ok := p.newFnIdentifier().(IdentExpr)
			if !ok {
				return FunctionDef{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
			}
			args = append(args, arg)
		}

		// block
		block := BlockStmt{Nodes: []Node{}}
		for {
			f, err = p.consume("end")
			if err != nil {
				return FunctionDef{}, err
			}
			if f {
				break
			}
			if p.curToken.Kind == token.EOF {
				return FunctionDef{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
			}

			n, err := p.stmt()
			if err != nil {
				return FunctionDef{}, err
			}
			block.Nodes = append(block.Nodes, n)
		}
		return FunctionDef{ident, block, args, false}, nil
	}
	node, err := p.stmt()
	if err != nil {
		return FunctionDef{}, err
	}
	return node, nil
}

func (p *Parser) stmt() (Node, error) {
	f, err := p.consume("if")
	if err != nil {
		return IfStmt{}, err
	}
	if f {
		block := BlockStmt{Nodes: []Node{}}
		node, err := p.expr()
		if err != nil {
			return node, err
		}

		f, err = p.consume("do")
		if err != nil {
			return IfStmt{}, err
		}

		for {
			f, err = p.consume("end")
			if err != nil {
				return IfStmt{}, err
			}
			if f {
				break
			}
			if p.curToken.Kind == token.EOF {
				return IfStmt{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
			}

			n, err := p.stmt()
			if err != nil {
				return IfStmt{}, err
			}
			block.Nodes = append(block.Nodes, n)
		}
		return IfStmt{Condition: node, Block: block}, nil
	}

	f, err = p.consume("while")
	if err != nil {
		return WhileStmt{}, err
	}
	if f {
		block := BlockStmt{Nodes: []Node{}}
		node, err := p.expr()
		if err != nil {
			return BlockStmt{}, err
		}
		_, err = p.consume("do")
		if err != nil {
			return BlockStmt{}, err
		}

		for {
			f, err = p.consume("end")
			if err != nil {
				return WhileStmt{}, err
			}
			if f {
				break
			}
			if p.curToken.Kind == token.EOF {
				return IfStmt{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
			}

			n, err := p.stmt()
			if err != nil {
				return WhileStmt{}, err
			}
			block.Nodes = append(block.Nodes, n)
		}
		return WhileStmt{Condition: node, Block: block}, nil
	}

	f, err = p.consume("return")
	if err != nil {
		return ReturnStmt{}, err
	}
	if f {
		node, err := p.expr()
		if err != nil {
			return node, err
		}
		return ReturnStmt{node}, nil
	}

	node, err := p.assign()
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) assign() (Node, error) {
	node, err := p.expr()
	if err != nil {
		return node, err
	}
	switch node.(type) {
	case IdentExpr:
		f, err := p.consume("=")
		if err != nil {
			return AssignStmt{}, err
		}
		if f {
			n, err := p.expr()
			if err != nil {
				return AssignStmt{}, err
			}
			return AssignStmt{node.(IdentExpr), n}, nil
		}
	}
	return node, nil
}

func (p *Parser) expr() (Node, error) {
	node, err := p.eq()
	if err != nil {
		return node, err
	}
	return node, nil
}

func (p *Parser) eq() (Node, error) {
	node, err := p.compare()
	if err != nil {
		return node, err
	}
	tok := p.curToken
	for {
		if f, err := p.consume("=="); f {
			if err != nil {
				return node, err
			}
			n, err := p.compare()
			if err != nil {
				return n, err
			}
			node = InfixExpr{tok, EQ, node, n}
		} else if f, err := p.consume("!="); f {
			if err != nil {
				return node, err
			}
			n, err := p.compare()
			if err != nil {
				return n, err
			}
			node = InfixExpr{tok, NEQ, node, n}
		} else {
			return node, nil
		}
	}
}

func (p *Parser) compare() (Node, error) {
	node, err := p.add()
	if err != nil {
		return node, err
	}
	tok := p.curToken
	for {
		if f, err := p.consume("<"); f {
			if err != nil {
				return node, err
			}
			n, err := p.add()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Less, node, n}
		} else if f, err := p.consume(">"); f {
			if err != nil {
				return node, err
			}
			n, err := p.add()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Greater, node, n}
		} else {
			return node, nil
		}
	}
}

func (p *Parser) add() (Node, error) {
	node, err := p.mul()
	if err != nil {
		return node, err
	}
	tok := p.curToken
	for {
		if f, err := p.consume("+"); f {
			if err != nil {
				return node, err
			}
			n, err := p.mul()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Add, node, n}
		} else if f, err := p.consume("-"); f {
			if err != nil {
				return node, err
			}
			n, err := p.mul()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Sub, node, n}
		} else {
			return node, nil
		}
	}
}

func (p *Parser) mul() (Node, error) {
	node, err := p.prim()
	if err != nil {
		return node, err
	}
	tok := p.curToken
	for {
		if f, err := p.consume("*"); f {
			if err != nil {
				return node, err
			}
			n, err := p.prim()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Mul, node, n}
		} else if f, err := p.consume("/"); f {
			if err != nil {
				return node, err
			}
			n, err := p.prim()
			if err != nil {
				return node, err
			}
			node = InfixExpr{tok, Div, node, n}
		} else {
			return node, nil
		}
	}
}

// prim ::= atom |
func (p *Parser) prim() (Node, error) {
	node, err := p.atom()
	return node, err
}

// atom ::= IntegerLiteral | Identifier
func (p *Parser) atom() (Node, error) {
	switch p.curToken.Kind {
	case token.Num:
		if p.peekToken.Kind == token.DotDot {
			return p.newIntegerRangeLiteral(), nil
		}
		return p.newIntegerLiteral(), nil
	case token.KeyTrue:
		return p.newBoolLiteral(), nil
	case token.KeyFalse:
		return p.newBoolLiteral(), nil
	case token.Identifier:
		var n Node
		// CallExpr
		if p.peekToken.Kind == token.LParen {
			literal := p.curToken.Literal
			p.nextToken()
			p.nextToken()
			args := []Node{}
			for {
				if p.curToken.Kind == token.EOF {
					return CallExpr{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
				}
				f, err := p.consume(")")
				if err != nil {
					return CallExpr{}, err
				}
				if f {
					break
				}
				if len(args) > 0 {
					f, err = p.consume(",")
					if err != nil || !f {
						return CallExpr{}, err
					}
				}
				arg, err := p.expr()
				if err != nil {
					return CallExpr{}, &ParseErr{ErrSyntax, p.curToken.Loc, p}
				}
				args = append(args, arg)
			}

			if p.curToken.Kind != token.Dot {
				if 'A' <= literal[0] && literal[0] <= 'Z' {
					n = InstantiationExpr{IdentExpr{variable, literal, false, Any, IntegerRangeLiteral{}}, args}
					return n, nil
				}
				n = CallExpr{IdentExpr{variable, literal, false, Any, IntegerRangeLiteral{}}, args}
				return n, nil
			}
		} else {
			n = p.newValIdentifier(false, Any, IntegerRangeLiteral{})
		}

		// call method
		f, err := p.consume(".")
		if err != nil {
			return CallExpr{}, err
		}
		if f {
			node, err := p.atom()
			if err != nil {
				return CallMethodExpr{}, err
			}
			return CallMethodExpr{n, node}, err
		}
		return n, nil
	case token.KeySelf:
		p.nextToken()
		_, err := p.consume(".")
		if err != nil {
			return IdentExpr{}, err
		}
		n, _ := p.newValIdentifier(true, Any, IntegerRangeLiteral{}).(IdentExpr)

		f, err := p.consume(":")
		if err != nil {
			return IdentExpr{}, err
		}
		if f {
			// instance val checker
			switch p.curToken.Kind {
			case token.KeyNumber:
				n.ValType = Num
				p.nextToken()
				return n, nil
			case token.KeyBool:
				n.ValType = Bool
				p.nextToken()
				return n, nil

			case token.Lbrace:
				p.nextToken()
				if p.curToken.Kind == token.KeyInclude {
					p.nextToken()
					// デモ用に動作固定
					// TODO: fix
					_, err := p.consume(":")
					if err != nil {
						return IdentExpr{}, err
					}
					lim, _ := p.newIntegerRangeLiteral().(IntegerRangeLiteral)
					_, err = p.consume("}")
					if err != nil {
						return IdentExpr{}, err
					}
					n.ValLimit = lim
					n.ValType = Include
					return n, nil
				}
				if p.curToken.Kind == token.KeyExclude {
					p.nextToken()
					_, err = p.consume(":")
					if err != nil {
						return IdentExpr{}, err
					}
					lim, _ := p.newIntegerRangeLiteral().(IntegerRangeLiteral)
					_, err = p.consume("}")
					if err != nil {
						return IdentExpr{}, err
					}
					n.ValLimit = lim
					n.ValType = Exclude
					return n, nil
				}
			}
		}

		return n, nil
	}
	return p.newValIdentifier(false, Any, IntegerRangeLiteral{}), nil
}

func (p *Parser) newIntegerLiteral() Node {
	val, _ := strconv.Atoi(p.curToken.Literal)
	node := IntegerLiteral{p.curToken, val}
	p.nextToken()
	return node
}

func (p *Parser) newBoolLiteral() Node {
	node := BoolLiteral{p.curToken}
	p.nextToken()
	return node
}

func (p *Parser) newIntegerRangeLiteral() Node {
	val, _ := strconv.Atoi(p.curToken.Literal)
	from := IntegerLiteral{p.curToken, val}
	p.nextToken()
	p.nextToken()
	val, _ = strconv.Atoi(p.curToken.Literal)
	to := IntegerLiteral{p.curToken, val}
	p.nextToken()
	return IntegerRangeLiteral{From: from, To: to}
}

func (p *Parser) newValIdentifier(flag bool, vt IdentValType, lim IntegerRangeLiteral) Node {
	node := IdentExpr{variable, p.curToken.Literal, flag, vt, lim}
	p.nextToken()
	return node
}

func (p *Parser) newFnIdentifier() Node {
	node := IdentExpr{fn, p.curToken.Literal, false, Any, IntegerRangeLiteral{}}
	p.nextToken()
	return node
}
