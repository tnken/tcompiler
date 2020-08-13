package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/tarm/serial"
)

type Serial struct {
	port   string
	baud   int
	isOpen bool
	p      *serial.Port
}

func newSerial(port string, baud int) *Serial {
	c := &serial.Config{Name: port, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return &Serial{port: port, baud: baud, isOpen: false, p: s}
	}
	return &Serial{port: port, baud: baud, isOpen: true, p: s}
}

func (s *Serial) write(b byte) error {
	if !s.isOpen {
		return errors.New("error: port is not opened")
	}
	_, err := s.p.Write([]byte{b})
	if err != nil {
		return err
	}
	return nil
}

func repl(port string) {
	stdin := bufio.NewScanner(os.Stdin)
	e := newEval(port)
	fmt.Print(">> ")
	for stdin.Scan() {
		text := stdin.Text()
		if text == "exit" {
			break
		}
		tokenizer := newTokenizer(text)
		p := newParser(tokenizer)
		stmt := p.stmt()
		fmt.Println("=> " + e.eval(stmt).stringVal())
		fmt.Print(">> ")
	}
}

func main() {
	port := "/dev/tty.usbmodem14501"
	repl(port)
}
