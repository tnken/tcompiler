package code

import (
	"encoding/binary"
	"fmt"
)

// Instructions is byte sequence of instruction that consists of Opcode and Operands
type Instructions []byte

type Opcode byte

// Define Opcode
const (
	OpConstant    Opcode = iota // 0
	OpAdd                       // 1
	OpSub                       // 2
	OpMul                       // 3
	OpDiv                       // 4
	OpDone                      // 5
	OpEQ                        // 6
	OpNEQ                       // 7
	OpLess                      // 8
	OpGreater                   // 9
	OpLoadGlobal                // 10
	OpStoreGlobal               // 11
	OpJNT                       // 12
	OpJMP                       // 13
	OpCall                      // 14
)

// Definition consits of Name and OperandWidths property
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:    {"OpConstant", []int{2}},
	OpAdd:         {"OpAdd", []int{}},
	OpSub:         {"OpSub", []int{}},
	OpMul:         {"OpMul", []int{}},
	OpDiv:         {"OpDiv", []int{}},
	OpDone:        {"OpDone", []int{}},
	OpEQ:          {"OpEQ", []int{}},
	OpNEQ:         {"OpNEQ", []int{}},
	OpLess:        {"OpLess", []int{}},
	OpGreater:     {"OpGreater", []int{}},
	OpLoadGlobal:  {"OpLoadGlobal", []int{1}},
	OpStoreGlobal: {"OpStoreGlobal", []int{1}},
	OpJNT:         {"OpJNT", []int{2}}, // false → OpJNTの位置+[]int{2}の分飛ぶ
	OpJMP:         {"OpJMP", []int{2}},
	OpCall:        {"OpCall", []int{1}},
}

// Lookup finds Definition of Opcode
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("Undefined Opcode: %d", op)
	}
	return def, nil
}

// Make convert Opcode and operands to byte
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}
