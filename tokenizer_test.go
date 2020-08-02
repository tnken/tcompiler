package main

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	input1 := "1+2-3"
	case1 := []struct {
		expectKind    TokenKind
		expectLiteral string
	}{
		{Num, "1"},
		{Plus, "+"},
		{Num, "2"},
		{Minus, "-"},
		{Num, "3"},
		{Eof, ""},
	}

	tokenizer := NewTokenizer(input1)
	for _, c := range case1 {
		token := tokenizer.next()
		if token.Kind != c.expectKind {
			t.Error("The token kind is wrong\n")
		}

		if token.Literal != c.expectLiteral {
			t.Error("The token literal is wrong\n")
		}
	}
}
