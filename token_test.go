package main

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	input1 := `
		a = 1 + 20 - 300 * 4 / 5
		testfn(a)
		loop {
			a = 1 + 20 - 300 * 4 / 5
			b = a
			print(b)
		}
		if a > 0 then
			b = 3
		end
		a == 3
		a != 3
		`

	case1 := []struct {
		expectKind    Kind
		expectLiteral string
	}{
		{Identifier, "a"},
		{Assign, "="},
		{Num, "1"},
		{Plus, "+"},
		{Num, "20"},
		{Minus, "-"},
		{Num, "300"},
		{Asterisk, "*"},
		{Num, "4"},
		{Slash, "/"},
		{Num, "5"},
		{Identifier, "testfn"},
		{LParen, "("},
		{Identifier, "a"},
		{RParen, ")"},
		{KeyLoop, "loop"},
		{Lbrace, "{"},
		{Identifier, "a"},
		{Assign, "="},
		{Num, "1"},
		{Plus, "+"},
		{Num, "20"},
		{Minus, "-"},
		{Num, "300"},
		{Asterisk, "*"},
		{Num, "4"},
		{Slash, "/"},
		{Num, "5"},
		{Identifier, "b"},
		{Assign, "="},
		{Identifier, "a"},
		{Identifier, "print"},
		{LParen, "("},
		{Identifier, "b"},
		{RParen, ")"},
		{Rbrace, "}"},
		{KeyIf, "if"},
		{Identifier, "a"},
		{GreaterThan, ">"},
		{Num, "0"},
		{KeyThen, "then"},
		{Identifier, "b"},
		{Assign, "="},
		{Num, "3"},
		{KeyEnd, "end"},
		{Identifier, "a"},
		{Eq, "=="},
		{Num, "3"},
		{Identifier, "a"},
		{NEq, "!="},
		{Num, "3"},
		{EOF, ""},
	}
	tokenizer := NewToken(input1)
	for _, c := range case1 {
		token := tokenizer.Next()
		if token.Kind != c.expectKind {
			t.Error("The token kind is wrong\n")
		}

		if token.Literal != c.expectLiteral {
			t.Error("The token literal is wrong\n")
		}
	}

	input2 := "hoge = [1, 2, 3, 2+2]"
	case2 := []struct {
		expectKind    Kind
		expectLiteral string
	}{
		{Identifier, "hoge"},
		{Assign, "="},
		{Lbracket, "["},
		{Num, "1"},
		{Comma, ","},
		{Num, "2"},
		{Comma, ","},
		{Num, "3"},
		{Comma, ","},
		{Num, "2"},
		{Plus, "+"},
		{Num, "2"},
		{Rbracket, "]"},
		{EOF, ""},
	}
	tokenizer = NewToken(input2)
	for _, c := range case2 {
		token := tokenizer.Next()
		if token.Kind != c.expectKind {
			t.Error("The token kind is wrong\n")
		}

		if token.Literal != c.expectLiteral {
			t.Error("The token literal is wrong\n")
		}
	}
}
