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
		{
			`if 3>1 then
  b = 3+5
  b+2
end`,
			`if (3 > 1) then
  b = (3 + 5)
  (b + 2)
end`},
	}

	for _, c := range cases {
		tokenizer := NewToken(c.input)
		p := NewParser(tokenizer)
		stmt := p.stmt()
		if stmt.string() != c.expected {
			fmt.Println(stmt.string())
			t.Error("The ast is wrong\n")
		}
	}
}
