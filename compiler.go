package main

import (
	"fmt"

	"github.com/takeru56/t/parser"
	"github.com/takeru56/t/table"
)

func emit(op Opcode, operands ...int) {
	ins := Make(op, operands...)
	for _, i := range ins {
		fmt.Printf("%02x", i)
	}
}

type Compiler struct {
	p           []parser.Node
	symbolTable *table.SymbolTable
}

// Compile generates bytecode
func Compile(program []parser.Node) {
	g := &Compiler{program, table.NewSymbolTable()}
	for _, node := range program {
		g.gen(node)
	}
	emit(OpDone, []int{}...)
}

func (c *Compiler) gen(node parser.Node) {
	switch node := node.(type) {
	case parser.IntegerLiteral:
		emit(OpConstant, []int{node.Val}...)
	case parser.InfixExpr:
		c.gen(node.Left)
		c.gen(node.Right)
		switch node.Op {
		case parser.Add:
			emit(OpAdd, []int{}...)
		case parser.Sub:
			emit(OpSub, []int{}...)
		case parser.Mul:
			emit(OpMul, []int{}...)
		case parser.Div:
			emit(OpDiv, []int{}...)
		case parser.EQ:
			emit(OpEQ, []int{}...)
		case parser.NEQ:
			emit(OpNEQ, []int{}...)
		case parser.Less:
			emit(OpLess, []int{}...)
		case parser.Greater:
			emit(OpGreater, []int{}...)
		}
	case parser.IdentExpr:
		symbol, ok := c.symbolTable.Resolve(node.Name)
		if ok {
			emit(OpLoadGlobal, []int{symbol.Index}...)
		}
		// TODO: do error handling, when ok is false
	case parser.AssignStmt:
		c.gen(node.Expr)
		global := c.symbolTable.Define(node.Ident.Name)
		emit(OpStoreGlobal, []int{global.Index}...)
	}
}
