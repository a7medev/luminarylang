package main

import "fmt"

type Error struct {
	Name, Details string
	Pos *Position
}

func NewError(n, d string, p *Position) *Error {
	e := &Error{
		Name: n,
		Details: d,
		Pos: p,
	}
	return e
}

func (e *Error) ToString() string {
	return fmt.Sprintf(
		"Error(%v): %v.\nFile: %v - Line: %v - Col: %v",
		e.Name,
		e.Details,
		e.Pos.FileName,
		e.Pos.Line,
		e.Pos.Col,
	)
}

func NewIlligalCharError(d string, p *Position) *Error {
	e := &Error{
		Name: "Illigal Char",
		Details: d,
		Pos: p,
	}
	return e
}
