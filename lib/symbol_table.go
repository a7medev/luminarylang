package main

type SymbolTable struct {
	Symbols map[string]Value
	Parent *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		Symbols: map[string]Value{},
	}
	return st
}

func (st *SymbolTable) Get(n string) Value {
	val := st.Symbols[n]
	if val == nil && st.Parent != nil {
		return st.Parent.Get(n)
	}
	return val
}

func (st *SymbolTable) Set(n string, v Value) Value {
	st.Symbols[n] = v
	return v
}

func (st *SymbolTable) Del(n string) {
	delete(st.Symbols, n)
}
