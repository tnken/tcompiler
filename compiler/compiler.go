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
	p            []parser.Node
	constantPool []obj.Object
	scopes       []CompilationScope
	scopeIndex   int
	symbolTable  *SymbolTable
}

func newCompiler(program []parser.Node) *Compiler {
	main := CompilationScope{}
	c := &Compiler{program, []obj.Object{}, []CompilationScope{main}, 0, NewSymbolTable()}
	return c
}

type CompilationScope struct {
	instructions code.Instructions
}

func (c *Compiler) enterScope() {
	c.scopeIndex++
	c.scopes = append(c.scopes, CompilationScope{})
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.scopes[c.scopeIndex].instructions
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	return instructions
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
		symbol, ok := c.symbolTable.Resolve(node.Name)
		if ok {
			c.emit(code.OpLoadGlobal, []int{symbol.Index}...)
			return
		}
		// TODO: do error handling, when ok is false
		fmt.Println("undefined variable")
		os.Exit(1)

	case parser.AssignStmt:
		c.gen(node.Expr)
		symbol, ok := c.symbolTable.Resolve(node.Ident.Name)
		if ok {
			c.emit(code.OpStoreGlobal, []int{symbol.Index}...)
			return
		}
		global := c.symbolTable.Define(node.Ident.Name)
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
		c.enterScope()
		for _, stmt := range node.Block.Nodes {
			c.gen(stmt)
		}
		instructions := c.leaveScope()
		objFunc := &obj.Function{Instructions: instructions}
		c.emit(code.OpConstant, []int{c.addConstant(objFunc)}...)

		symbol, ok := c.symbolTable.Resolve(node.Ident.Name)
		if ok {
			c.emit(code.OpStoreGlobal, []int{symbol.Index}...)
			return
		}
		global := c.symbolTable.Define(node.Ident.Name)
		c.emit(code.OpStoreGlobal, []int{global.Index}...)
	case parser.CallExpr:
		symbol, ok := c.symbolTable.Resolve(node.Ident.Name)
		if ok {
			c.emit(code.OpCall, []int{symbol.Index}...)
			return
		}
		// TODO: 未定義関数呼び出しのエラーハンドル
		// ひとまず握りつぶす
		fmt.Println("undefined function")
		os.Exit(1)
	}
}
