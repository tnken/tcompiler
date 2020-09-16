package main

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
)

func TestCompile(t *testing.T) {
	cases := []struct {
		source   string
		bytecode []byte
	}{
		{"23", []byte{0, 0, 23, 5}},
		{"256+1", []byte{0, 1, 0, 0, 0, 1, 1, 5}},
		{"1-1", []byte{0, 0, 1, 0, 0, 1, 2, 5}},
		{"1*1", []byte{0, 0, 1, 0, 0, 1, 3, 5}},
		{"1/1", []byte{0, 0, 1, 0, 0, 1, 4, 5}},
		{"1>1", []byte{0, 0, 1, 0, 0, 1, 9, 5}},
		{"a = 1", []byte{0, 0, 1, 11, 0, 5}},
		{"a = 2 a == 2", []byte{0, 0, 2, 11, 0, 10, 0, 0, 0, 2, 6, 5}},
		{"a = 1 b = 2 b", []byte{0, 0, 1, 11, 0, 0, 0, 2, 11, 1, 10, 1, 5}},
		{"if 1 > 1 do 1+1 end a = 1", []byte{0, 0, 1, 0, 0, 1, 9, 12, 0, 8, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 11, 0, 5}},
	}

	for _, c := range cases {
		out, err := exec.Command("go", "run", ".", c.source).Output()
		if err != nil {
			log.Fatal(err)
		}
		s := ""
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
