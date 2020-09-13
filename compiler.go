package main

import (
	"fmt"
)

func emit(op Opcode, operands ...int) {
	ins := Make(op, operands...)
	for _, i := range ins {
		fmt.Printf("%02x", i)
	}
}

type Compiler struct {
	p           []Node
	symbolTable *SymbolTable
}

// Compile generates bytecode
func Compile(program []Node) {
	g := &Compiler{program, NewSymbolTable()}
	for _, node := range program {
		g.gen(node)
	}
	emit(OpDone, []int{}...)
}

func (c *Compiler) gen(node Node) {
	switch node := node.(type) {
	case IntegerLiteral:
		emit(OpConstant, []int{node.Val}...)
	case InfixExpr:
		c.gen(node.Left)
		c.gen(node.Right)
		switch node.Op {
		case Add:
			emit(OpAdd, []int{}...)
		case Sub:
			emit(OpSub, []int{}...)
		case Mul:
			emit(OpMul, []int{}...)
		case Div:
			emit(OpDiv, []int{}...)
		case EQ:
			emit(OpEQ, []int{}...)
		case NEQ:
			emit(OpNEQ, []int{}...)
		case Less:
			emit(OpLess, []int{}...)
		case Greater:
			emit(OpGreater, []int{}...)
		}
	case IdentExpr:
		symbol, ok := c.symbolTable.Resolve(node.Name)
		if ok {
			emit(OpLoadGlobal, []int{symbol.Index}...)
		}
		// TODO: do error handling, when ok is false
	case AssignStmt:
		c.gen(node.Expr)
		global := c.symbolTable.Define(node.Ident.Name)
		emit(OpStoreGlobal, []int{global.Index}...)
	}
}
