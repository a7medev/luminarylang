package main

import "fmt"

const TTNum     = "NUM"
const TTOp      = "OP"
const TTId      = "ID"
const TTEOF     = "EOF"
const TTKeyword = "KEYWORD"

type Token struct {
	Type string
	Value interface{}
	StartPos, EndPos *Position
}

func NewToken(t string, v interface{}, sp, ep *Position) *Token {
	token := &Token{
		Type: t,
		Value: v,
		StartPos: sp,
		EndPos: ep,
	}

	if ep == nil {
		endPos := *sp
		endPos.Index += 1
		endPos.Col += 1
		token.EndPos = &endPos
	}

	return token
}

func (t *Token) String() string {
	return fmt.Sprintf("[%v: %v]", t.Type, t.Value)
}
