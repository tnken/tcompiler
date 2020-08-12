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
		{"[1+1]", "[(1 + 1)]"},
		{"[1+2*3, 2, [1, 2]]", "[(1 + (2 * 3)) 2 [1 2]]"},
		{"a = 1", "a = 1"},
		{"testfn(1+2*3)", "testfn((1 + (2 * 3)))"},
		{"loop { print(3) }", "loop { print(3) }"},
	}

	for _, c := range cases {
		tokenizer := newTokenizer(c.input)
		p := NewParser(tokenizer)
		stmt := p.stmt()
		fmt.Println("expected: " + stmt.string() + ", actual: " + c.expected)
		if stmt.string() != c.expected {
			fmt.Println(stmt.string())
			t.Error("The ast is wrong\n")
		}
	}
}
