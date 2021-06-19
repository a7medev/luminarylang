package main

import "fmt"

const TTInt    = "INT"
const TTFloat  = "FLOAT"
const TTPlus   = "PLUS"
const TTMinus  = "MINUS"
const TTMul    = "MUL"
const TTDiv    = "DIV"
const TTLParen = "LPAREN"
const TTRParen = "RPAREN"

type Token struct {
	Type, Value string
}

func NewToken(t, v string) *Token {
	token := &Token{
		Type: t,
		Value: v,
	}

	return token
}

func (t *Token) ToString() string {
	return fmt.Sprintf("[%v: %v]", t.Type, t.Value)
}