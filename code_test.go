package main

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65532}, []byte{byte(OpConstant), 255, 252}},
	}

	for _, tt := range tests {
		instructions := Make(tt.op, tt.operands...)
		for i, b := range tt.expected {
			if instructions[i] != b {
				t.Error("wrong byte code in instruction\n")
			}
		}
	}
}
