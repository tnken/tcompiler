package compiler

import (
	"fmt"
	"os"

	"github.com/takeru56/tcompiler/code"
	"github.com/takeru56/tcompiler/obj"
	"github.com/takeru56/tcompiler/parser"
)

func (c *Compiler) emit(op code.Opcode, operands ...int) {
	ins := code.Make(op, operands...)
	for _, i := range ins {
		c.scopes[c.scopeIndex].instructions = append(c.scopes[c.scopeIndex].instructions, i)
	}
}

type Compiler struct {
	p              []parser.Node
	constantPool   []obj.Object
	scopes         []CompilationScope
	scopeIndex     int
	cTable         *ClassTable
	classPool      []obj.Class
	FlagClassScope bool
}

func newCompiler(program []parser.Node) *Compiler {
	main := CompilationScope{table: NewSymbolTable()}
	c := &Compiler{program, []obj.Object{}, []CompilationScope{main}, 0, NewClassTable(), []obj.Class{}, false}
	return c
}

type CompilationScope struct {
	instructions code.Instructions
	numLocal     int
	table        *SymbolTable
}

func (c *Compiler) enterClass() {
	c.FlagClassScope = true
}

func (c *Compiler) leaveClass() obj.Class {
	c.FlagClassScope = false
	return c.classPool[len(c.classPool)-1]
}

func (c *Compiler) enterScope() {
	t := NewSymbolTable()
	t.outerScope = c.currentScope().table
	c.scopeIndex++
	c.scopes = append(c.scopes, CompilationScope{numLocal: 0, table: t})
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.scopes[c.scopeIndex].instructions
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	return instructions
}

func (c *Compiler) currentScope() *CompilationScope {
	return &c.scopes[c.scopeIndex]
}

func Exec(program []parser.Node) *Compiler {
	c := newCompiler(program)
	for _, node := range program {
		c.gen(node)
	}
	c.emit(code.OpDone, []int{}...)
	return c
}

func (c *Compiler) addConstant(obj obj.Object) int {
	if c.FlagClassScope {
		c.classPool[len(c.classPool)-1].ConstantPool = append(c.classPool[len(c.classPool)-1].ConstantPool, obj)
		return len(c.constantPool)
	}
	c.constantPool = append(c.constantPool, obj)
	return len(c.constantPool)
}

func (c *Compiler) gen(n parser.Node) {
	switch node := n.(type) {
	case parser.IntegerLiteral:
		integer := &obj.Integer{Value: node.Val}
		c.emit(code.OpConstant, []int{c.addConstant(integer)}...)
	case parser.InfixExpr:
		c.gen(node.Left)
		c.gen(node.Right)
		switch node.Op {
		case parser.Add:
			c.emit(code.OpAdd, []int{}...)
		case parser.Sub:
			c.emit(code.OpSub, []int{}...)
		case parser.Mul:
			c.emit(code.OpMul, []int{}...)
		case parser.Div:
			c.emit(code.OpDiv, []int{}...)
		case parser.EQ:
			c.emit(code.OpEQ, []int{}...)
		case parser.NEQ:
			c.emit(code.OpNEQ, []int{}...)
		case parser.Less:
			c.emit(code.OpLess, []int{}...)
		case parser.Greater:
			c.emit(code.OpGreater, []int{}...)
		}
	case parser.IdentExpr:
		if c.scopeIndex > 0 {
			symbol, ok := c.currentScope().table.Resolve(node.Name)
			if ok {
				c.emit(code.OpLoadLocal, []int{symbol.Index}...)
				return
			}
			symbol, ok = c.currentScope().table.outerScope.Resolve(node.Name)
			if ok {
				c.emit(code.OpLoadGlobal, []int{symbol.Index}...)
				return
			}
		}
		symbol, ok := c.currentScope().table.Resolve(node.Name)
		if ok {
			c.emit(code.OpLoadGlobal, []int{symbol.Index}...)
			return
		}

		// ひとまず握りつぶしとく
		fmt.Println("Undefined identifier")
		os.Exit(1)

	case parser.AssignStmt:
		c.gen(node.Expr)

		// local variable
		if c.scopeIndex > 0 {
			symbol, ok := c.currentScope().table.Resolve(node.Ident.Name)
			if ok {
				c.emit(code.OpStoreLocal, []int{symbol.Index}...)
				return
			}
			local := c.currentScope().table.DefineLocal(node.Ident.Name)
			c.emit(code.OpStoreLocal, []int{local.Index}...)
			return
		}
		// global variable
		symbol, ok := c.currentScope().table.Resolve(node.Ident.Name)
		if ok {
			c.emit(code.OpStoreGlobal, []int{symbol.Index}...)
			return
		}
		global := c.currentScope().table.DefineGlobal(node.Ident.Name)
		c.emit(code.OpStoreGlobal, []int{global.Index}...)
	case parser.IfStmt:
		c.gen(node.Condition)
		c.emit(code.OpJNT, []int{0}...)
		blockHead := len(c.scopes[c.scopeIndex].instructions)
		ifHead := blockHead - 3
		for _, stmt := range node.Block.Nodes {
			c.gen(stmt)
		}
		ins := code.Make(code.OpJNT, []int{len(c.scopes[c.scopeIndex].instructions)}...)

		c.scopes[c.scopeIndex].instructions[ifHead+1] = ins[1]
		c.scopes[c.scopeIndex].instructions[ifHead+2] = ins[2]
	case parser.WhileStmt:
		head := len(c.scopes[c.scopeIndex].instructions)
		c.gen(node.Condition)
		c.emit(code.OpJNT, []int{0}...)
		blockHead := len(c.scopes[c.scopeIndex].instructions)
		whileHead := blockHead - 3
		for _, stmt := range node.Block.Nodes {
			c.gen(stmt)
		}
		c.emit(code.OpJMP, []int{head}...)

		ins := code.Make(code.OpJNT, []int{len(c.scopes[c.scopeIndex].instructions)}...)
		c.scopes[c.scopeIndex].instructions[whileHead+1] = ins[1]
		c.scopes[c.scopeIndex].instructions[whileHead+2] = ins[2]
	case parser.FunctionDef:
		if c.FlagClassScope {
			c.enterScope()
			for _, arg := range node.Args {
				c.currentScope().table.DefineLocal(arg.Name)
			}
			for _, stmt := range node.Block.Nodes {
				c.gen(stmt)
			}
			instructions := c.leaveScope()
			objFunc := &obj.Function{Instructions: instructions, NumArg: len(node.Args)}
			c.classPool[len(c.classPool)-1].ConstantPool = append(c.classPool[len(c.classPool)-1].ConstantPool, objFunc)
			return
		}
		symbol, ok := c.currentScope().table.Resolve(node.Ident.Name)
		if !ok {
			symbol = c.currentScope().table.DefineGlobal(node.Ident.Name)
		}

		c.enterScope()
		for _, arg := range node.Args {
			c.currentScope().table.DefineLocal(arg.Name)
		}
		for _, stmt := range node.Block.Nodes {
			c.gen(stmt)
		}
		instructions := c.leaveScope()
		objFunc := &obj.Function{Instructions: instructions, NumArg: len(node.Args)}
		c.emit(code.OpConstant, []int{c.addConstant(objFunc)}...)

		if ok {
			c.emit(code.OpStoreGlobal, []int{symbol.Index}...)
			return
		}
		c.emit(code.OpStoreGlobal, []int{symbol.Index}...)
	case parser.CallExpr:
		c.gen(node.Ident)
		for _, expr := range node.Args {
			c.gen(expr)
		}
		c.emit(code.OpCall, []int{len(node.Args)}...)
	case parser.ReturnStmt:
		c.gen(node.Expr)
		c.emit(code.OpReturn, []int{}...)
	case parser.ClassDef:
		ct := c.cTable.DefineClass(node.Ident.Name)
		c.classPool = append(c.classPool, obj.Class{Index: ct.Index, NumInstanceVal: 0, NumMethod: 0, ConstantPool: []obj.Object{}})
		c.enterClass()
		for _, method := range node.Methods {
			c.gen(method)
		}
		c.leaveClass()
	}
}
