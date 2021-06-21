package main

import (
	"strconv"
	"strings"
)

const Digits = "0123456789"
const Letters = "abcdefghijklmnopqrstunwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const IdAllowedChars = Letters + Digits + "_"

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

func (l *Lexer) MakeId() *Token {
	idStr := ""

	for l.CurrChar != "" && strings.Contains(IdAllowedChars, l.CurrChar) {
		idStr += l.CurrChar
		l.Advance()
	}

	endPos := *l.Pos
	endPos.Index += len(idStr)
	endPos.Col += len(idStr)
	return NewToken(TTId, idStr, l.Pos, &endPos)
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

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		panic(err)
	}
	endPos := l.Pos
	endPos.Col += len(numStr)
	return NewToken(TTNum, val, l.Pos, endPos)
}

func (l *Lexer) MakeTokens() ([]*Token, *Error) {
	tokens := []*Token{}

	addToken  := func(t *Token) {
		tokens = append(tokens, t)
		// don't advance if token is a number cuz the MakeNumber method already advances
		if t.Type != TTNum && t.Type != TTId {
			l.Advance()
		}
	}

	for l.CurrChar != "" {
		if strings.Contains("\t\n ", l.CurrChar) {
			l.Advance()
		} else if strings.Contains(Letters, l.CurrChar) {
			addToken(l.MakeId())
		} else if strings.Contains(Digits, l.CurrChar) {
			addToken(l.MakeNumber())
		} else if l.CurrChar == "+" {
			addToken(NewToken(TTOp, "+", l.Pos, nil))
		} else if l.CurrChar == "-" {
			addToken(NewToken(TTOp, "-", l.Pos, nil))
		} else if l.CurrChar == "*" {
			addToken(NewToken(TTOp, "*", l.Pos, nil))
		} else if l.CurrChar == "/" {
			addToken(NewToken(TTOp, "/", l.Pos, nil))
		} else if l.CurrChar == "%" {
			addToken(NewToken(TTOp, "%", l.Pos, nil))
		} else if l.CurrChar == "^" {
			addToken(NewToken(TTOp, "^", l.Pos, nil))
		} else if l.CurrChar == "(" {
			addToken(NewToken(TTParen, "(", l.Pos, nil))
		} else if l.CurrChar == ")" {
			addToken(NewToken(TTParen, ")", l.Pos, nil))
		} else {
			endPos := *l.Pos
			endPos.Advance(l.CurrChar)
			text := l.Text[l.Pos.Index:endPos.Index]
			return []*Token{}, NewIlligalCharError("'" + text + "'", l.Pos, &endPos)
		}
	}

	tokens = append(tokens, NewToken(TTEOF, TTEOF, l.Pos, nil))

	return tokens, nil
}
