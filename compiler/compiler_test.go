package compiler

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
)

// output tarto IR bytecode Format
// ***************************************

// struct {
// 	u4 magic
// 	u2 constant_pool_count
// 	cp constant_pool[constant_pool_count]
// 	u2 instruction_count
// 	ins instructions[instruction_count]
// }

// struct constant_pool {
// 	u1 constant type
// 	u2 constant size
// 	c [const size]constants
// }
// ***************************************

func TestCompile(t *testing.T) {
	cases := []struct {
		source   string
		bytecode []byte
	}{
		{"23", []byte{0, 1, 0, 0, 2, 0, 23, 0, 4, 0, 0, 1, 5}},
		{"256+1", []byte{0, 2, 0, 0, 2, 1, 0, 0, 0, 2, 0, 1, 0, 8, 0, 0, 1, 0, 0, 2, 1, 5}},
		{"1-1", []byte{0, 2, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 8, 0, 0, 1, 0, 0, 2, 2, 5}},
		{"1*1", []byte{0, 2, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 8, 0, 0, 1, 0, 0, 2, 3, 5}},
		{"1/1", []byte{0, 2, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 8, 0, 0, 1, 0, 0, 2, 4, 5}},
		{"1>1", []byte{0, 2, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 8, 0, 0, 1, 0, 0, 2, 9, 5}},
		{"a = 1", []byte{0, 1, 0, 0, 2, 0, 1, 0, 6, 0, 0, 1, 11, 0, 5}},
		{"a = 2 a == 2", []byte{0, 2, 0, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 12, 0, 0, 1, 11, 0, 10, 0, 0, 0, 2, 6, 5}},
		{"a = 1 b = 2 b", []byte{0, 2, 0, 0, 2, 0, 1, 0, 0, 2, 0, 2, 0, 13, 0, 0, 1, 11, 0, 0, 0, 2, 11, 1, 10, 1, 5}},
		{"if 1 > 1 do 1+1 end a = 1", []byte{0, 5, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 23, 0, 0, 1, 0, 0, 2, 9, 12, 0, 17, 0, 0, 3, 0, 0, 4, 1, 0, 0, 5, 11, 0, 5}},
		{"while 1 > 0 do 1 end 1", []byte{0, 4, 0, 0, 2, 0, 1, 0, 0, 2, 0, 0, 0, 0, 2, 0, 1, 0, 0, 2, 0, 1, 0, 20, 0, 0, 1, 0, 0, 2, 9, 12, 0, 16, 0, 0, 3, 13, 0, 0, 0, 0, 4, 5}},
		{"a = 1 while 5 > a do a=a+1 end a", []byte{0, 3, 0, 0, 2, 0, 1, 0, 0, 2, 0, 5, 0, 0, 2, 0, 1, 0, 28, 0, 0, 1, 11, 0, 0, 0, 2, 10, 0, 9, 12, 0, 25, 10, 0, 0, 0, 3, 1, 11, 0, 13, 0, 5, 10, 0, 5}},
		{"def myFunc() 2+3 end", []byte{0, 3, 0, 0, 2, 0, 2, 0, 0, 2, 0, 3, 1, 0, 7, 0, 0, 1, 0, 0, 2, 1, 0, 6, 0, 0, 3, 11, 0, 5}},
		{"def myFunc() 2+3 end myFunc()", []byte{0, 3, 0, 0, 2, 0, 2, 0, 0, 2, 0, 3, 1, 0, 7, 0, 0, 1, 0, 0, 2, 1, 0, 8, 0, 0, 3, 11, 0, 14, 0, 5}},
	}

	for _, c := range cases {
		out, err := exec.Command("go", "run", "../", c.source).Output()
		if err != nil {
			log.Fatal(err)
		}
		s := "ffffffff"
		for _, b := range c.bytecode {
			s += fmt.Sprintf("%02x", b)
		}

		if string(out) != s {
			fmt.Println("expected: " + s)
			fmt.Println("but actual: " + string(out))
			t.Error("not match\n")
		}
	}
}
