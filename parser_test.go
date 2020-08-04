package main

import (
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
	}

	for _, c := range cases {
		tokenizer := NewTokenizer(c.input)
		p := NewParser(tokenizer)
		exp := p.expr(Lowest)
		if exp.string() != c.expected {
			t.Error("The ast is wrong\n")
		}
	}
}
