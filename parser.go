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
		return pr.Failure(NewInvalidSyntaxError("Expected '+', '-', '*' or '/'", p.CurrToken.Pos))
	}

	return pr
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
	} else if t.Type == TTParen && t.Value == "(" {
		pr.Register(p.Advance())
		exp := pr.Register(p.Exp())
		if pr.Error != nil {
			return pr
		}
		if p.CurrToken.Type == TTParen && p.CurrToken.Value == ")" {
			pr.Register(p.Advance())
			return pr.Success(exp)
		} else {
			return pr.Failure(NewInvalidSyntaxError("Expected ')'", p.CurrToken.Pos))
		}
	} else if t.Type == TTNum {
		pr.Register(p.Advance())
		return pr.Success(NewNumberNode(t))
	}
	return pr.Failure(NewInvalidSyntaxError("Expected int or float", t.Pos))
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []string{"*", "/", "%"})
}

func (p *Parser) Exp() *ParseResult {
	return p.BinOp(p.Term, []string{"+", "-"})
}

func (p *Parser) BinOp(f func() *ParseResult, ops []string) *ParseResult {
	pr := NewParseResult()
  l := pr.Register(f())

	if pr.Error != nil {
		return pr
	}

	for p.CurrToken.Type == TTOp && Contains(ops, p.CurrToken.Value) {
		op := p.CurrToken
		pr.Register(p.Advance())
		r := pr.Register(f())

		if pr.Error != nil {
			return pr
		}

		l = NewBinNode(l, r, op)
	}

	return pr.Success(l)
}
