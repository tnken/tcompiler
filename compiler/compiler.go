package compiler

import (
	"encoding/binary"
	"fmt"

	"github.com/takeru56/tcompiler/obj"
	"github.com/takeru56/tcompiler/parser"
)

func (c *Compiler) emit(op Opcode, operands ...int) {
	ins := Make(op, operands...)
	for _, i := range ins {
		c.instructions = append(c.instructions, i)
	}
}

type Compiler struct {
	p            []parser.Node
	constantPool []obj.Object
	instructions []byte
	symbolTable  *SymbolTable
}

func Exec(program []parser.Node) *Compiler {
	c := &Compiler{program, []obj.Object{}, []byte{}, NewSymbolTable()}
	for _, node := range program {
		c.gen(node)
	}
	c.emit(OpDone, []int{}...)
	return c
}

func (c *Compiler) addConstant(obj obj.Object) int {
	c.constantPool = append(c.constantPool, obj)
	return len(c.constantPool)
}

func (c *Compiler) gen(n parser.Node) {
	switch node := n.(type) {
	case parser.IntegerLiteral:
		integer := &obj.Integer{Value: node.Val}
		c.emit(OpConstant, []int{c.addConstant(integer)}...)
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

type ConstantType byte

// Define Opcode
const (
	CONST_INT ConstantType = iota
)

func toUint16(num int) [2]byte {
	b := [2]byte{}
	binary.BigEndian.PutUint16(b[0:], uint16(num))
	return b
}

// output tarto IR bytecode Format
// ***************************************

// struct {
// 	u4 magic
// 	u2 constant_pool_count
// 	cp constant_pool[constant_pool_count]
// 	u2 instruction_count
// 	ins instructions[instruction_count]
// }

// struct constant_pool {
// 	u1 constant type
// 	u2 constant size
// 	c [const size]constants
// }
// ***************************************

// TODO: 32bitに拡張+エラー処理

func (c *Compiler) Bytecode() string {
	b := ""
	// u4 magic（特に意味無し）
	b += fmt.Sprintf("%02x", []byte{255, 255, 255, 255})
	// u2 constant_pool_count
	b += fmt.Sprintf("%02x", toUint16(len(c.constantPool)))
	// const pool
	for _, constant := range c.constantPool {
		switch constant := constant.(type) {
		case *obj.Integer:
			// u1
			b += fmt.Sprintf("%02x", CONST_INT)
			// u2
			b += fmt.Sprintf("%02x", toUint16(constant.Size()))
			// u2
			b += fmt.Sprintf("%02x", toUint16(constant.Value))
		}
	}
	// u2 instruction_count
	b += fmt.Sprintf("%02x", toUint16(len(c.instructions)))

	// instruction
	for _, bytecode := range c.instructions {
		b += fmt.Sprintf("%02x", bytecode)
	}
	return b
}

func (c *Compiler) Output() {
	fmt.Print(c.Bytecode())
}

func (c *Compiler) Dump() {
	b := c.Bytecode()
	p := 0
	size := 0
	for p < len(b) {
		if size%16 == 0 {
			if size != 0 {
				fmt.Print("\n")
			}
			fmt.Printf("%02X: ", size)
		}
		if size%16 != 0 && size%8 == 0 {
			fmt.Print(" ")
		}
		fmt.Print(b[p : p+2])
		p += 2
		size++
	}
}
