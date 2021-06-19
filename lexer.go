package main

import (
	"strconv"
	"strings"
)

const Digits = "0123456789"

type Lexer struct {
	CurrChar, Text,	FileName,	FileText string
	Pos *Position
}

func NewLexer(txt, fn, ftxt string) *Lexer {
	lexer := &Lexer{
		Text: txt,
		FileName: fn,
		FileText: ftxt,
		Pos: NewPosition(-1, 1, 0, fn, ftxt),
	}

	lexer.Advance()

	return lexer
}

func (l *Lexer) Advance() {
	l.Pos.Advance(l.CurrChar)
	if len(l.Text) > l.Pos.Index {
		l.CurrChar = l.Text[l.Pos.Index:l.Pos.Index + 1]
	} else {
		l.CurrChar = ""
	}
}

func (l *Lexer) MakeNumber() *Token {
	numStr := ""
	hasDot := false
	for l.CurrChar != "" && strings.Contains(Digits + ".", l.CurrChar) {
		if l.CurrChar == "." {
			if hasDot {
				break
			}
			hasDot = true
		}
		numStr += l.CurrChar
		l.Advance()
	}
	if hasDot {
		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			panic(err)
		}
		return NewToken(TTFloat, val)
	}
	val, err := strconv.ParseInt(numStr, 10, 32)
	if err != nil {
		panic(err)
	}
	return NewToken(TTInt, val)
}

func (l *Lexer) MakeTokens() ([]*Token, *Error) {
	tokens := []*Token{}

	addToken  := func(t *Token) {
		tokens = append(tokens, t)
		l.Advance()
	}

	for l.CurrChar != "" {
		if strings.Contains("\t\n ", l.CurrChar) {
			l.Advance()
		} else if strings.Contains(Digits, l.CurrChar) {
			addToken(l.MakeNumber())
		} else if l.CurrChar == "+" {
			addToken(NewToken(TTOp, "+"))
		} else if l.CurrChar == "-" {
			addToken(NewToken(TTOp, "-"))
		} else if l.CurrChar == "*" {
			addToken(NewToken(TTOp, "*"))
		} else if l.CurrChar == "/" {
			addToken(NewToken(TTOp, "/"))
		} else if l.CurrChar == "(" {
			addToken(NewToken(TTParen, "("))
		} else if l.CurrChar == ")" {
			addToken(NewToken(TTParen, ")"))
		} else {
			return []*Token{}, NewIlligalCharError("'" + l.CurrChar + "'", l.Pos)
		}
	}

	return tokens, nil
}
