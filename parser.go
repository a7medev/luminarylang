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

func (p *Parser) Advance() *Token {
	p.TokenIndex += 1
	if p.TokenIndex < len(p.Tokens) {
		p.CurrToken = p.Tokens[p.TokenIndex]
	}
	return p.CurrToken
}

func (p *Parser) Parse() interface{} {
	res := p.Expr()
	return res
}

func (p *Parser) Factor() interface{} {
	t := p.CurrToken
	if t.Type == TTInt || t.Type == TTFloat {
		p.Advance()
		return NewNumberNode(t)
	}
	return nil
}

func (p *Parser) Term() interface{} {
	return p.BinOp(p.Factor, [2]string{"*", "/"});
}

func (p *Parser) Expr() interface{} {
	return p.BinOp(p.Term, [2]string{"+", "-"});
}

func (p *Parser) BinOp(f func() interface{}, ops [2]string) interface{} {
  var l interface{} = f()

	for p.CurrToken.Type == TTOp && (p.CurrToken.Value == ops[0] || p.CurrToken.Value == ops[1]) {
		op := p.CurrToken
		p.Advance()
		r := f()
		l = NewBinNode(l, r, op)
	}

	return l;
}
