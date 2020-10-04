package compiler

import (
	"fmt"

	"github.com/takeru56/t/parser"
)

func (c *Compiler) emit(op Opcode, operands ...int) {
	ins := Make(op, operands...)
	for _, i := range ins {
		c.instructions = append(c.instructions, i)
	}
}

type Compiler struct {
	p            []parser.Node
	instructions []byte
	symbolTable  *SymbolTable
}

// Compile generates bytecode
func Exec(program []parser.Node) *Compiler {
	c := &Compiler{program, []byte{}, NewSymbolTable()}
	for _, node := range program {
		c.gen(node)
	}
	c.emit(OpDone, []int{}...)
	return c
}

func (c *Compiler) gen(n parser.Node) {
	switch node := n.(type) {
	case parser.IntegerLiteral:
		c.emit(OpConstant, []int{node.Val}...)
	case parser.InfixExpr:
		c.gen(node.Left)
		c.gen(node.Right)
		switch node.Op {
		case parser.Add:
			c.emit(OpAdd, []int{}...)
		case parser.Sub:
			c.emit(OpSub, []int{}...)
		case parser.Mul:
			c.emit(OpMul, []int{}...)
		case parser.Div:
			c.emit(OpDiv, []int{}...)
		case parser.EQ:
			c.emit(OpEQ, []int{}...)
		case parser.NEQ:
			c.emit(OpNEQ, []int{}...)
		case parser.Less:
			c.emit(OpLess, []int{}...)
		case parser.Greater:
			c.emit(OpGreater, []int{}...)
		}
	case parser.IdentExpr:
		symbol, ok := c.symbolTable.Resolve(node.Name)
		if ok {
			c.emit(OpLoadGlobal, []int{symbol.Index}...)
		}
		// TODO: do error handling, when ok is false
	case parser.AssignStmt:
		c.gen(node.Expr)
		symbol, ok := c.symbolTable.Resolve(node.Ident.Name)
		if ok {
			c.emit(OpStoreGlobal, []int{symbol.Index}...)
			return
		}
		global := c.symbolTable.Define(node.Ident.Name)
		c.emit(OpStoreGlobal, []int{global.Index}...)
	case parser.IfStmt:
		c.gen(node.Condition)
		c.emit(OpJNT, []int{0}...)
		blockHead := len(c.instructions)
		ifHead := blockHead - 3
		for _, stmt := range node.Block.Nodes {
			c.gen(stmt)
		}
		ins := Make(OpJNT, []int{len(c.instructions)}...)

		c.instructions[ifHead+1] = ins[1]
		c.instructions[ifHead+2] = ins[2]
	case parser.WhileStmt:
		head := len(c.instructions)
		c.gen(node.Condition)
		c.emit(OpJNT, []int{0}...)
		blockHead := len(c.instructions)
		whileHead := blockHead - 3
		for _, stmt := range node.Block.Nodes {
			c.gen(stmt)
		}
		c.emit(OpJMP, []int{head}...)

		ins := Make(OpJNT, []int{len(c.instructions)}...)
		c.instructions[whileHead+1] = ins[1]
		c.instructions[whileHead+2] = ins[2]
	}
}

func (c *Compiler) Output() {
	for _, bytecode := range c.instructions {
		fmt.Printf("%02x", bytecode)
	}
}
