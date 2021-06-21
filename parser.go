package main

type Parser struct {
	Tokens []*Token
	TokenIndex int
	CurrToken *Token
}

func NewParser(t []*Token, i int) *Parser {
	p := &Parser{
		Tokens: t,
		TokenIndex: i,
	}

	p.Advance()

	return p
}

type ParseResult struct {
	Error *Error
	Node interface{}
}

func NewParseResult() *ParseResult {
	pr := &ParseResult{}
	return pr
}

func (pr *ParseResult) Register(res interface{}) interface{} {
	if r, ok := res.(*ParseResult); ok {
		if r.Error != nil {
			pr.Error = r.Error
		}
		return r.Node
	}
	return res
}

func (pr *ParseResult) Success(node interface{}) *ParseResult {
	pr.Node = node
	return pr
}

func (pr *ParseResult) Failure(err *Error) *ParseResult {
	pr.Error = err
	return pr
}

func (p *Parser) Advance() *Token {
	p.TokenIndex += 1
	if p.TokenIndex < len(p.Tokens) {
		p.CurrToken = p.Tokens[p.TokenIndex]
	}
	return p.CurrToken
}

func (p *Parser) Parse() *ParseResult {
	pr := p.Exp()

	if pr.Error == nil && p.CurrToken.Type != TTEOF {
		return pr.Failure(
			NewInvalidSyntaxError(
				"Expected '+', '-', '*', '/', '%' or '^'",
				p.CurrToken.StartPos,
				p.CurrToken.EndPos))
	}

	return pr
}

func (p *Parser) Atom() *ParseResult {
	pr := NewParseResult()
	t := p.CurrToken

	if t.Type == TTOp && t.Value == "(" {
		pr.Register(p.Advance())
		exp := pr.Register(p.Exp())
		if pr.Error != nil {
			return pr
		}
		if p.CurrToken.Type == TTOp && p.CurrToken.Value == ")" {
			pr.Register(p.Advance())
			return pr.Success(exp)
		} else {
			return pr.Failure(
				NewInvalidSyntaxError(
					"Expected ')'",
					p.CurrToken.StartPos,
					p.CurrToken.EndPos))
		}
	} else if t.Type == TTNum {
		pr.Register(p.Advance())
		return pr.Success(NewNumberNode(t))
	} else if t.Type == TTId {
		pr.Register(p.Advance())
		return pr.Success(NewVarAccessNode(t))
	}

	return pr.Failure(NewInvalidSyntaxError(
		"Expected int, float, identifier, '+', '-' or '('",
		t.StartPos,
		t.EndPos))
}

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Atom, p.Factor, TTOp, []string{"^"})
}

func (p *Parser) Factor() *ParseResult {
	pr := NewParseResult()
	t := p.CurrToken

	if t.Type == TTOp && t.Value == "+" || t.Value == "-" {
		pr.Register(p.Advance())
		fc := pr.Register(p.Factor())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(NewUnaryOpNode(t, fc))
	}

	return p.Power()
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, p.Factor, TTOp, []string{"*", "/", "%"})
}

func (p *Parser) Exp() *ParseResult {
	pr := NewParseResult()

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "set" {
		pr.Register(p.Advance())

		varName := p.CurrToken

		if p.CurrToken.Type == TTId {
			pr.Register(p.Advance())
			
			if p.CurrToken.Type != TTOp || p.CurrToken.Value != "=" {
				return pr.Failure(NewInvalidSyntaxError(
					"Expected =", p.CurrToken.StartPos, p.CurrToken.EndPos))
			}

			pr.Register(p.Advance())
			exp := pr.Register(p.Exp())

			if pr.Error != nil {
				return pr
			}

			return pr.Success(NewVarAssignNode(varName, exp))
		}
		return pr.Failure(NewInvalidSyntaxError(
			"Expected identifier", p.CurrToken.StartPos, p.CurrToken.EndPos))
	}

	node := pr.Register(p.BinOp(p.CompExp, p.CompExp, TTKeyword, []string{"and", "or"}))

	if pr.Error != nil {
		return pr.Failure(NewInvalidSyntaxError(
			"Expected 'set', number, identifier, '+', '-', '(' or 'not'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	return pr.Success(node)
}

func (p *Parser) CompExp() *ParseResult {
	pr := NewParseResult()

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "not" {
		op := p.CurrToken
		pr.Register(p.Advance())
		node := pr.Register(p.CompExp())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(NewUnaryOpNode(op, node))
	}

	node := pr.Register(p.BinOp(p.ArithExp, p.ArithExp, TTOp, []string{"==", "!=", ">", ">=", "<", "<="}))
	if pr.Error != nil {
		return pr.Failure(NewInvalidSyntaxError(
			"Expected number, identifier, '+', '-', '(' or 'not'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}
	return pr.Success(node)
}

func (p *Parser) ArithExp() *ParseResult {
	return p.BinOp(p.Term, p.Term, TTOp, []string{"+", "-"})
}

func (p *Parser) BinOp(rf, lf func() *ParseResult, opType string, ops []string) *ParseResult {
	pr := NewParseResult()
  r := pr.Register(rf())

	if pr.Error != nil {
		return pr
	}

	for p.CurrToken.Type == opType && Contains(ops, p.CurrToken.Value) {
		op := p.CurrToken
		pr.Register(p.Advance())
		l := pr.Register(lf())

		if pr.Error != nil {
			return pr
		}

		r = NewBinNode(r, l, op)
	}

	return pr.Success(r)
}
