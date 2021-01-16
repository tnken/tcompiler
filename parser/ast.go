package parser

import (
	"strconv"

	"github.com/takeru56/tcompiler/token"
)

//
// Interface
//

// Node abstract Stmt and Expr
type Node interface {
	string() string
}

// Expr abstructs expression
type Expr interface {
	Node
	nodeExpr()
}

// Stmt abstructs some kinds of statements
type Stmt interface {
	Node
	nodeStmt()
}

// For Debugging
func (i InfixExpr) nodeExpr()      {}
func (i IntegerLiteral) nodeExpr() {}
func (i IdentExpr) nodeExpr()      {}
func (c CallExpr) nodeExpr()       {}
func (l LoopStmt) nodeStmt()       {}
func (a AssignStmt) nodeStmt()     {}
func (b BlockStmt) nodeStmt()      {}
func (i IfStmt) nodeStmt()         {}
func (w WhileStmt) nodeStmt()      {}
func (f FunctionDef) nodeStmt()    {}

//
// Expr
//

// OpKind express kind of operands as enum
type OpKind int

const (
	Add OpKind = iota
	Sub
	Mul
	Div
	EQ
	NEQ
	Less
	Greater
)

// InfixExpr has a operand and two nodes.
type InfixExpr struct {
	tok   token.Token
	Op    OpKind
	Left  Node
	Right Node
}

func (i InfixExpr) string() string {
	return "(" + i.Left.string() + " " + i.tok.Literal + " " + i.Right.string() + ")"
}

// IntegerLiteral express unsigned number
type IntegerLiteral struct {
	Tok token.Token
	Val int
}

func (i IntegerLiteral) string() string {
	return strconv.Itoa(i.Val)
}

// IdentKind show kind of the Identifier as enum
type IdentKind int

const (
	variable IdentKind = iota
	fn
)

// IdentExpr has kind and name
type IdentExpr struct {
	kind IdentKind
	Name string
}

func (i IdentExpr) string() string {
	return i.Name
}

type CallExpr struct {
	Ident IdentExpr
}

func (c CallExpr) string() string {
	return c.Ident.Name + "()"
}

//
// Stmt
//

type AssignStmt struct {
	Ident IdentExpr
	Expr  Node
}

func (a AssignStmt) string() string {
	return a.Ident.string() + " = " + a.Expr.string()
}

type BlockStmt struct {
	Nodes []Node
}

func (b BlockStmt) string() string {
	s := "do\n"
	for _, stmt := range b.Nodes {
		s += "  " + stmt.string() + "\n"
	}
	return s + "end"
}

type IfStmt struct {
	Block     BlockStmt
	Condition Node
}

func (i IfStmt) string() string {
	s := "if " + i.Condition.string() + " then\n"
	for _, node := range i.Block.Nodes {
		s += "  " + node.string() + "\n"
	}
	return s + "end"
}

type WhileStmt struct {
	Block     BlockStmt
	Condition Node
}

func (w WhileStmt) string() string {
	s := "while " + w.Condition.string() + " do\n"
	for _, node := range w.Block.Nodes {
		s += "  " + node.string() + "\n"
	}
	return s + "end"
}

// LoopStmt has a block
type LoopStmt struct {
	block []Stmt
}

func (l LoopStmt) string() string {
	s := "loop {"
	for _, b := range l.block {
		s += " " + b.string()
	}
	return s + " }"
}

type FunctionDef struct {
	Ident IdentExpr
	Block BlockStmt
}

func (f FunctionDef) string() string {
	s := "def " + f.Ident.Name + "()\n"
	for _, b := range f.Block.Nodes {
		s += "  " + b.string() + "\n"
	}
	return s + "end"
}
