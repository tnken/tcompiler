package main

import (
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"1", "1"},
		{"1+2*3", "7"},
		{" 1 * 2 +  3  ", "5"},
		{"[1 * 2 +  3, [1+1, 2+2, 3*3]]", "[5 [2 4 9]]"},
	}

	for _, c := range cases {
		tokenizer := NewTokenizer(c.input)
		p := NewParser(tokenizer)
		exp := p.expr(Lowest)
		if eval(exp).stringVal() != c.expected {
			t.Error("The evaluate result is wrong\n")
		}
	}
}
