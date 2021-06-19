package main

import (
	"strconv"
	"strings"
)

const Digits = "0123456789"

type Lexer struct {
	Text string
	Pos int
	CurrChar string
}

func NewLexer(text string) *Lexer {
	lexer := &Lexer{
		Text: text,
		Pos: -1,
	}

	lexer.Advance()

	return lexer
}

func (l *Lexer) Advance() {
	l.Pos += 1
	if len(l.Text) > l.Pos {
		l.CurrChar = l.Text[l.Pos:l.Pos + 1]
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
			addToken(NewToken(TTPlus, "+"))
		} else if l.CurrChar == "-" {
			addToken(NewToken(TTMinus, "-"))
		} else if l.CurrChar == "*" {
			addToken(NewToken(TTMul, "*"))
		} else if l.CurrChar == "/" {
			addToken(NewToken(TTDiv, "/"))
		} else if l.CurrChar == "(" {
			addToken(NewToken(TTLParen, "("))
		} else if l.CurrChar == ")" {
			addToken(NewToken(TTRParen, ")"))
		} else {
			return []*Token{}, NewIlligalCharError("'" + l.CurrChar + "'")
		}
	}

	return tokens, nil
}
