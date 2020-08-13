package main

// Parser has the information of curToken and peekToken
type Parser struct {
	tokenizer *Tokenizer
	curToken  Token
	peekToken Token
}

func newParser(t *Tokenizer) *Parser {
	p := &Parser{tokenizer: t}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.tokenizer.next()
}

// 最初に単体のstmtを返せるようにして，あとからファイルを導入して配列で返せるようにする
func (p *Parser) stmt() Stmt {
	switch p.curToken.Kind {
	case KeyLoop:
		p.nextToken()
		p.check(Lbrace)
		b := []Stmt{}

		for p.curToken.Kind != Rbrace {
			b = append(b, p.stmt())
			p.nextToken()
		}
		p.check(Rbrace)
		return LoopStmt{block: b}
	}

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
			return p.fnCallExpr(ident)
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
