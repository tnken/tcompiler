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
		{"a=1+2*3", "7"},
		{"b = [1 * 2 +  3, [1+1, 2+2, 3*3]]", "[5 [2 4 9]]"},
		{"c = 2 + a", "9"},
		{"c", "9"},
		{"1+c", "10"},
	}

	port := ""
	e := newEval(port)
	for _, c := range cases {
		tokenizer := newTokenizer(c.input)
		p := NewParser(tokenizer)
		stmt := p.stmt()
		if e.eval(stmt).stringVal() != c.expected {
			t.Error("The evaluate result is wrong\n")
		}
	}
}
