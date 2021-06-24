package main

import (
	"strconv"
	"strings"
)

const Digits = "0123456789"
const Letters = "abcdefghijklmnopqrstunwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const IdAllowedChars = Letters + Digits + "_"

var Keywords = [11]string{"set", "and", "or", "not", "if", "else", "elif", "while", "for", "by", "fun"}

const SimpleOps = "+-*/%^(){}?:,[]"

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
	startPos := *l.Pos

	for l.CurrChar != "" && strings.Contains(IdAllowedChars, l.CurrChar) {
		idStr += l.CurrChar
		l.Advance()
	}

	endPos := *l.Pos
	if Contains(Keywords, idStr) {
		return NewToken(TTKeyword, idStr, &startPos, &endPos)
	}
	return NewToken(TTId, idStr, &startPos, &endPos)
}

func (l *Lexer) MakeNumber() *Token {
	numStr := ""
	hasDot := false
	startPos := *l.Pos

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
	endPos := *l.Pos
	return NewToken(TTNum, val, &startPos, &endPos)
}

func (l *Lexer) MakeString() (*Token, *Error) {
	str := ""
	escape := false
	startPos := *l.Pos

	escapeChars := map[string]string{
		"n": "\n",
		"t": "\t",
		"\"": "\"",
		"\\": "\\",
	}

	l.Advance()

	for l.CurrChar != "" && (l.CurrChar != "\"" || escape) {
		if escape {
			str += escapeChars[l.CurrChar]
			if str == "" {
				return nil, NewInvalidSyntaxError("Expected 'n' or 't' or '\"' after '\\'", &startPos, l.Pos)
			}
			escape = false
		} else if l.CurrChar == "\\" {
			escape = true
		} else {
			str += l.CurrChar
			escape = false
		}
		l.Advance()
	}

	if l.CurrChar != "\"" {
		return nil, NewInvalidSyntaxError("Expected '\"'", &startPos, l.Pos)
	}

	l.Advance()
	endPos := *l.Pos
	return NewToken(TTStr, str, &startPos, &endPos), nil
}

func (l *Lexer) MakeNotEquals() (*Token, *Error) {
	startPos := *l.Pos

	l.Advance()

	if l.CurrChar == "=" {
		l.Advance()
		return NewToken(TTOp, "!=", &startPos, l.Pos), nil
	}

	return nil, NewInvalidSyntaxError("Expected '=' after '!'", &startPos, l.Pos)
}

func (l *Lexer) MakeEquals() *Token {
	startPos := *l.Pos

	l.Advance()

	if l.CurrChar == "=" {
		l.Advance()
		return NewToken(TTOp, "==", &startPos, l.Pos)
	}

	return NewToken(TTOp, "=", &startPos, l.Pos)
}

func (l *Lexer) MakeGreaterThan() *Token {
	startPos := *l.Pos

	l.Advance()

	if l.CurrChar == "=" {
		l.Advance()
		return NewToken(TTOp, ">=", &startPos, l.Pos)
	}

	return NewToken(TTOp, ">", &startPos, l.Pos)
}

func (l *Lexer) MakeLessThan() *Token {
	startPos := *l.Pos

	l.Advance()

	if l.CurrChar == "=" {
		l.Advance()
		return NewToken(TTOp, "<=", &startPos, l.Pos)
	}

	return NewToken(TTOp, "<", &startPos, l.Pos)	
}

func (l *Lexer) MakeTokens() ([]*Token, *Error) {
	tokens := []*Token{}

	addToken  := func(t *Token, adv bool) {
		tokens = append(tokens, t)
		if adv {
			l.Advance()
		}
	}

	for l.CurrChar != "" {
		if strings.Contains("\t\n ", l.CurrChar) {
			l.Advance()
		} else if strings.Contains(Letters, l.CurrChar) {
			addToken(l.MakeId(), false)
		} else if strings.Contains(Digits, l.CurrChar) {
			addToken(l.MakeNumber(), false)
		} else if l.CurrChar == "\"" {
			tok, err := l.MakeString()
			if err != nil {
				return []*Token{}, err
			}
			addToken(tok, false)
		} else if strings.Contains(SimpleOps, l.CurrChar) {
			addToken(NewToken(TTOp, l.CurrChar, l.Pos, nil), true)
		} else if l.CurrChar == "!" {
			tok, err := l.MakeNotEquals()
			if err != nil {
				return []*Token{}, err
			}
			addToken(tok, false)
		} else if l.CurrChar == "=" {
			addToken(l.MakeEquals(), false)
		} else if l.CurrChar == ">" {
			addToken(l.MakeGreaterThan(), false)
		} else if l.CurrChar == "<" {
			addToken(l.MakeLessThan(), false)
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
