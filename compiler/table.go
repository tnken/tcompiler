package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store       map[string]Symbol
	symbolCount int
	outerScope  *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s, symbolCount: 0}
}

func (st *SymbolTable) DefineGlobal(name string) Symbol {
	symbol := Symbol{Name: name, Scope: GlobalScope, Index: st.symbolCount}
	st.store[name] = symbol
	st.symbolCount += 1
	return symbol
}

func (st *SymbolTable) DefineLocal(name string) Symbol {
	symbol := Symbol{Name: name, Scope: LocalScope, Index: st.symbolCount}
	st.store[name] = symbol
	st.symbolCount++
	return symbol
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	return sym, ok
}

type Class struct {
	Name             string
	Index            int
	instanceValTable map[string]int
	instanceValCount int
	hasInit          bool
}

func NewClass(name string, index int, hasInit bool) *Class {
	t := make(map[string]int)
	return &Class{Name: name, Index: index, instanceValTable: t, instanceValCount: 0, hasInit: hasInit}
}

type ClassTable struct {
	store      map[string]*Class
	classCount int
}

func NewClassTable() *ClassTable {
	c := make(map[string]*Class)
	return &ClassTable{store: c, classCount: 0}
}

func (ct *ClassTable) DefineClass(name string) Class {
	class := NewClass(name, ct.classCount, false)
	ct.store[name] = class
	ct.classCount++
	return *class
}

func (ct *ClassTable) Resolve(name string) (*Class, bool) {
	class, ok := ct.store[name]
	return class, ok
}

func (c *Class) DefineInstanceVal(name string) int {
	c.instanceValTable[name] = c.instanceValCount
	c.instanceValCount++
	return c.instanceValTable[name]
}

func (c *Class) ResolveInstanceVal(name string) (int, bool) {
	id, ok := c.instanceValTable[name]
	return id, ok
}

type MethodTable struct {
	store       map[string]int
	methodCount int
}

func NewMethodTable() *MethodTable {
	m := make(map[string]int)
	// for constructor
	m["init"] = 0
	return &MethodTable{store: m, methodCount: 1}
}

func (mt *MethodTable) DefineMethodId(name string) int {
	_, ok := mt.store[name]
	if ok {
		return mt.store[name]
	}
	mt.store[name] = mt.methodCount
	mt.methodCount++
	return mt.store[name]
}

func (mt *MethodTable) ResolveMethodId(name string) (int, bool) {
	id, ok := mt.store[name]
	return id, ok
}
