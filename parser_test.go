package main

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"1", "1"},
		{"1+2*3", "(1 + (2 * 3))"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
		{"a=1+1", "a = (1 + 1)"},
	}

	for _, c := range cases {
		tokenizer := NewToken(c.input)
		p := NewParser(tokenizer)
		stmt := p.stmt()
		fmt.Println("actual: " + stmt.string() + ", expected: " + c.expected)
		if stmt.string() != c.expected {
			fmt.Println(stmt.string())
			t.Error("The ast is wrong\n")
		}
	}
}
