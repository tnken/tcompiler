package main

import (
	"fmt"

	"github.com/takeru56/t/parser"
)

func emit(op Opcode, operands ...int) {
	ins := Make(op, operands...)
	for _, i := range ins {
		fmt.Printf("%02x", i)
	}
}

type Gen struct {
	p []parser.Node
}

// GenProgram generates bytecode
func Program(program []parser.Node) {
	g := &Gen{program}
	for _, node := range program {
		g.gen(node)
	}
	emit(OpDone, []int{}...)
}

func (g *Gen) gen(node parser.Node) {
	switch node := node.(type) {
	case parser.IntegerLiteral:
		emit(OpConstant, []int{node.Val}...)
	case parser.InfixExpr:
		g.gen(node.Left)
		g.gen(node.Right)
		switch node.Op {
		case parser.Add:
			emit(OpAdd, []int{}...)
		case parser.Sub:
			emit(OpSub, []int{}...)
		case parser.Mul:
			emit(OpMul, []int{}...)
		case parser.Div:
			emit(OpDiv, []int{}...)
		}
	}
}
