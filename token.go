package main

import "fmt"

const TTInt    = "INT"
const TTFloat  = "FLOAT"
const TTOp     = "OP"
const TTParen  = "PAREN"

type Token struct {
	Type string
	Value interface{}
}

func NewToken(t string, v interface{}) *Token {
	token := &Token{
		Type: t,
		Value: v,
	}

	return token
}

func (t *Token) ToString() string {
	return fmt.Sprintf("[%v: %v]", t.Type, t.Value)
}
