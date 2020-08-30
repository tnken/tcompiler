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
		{"23", []byte{0, 0, 23}},
	}

	for _, c := range cases {
		out, err := exec.Command("go", "run", ".", c.source).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(c.bytecode))

		s := ""
		for _, b := range c.bytecode {
			s += fmt.Sprintf("%02x\n", b)
		}

		if string(out) != s {
			fmt.Println(string(c.bytecode))
			t.Error("not match\n")
		}
	}
}
