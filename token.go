package main

import "fmt"

const TTInt   = "INT"
const TTFloat = "FLOAT"
const TTOp    = "OP"
const TTParen = "PAREN"
const TTEOF   = "EOF"

type Token struct {
	Type string
	Value interface{}
	Pos *Position
}

func NewToken(t string, v interface{}, p *Position) *Token {
	token := &Token{
		Type: t,
		Value: v,
		Pos: p,
	}

	return token
}

func (t *Token) ToString() string {
	return fmt.Sprintf("[%v: %v]", t.Type, t.Value)
}
