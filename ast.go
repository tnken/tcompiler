package main

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

func (p *Parser) fnCallExpr(ident Ident) FnCallExpr {
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
	return FnCallExpr{args: args, ident: ident}
}

type VarDecl struct {
	left  Ident
	right Expr
}

func (vd VarDecl) string() string {
	return vd.left.string() + " = " + vd.right.string()
}

func (vd VarDecl) nodeStmt() {}

func (p *Parser) varDecl(lhd Ident) VarDecl {
	p.nextToken()
	return VarDecl{lhd, p.expr(lowest)}
}

type LoopStmt struct {
	block []Stmt
}

func (ls LoopStmt) string() string {
	s := "loop {"
	for _, b := range ls.block {
		s += " " + b.string()
	}
	return s + " }"
}

func (ls LoopStmt) nodeStmt() {}

type ExprStmt struct {
	val Expr
}

func (es ExprStmt) string() string {
	return es.val.string()
}

func (es ExprStmt) nodeStmt() {}

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
