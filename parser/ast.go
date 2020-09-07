package parser

import (
	"fmt"
	"strconv"

	"github.com/takeru56/t/token"
)

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
// Ignore
func (i InfixExpr) nodeExpr()      {}
func (i IntegerLiteral) nodeExpr() {}
func (i IdentExpr) nodeExpr()      {}
func (l LoopStmt) nodeStmt()       {}

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
)

// InfixExpr has a operand and two nodes.
type InfixExpr struct {
	tok   token.Token
	Op    OpKind
	Left  Node
	Right Node
}

func (i InfixExpr) string() string {
	fmt.Println(i.tok)
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
	name string
}

func (i IdentExpr) string() string {
	return i.name
}

//
// Stmt
//

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
