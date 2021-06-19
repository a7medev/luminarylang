package main

type Error struct {
	Name, Details string
}

func NewError(n, d string) *Error {
	e := &Error{
		Name: n,
		Details: d,
	}
	return e
}

func NewIlligalCharError(d string) *Error {
	e := &Error{
		Name: "Illigal Char",
		Details: d,
	}
	return e
}
