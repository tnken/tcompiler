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
