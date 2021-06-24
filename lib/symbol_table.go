package main

type SymbolTable struct {
	Symbols map[string]Value
	Parent *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		Symbols: map[string]Value{},
	}
	st.Init()
	return st
}

func (st *SymbolTable) Init() {
	st.Set("true", NewNumber(1))
	st.Set("false", NewNumber(0))
	

	//* BUILTIN FUNCTIONS

	// Stdin/Stdout/System
	st.Set("print", BuiltinPrint)
	st.Set("println", BuiltinPrintln)
	st.Set("scan", BuiltinScan)
	st.Set("exit", BuiltinExit)

	// Arrays
	st.Set("len", BuiltinLen)
	st.Set("append", BuiltinAppend)
	st.Set("prepend", BuiltinPrepend)
	st.Set("pop", BuiltinPop)
	st.Set("shift", BuiltinShift)
	st.Set("map", BuiltinMap)
	st.Set("reduce", BuiltinReduce)

	// Strings
	st.Set("trim", BuiltinTrim)
	st.Set("replace", BuiltinReplace)
	st.Set("upper", BuiltinUpper)
	st.Set("lower", BuiltinLower)

	// Conversion
	st.Set("num", BuiltinNum)
	st.Set("str", BuiltinStr)
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
