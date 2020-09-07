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
		{"1+1", []byte{0, 0, 1, 0, 0, 1, 1, 5}},
		{"1-1", []byte{0, 0, 1, 0, 0, 1, 2, 5}},
		{"1*1", []byte{0, 0, 1, 0, 0, 1, 3, 5}},
		{"1/1", []byte{0, 0, 1, 0, 0, 1, 4, 5}},
	}

	for _, c := range cases {
		out, err := exec.Command("go", "run", ".", c.source).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(out))
		s := ""
		for _, b := range c.bytecode {
			s += fmt.Sprintf("%02x", b)
		}
		fmt.Println(s)

		if string(out) != s {
			fmt.Println(string(c.bytecode))
			t.Error("not match\n")
		}
	}
}
