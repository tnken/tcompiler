package token

import (
	"fmt"
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
		while a > 10 do
			a = a + 3
		end
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
		{KeyWhile, "while"},
		{Identifier, "a"},
		{GreaterThan, ">"},
		{Num, "10"},
		{KeyDo, "do"},
		{Identifier, "a"},
		{Assign, "="},
		{Identifier, "a"},
		{Plus, "+"},
		{Num, "3"},
		{KeyEnd, "end"},
		{EOF, ""},
	}
	tokenizer := New(input1)
	for _, c := range case1 {
		token, _ := tokenizer.Next()
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
	tokenizer = New(input2)
	for _, c := range case2 {
		token, _ := tokenizer.Next()
		if token.Kind != c.expectKind {
			t.Error("The token kind is wrong\n")
		}

		if token.Literal != c.expectLiteral {
			t.Error("The token literal is wrong\n")
		}
	}

	input3 := "def myFunc(hoge) a = 33 return hoge+a end myFunc()"
	case3 := []struct {
		expectKind    Kind
		expectLiteral string
	}{
		{KeyDef, "def"},
		{Identifier, "myFunc"},
		{LParen, "("},
		{Identifier, "hoge"},
		{RParen, ")"},
		{Identifier, "a"},
		{Assign, "="},
		{Num, "33"},
		{KeyReturn, "return"},
		{Identifier, "hoge"},
		{Plus, "+"},
		{Identifier, "a"},
		{KeyEnd, "end"},
		{Identifier, "myFunc"},
		{LParen, "("},
		{RParen, ")"},
		{EOF, ""},
	}
	tokenizer = New(input3)
	for _, c := range case3 {
		token, _ := tokenizer.Next()
		if token.Kind != c.expectKind {
			t.Error("The token kind is wrong\n")
		}

		if token.Literal != c.expectLiteral {
			t.Error("The token literal is wrong\n")
		}
	}

	input4 := `
class LED
	def on(num)
		# this is comment
		self.pin = num
	end
end

# this is also comment
a = LED()
a.on(3) # call on method`
	case4 := []struct {
		expectKind    Kind
		expectLiteral string
	}{
		{KeyClass, "class"},
		{Identifier, "LED"},
		{KeyDef, "def"},
		{Identifier, "on"},
		{LParen, "("},
		{Identifier, "num"},
		{RParen, ")"},
		{KeySelf, "self"},
		{Dot, "."},
		{Identifier, "pin"},
		{Assign, "="},
		{Identifier, "num"},
		{KeyEnd, "end"},
		{KeyEnd, "end"},
		{Identifier, "a"},
		{Assign, "="},
		{Identifier, "LED"},
		{LParen, "("},
		{RParen, ")"},
		{Identifier, "a"},
		{Dot, "."},
		{Identifier, "on"},
		{LParen, "("},
		{Num, "3"},
		{RParen, ")"},
		{EOF, ""},
	}
	tokenizer = New(input4)
	for _, c := range case4 {
		token, _ := tokenizer.Next()
		if token.Kind != c.expectKind {
			t.Error("The token kind is wrong\n")
			fmt.Println(token.Kind)
		}

		if token.Literal != c.expectLiteral {
			fmt.Println("expected: " + c.expectLiteral)
			fmt.Println("but actual: " + token.Literal)
			t.Error("The token literal is wrong\n")
		}
	}
}
