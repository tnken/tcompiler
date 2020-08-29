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
	opAdd OpKind = iota
	opSub
	opMul
	opDiv
)

// InfixExpr has a operand and two nodes.
type InfixExpr struct {
	tok   token.Token
	op    OpKind
	left  Node
	right Node
}

func (i InfixExpr) string() string {
	fmt.Println(i.tok)
	return "(" + i.left.string() + " " + i.tok.Literal + " " + i.right.string() + ")"
}

// IntegerLiteral express unsigned number
type IntegerLiteral struct {
	tok token.Token
	val int
}

func (i IntegerLiteral) string() string {
	return strconv.Itoa(i.val)
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
