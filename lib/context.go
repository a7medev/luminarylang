package main

type Context struct {
	Name string
	SymbolTable *SymbolTable
	Parent *Context
}

func NewContext(n string, ) *Context {
	c := &Context{
		Name: n,
	}

	return c
}
