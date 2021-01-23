package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65532}, []byte{byte(OpConstant), 255, 252}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpMul, []int{}, []byte{byte(OpMul)}},
		{OpDiv, []int{}, []byte{byte(OpDiv)}},
		{OpDone, []int{}, []byte{byte(OpDone)}},
		{OpStoreGlobal, []int{0}, []byte{byte(OpStoreGlobal), 0}},
		{OpCallMethod, []int{3}, []byte{byte(OpCallMethod), 3}},
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
