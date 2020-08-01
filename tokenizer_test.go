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
		{Reserved, "+"},
		{Num, "2"},
		{Reserved, "-"},
		{Num, "3"},
		{Eof, ""},
	}

	tokenizer := NewTokenizer(input1)
	tokens := tokenizer.Run()
	for i, c := range case1 {
		if tokens[i].Kind != c.expectKind {
			t.Error("The token kind is wrong\n")
		}

		if tokens[i].Literal != c.expectLiteral {
			t.Error("The token literal is wrong\n")
		}
	}
}
