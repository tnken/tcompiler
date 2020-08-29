package parser

import (
	"fmt"
	"testing"

	"github.com/takeru56/t/token"
)

func TestParser(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"1", "1"},
		{"1+2*3", "(1 + (2 * 3))"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
	}

	for _, c := range cases {
		tokenizer := token.New(c.input)
		p := New(tokenizer)
		stmt := p.stmt()
		fmt.Println("actual: " + stmt.string() + ", expected: " + c.expected)
		if stmt.string() != c.expected {
			fmt.Println(stmt.string())
			t.Error("The ast is wrong\n")
		}
	}
}
