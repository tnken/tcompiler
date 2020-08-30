package main

import (
	"fmt"

	"github.com/takeru56/t/parser"
)

func emit(op Opcode, operands ...int) {
	ins := Make(op, operands...)
	for _, i := range ins {
		fmt.Printf("%02x\n", i)
	}
}

type BCGen struct {
	p []parser.Node
}

// GenProgram generates bytecode
func GenProgram(program []parser.Node) {
	bcg := &BCGen{program}
	for _, node := range program {
		bcg.gen(node)
	}
}

func (g *BCGen) gen(node parser.Node) {
	switch n := node.(type) {
	case parser.IntegerLiteral:
		emit(OpConstant, []int{n.Val}...)
	}
}
