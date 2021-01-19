package parser

import (
	"fmt"
	"testing"

	"github.com/takeru56/tcompiler/token"
)

func TestParser(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"1", []string{"1"}},
		{"1+2*3", []string{"(1 + (2 * 3))"}},
		{"1 * 2 + 3", []string{"((1 * 2) + 3)"}},
		{"a=1+1", []string{"a = (1 + 1)"}},
		{
			`if 3>1 do
  b = 3+5
  b+2
end`,
			[]string{`if (3 > 1) then
  b = (3 + 5)
  (b + 2)
end`}},
		{
			`while 3 > 1 do
  b = 3+5
  b+2
end`,
			[]string{`while (3 > 1) do
  b = (3 + 5)
  (b + 2)
end`}},
		{
			`def myFunc()
  b = 1+1
  b+2
  return b
end
myFunc()`,
			[]string{`def myFunc()
  b = (1 + 1)
  (b + 2)
  return b
end`,
				"myFunc()"}},
		{
			`
def myFunc(a)
  return a+1
end
return myFunc()+1`,
			[]string{`def myFunc(a)
  return (a + 1)
end`,
				"return (myFunc() + 1)"}},
	}

	for _, c := range cases {
		tokenizer := token.New(c.input)
		p, _ := New(tokenizer)

		i := 0
		for p.curToken.Kind != token.EOF {
			stmt, _ := p.stmt()
			fmt.Println(p.curToken.Literal)
			if stmt.string() != c.expected[i] {
				fmt.Println("expecting: \n" + c.expected[i])
				fmt.Println("but actual: \n" + stmt.string())
				t.Error("The ast is wrong\n")
			}
			i++
		}
	}
}
