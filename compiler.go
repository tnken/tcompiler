package main

import "fmt"

func (c *Compiler) emit(op Opcode, operands ...int) {
	ins := Make(op, operands...)
	for _, i := range ins {
		c.instructions = append(c.instructions, i)
	}
}

type Compiler struct {
	p            []Node
	instructions []byte
	symbolTable  *SymbolTable
}

// Compile generates bytecode
func Compile(program []Node) *Compiler {
	c := &Compiler{program, []byte{}, NewSymbolTable()}
	for _, node := range program {
		c.gen(node)
	}
	c.emit(OpDone, []int{}...)
	return c
}

func (c *Compiler) gen(node Node) {
	switch node := node.(type) {
	case IntegerLiteral:
		c.emit(OpConstant, []int{node.Val}...)
	case InfixExpr:
		c.gen(node.Left)
		c.gen(node.Right)
		switch node.Op {
		case Add:
			c.emit(OpAdd, []int{}...)
		case Sub:
			c.emit(OpSub, []int{}...)
		case Mul:
			c.emit(OpMul, []int{}...)
		case Div:
			c.emit(OpDiv, []int{}...)
		case EQ:
			c.emit(OpEQ, []int{}...)
		case NEQ:
			c.emit(OpNEQ, []int{}...)
		case Less:
			c.emit(OpLess, []int{}...)
		case Greater:
			c.emit(OpGreater, []int{}...)
		}
	case IdentExpr:
		symbol, ok := c.symbolTable.Resolve(node.Name)
		if ok {
			c.emit(OpLoadGlobal, []int{symbol.Index}...)
		}
		// TODO: do error handling, when ok is false
	case AssignStmt:
		c.gen(node.Expr)
		symbol, ok := c.symbolTable.Resolve(node.Ident.Name)
		if ok {
			c.emit(OpStoreGlobal, []int{symbol.Index}...)
			return
		}
		global := c.symbolTable.Define(node.Ident.Name)
		c.emit(OpStoreGlobal, []int{global.Index}...)
	case IfStmt:
		c.gen(node.condition)
		c.emit(OpJNT, []int{0}...)
		blockHead := len(c.instructions)
		ifHead := blockHead - 3
		for _, stmt := range node.block.nodes {
			c.gen(stmt)
		}
		ins := Make(OpJNT, []int{len(c.instructions)}...)

		c.instructions[ifHead+1] = ins[1]
		c.instructions[ifHead+2] = ins[2]
	case WhileStmt:
		head := len(c.instructions)
		c.gen(node.condition)
		c.emit(OpJNT, []int{0}...)
		blockHead := len(c.instructions)
		whileHead := blockHead - 3
		for _, stmt := range node.block.nodes {
			c.gen(stmt)
		}
		c.emit(OpJMP, []int{head}...)

		ins := Make(OpJNT, []int{len(c.instructions)}...)
		c.instructions[whileHead+1] = ins[1]
		c.instructions[whileHead+2] = ins[2]
	}
}

func (c *Compiler) output() {
	for _, bytecode := range c.instructions {
		fmt.Printf("%02x", bytecode)
	}
}
