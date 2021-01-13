package parser

import (
	"fmt"
	"testing"

	"github.com/takeru56/tcompiler/token"
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
		{
			`while 3 > 1 do
  b = 3+5
  b+2
end`,
			`while (3 > 1) do
  b = (3 + 5)
  (b + 2)
end`},
	}

	for _, c := range cases {
		tokenizer := token.New(c.input)
		p, _ := New(tokenizer)
		stmt, _ := p.stmt()
		if stmt.string() != c.expected {
			fmt.Println(stmt.string())
			t.Error("The ast is wrong\n")
		}
	}
}
