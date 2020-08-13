package main

import (
	"fmt"
	"strconv"
	"time"
)

type Object interface {
	stringVal() string
}

type Integer struct {
	value int
}

func (i Integer) stringVal() string { return strconv.Itoa(i.value) }

type Array struct {
	val []Object
}

func (arr Array) stringVal() string {
	s := "["
	for i, v := range arr.val {
		switch ele := v.(type) {
		//TODO: remove redundancy
		case Integer:
			if i == 0 {
				s += strconv.Itoa(ele.value)
			} else {
				s += " " + strconv.Itoa(ele.value)
			}
		case Array:
			if i == 0 {
				s += ele.stringVal()
			} else {
				s += " " + ele.stringVal()
			}
		}
	}
	return s + "]"
}

type Var struct {
	name string
	obj  Object
}

func (v Var) stringVal() string { return v.obj.stringVal() }

type Nil struct {
	name string
}

func (n Nil) stringVal() string { return "nil" }

// eval
type Eval struct {
	port string
	vars map[string]Var
}

func newEval(p string) Eval {
	return Eval{port: p, vars: map[string]Var{}}
}

func (e Eval) eval(node Node) Object {
	return e.stmt(node.(Stmt))
}

func (e Eval) stmt(stmt Stmt) Object {
	switch s := stmt.(type) {
	case VarDecl:
		name := s.left.name
		v := Var{name: name, obj: e.expr(s.right)}
		e.vars[name] = v
		return v
	case ExprStmt:
		return e.expr(s.val)
	case LoopStmt:
		for {
			for _, line := range s.block {
				e.stmt(line)
			}
		}
	}
	panic("error")
}

// Tree Walk
func (e Eval) expr(expr Expr) Object {
	switch v := expr.(type) {
	case InfixExpr:
		l := e.expr(v.left).(Integer).value
		r := e.expr(v.right).(Integer).value
		switch v.op {
		case Plus:
			return Integer{value: l + r}
		case Minus:
			return Integer{value: l - r}
		case Asterisk:
			return Integer{value: l * r}
		case Slash:
			return Integer{value: int(l / r)}
		}
	case NumberLiteral:
		i, _ := strconv.Atoi(v.val)
		return Integer{value: i}
	case ArrayInit:
		val := []Object{}
		for _, ele := range v.exprs {
			val = append(val, e.expr(ele))
		}
		return Array{val: val}
	case Ident:
		return e.vars[v.name].obj
	case FnCallExpr:
		switch v.ident.name {
		// builtin functions
		case "digitalwrite":
			// TODO: to be simple
			serial := newSerial(e.port, 9600)
			if v.args[0].(NumberLiteral).string() == "1" {
				serial.write('1')
			} else {
				serial.write('0')
			}
		case "sleep":
			t := time.Duration(e.expr(v.args[0]).(Integer).value) * time.Second
			time.Sleep(t)
		case "print":
			fmt.Println(e.expr(v.args[0]).stringVal())
		}
		return Nil{}
	}
	return nil
}
