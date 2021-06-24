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
	AdvanceCount int
	LastAdvanceCount int
	ToReverseCount int
}

func NewParseResult() *ParseResult {
	pr := &ParseResult{
		AdvanceCount: 0,
		LastAdvanceCount: 0,
		ToReverseCount: 0,
	}
	return pr
}

func (pr *ParseResult) Register(res interface{}) interface{} {
	if r, ok := res.(*ParseResult); ok {
		pr.AdvanceCount += r.AdvanceCount
		pr.LastAdvanceCount = r.LastAdvanceCount

		if r.Error != nil {
			pr.Error = r.Error
		}
		return r.Node
	}
	return res
}

func (pr *ParseResult) RegisterAdvance() {
	pr.AdvanceCount += 1
	pr.LastAdvanceCount = 1
}

func (pr *ParseResult) TryRegister(res interface{}) interface{} {
	if r, ok := res.(*ParseResult); ok {
		if r.Error != nil {
			pr.ToReverseCount = r.AdvanceCount
			return nil
		}
		return pr.Register(r.Node)
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
	p.UpdateToken()
	return p.CurrToken
}

func (p *Parser) Reverse(amount int) *Token {
	p.TokenIndex -= amount
	p.UpdateToken()
	return p.CurrToken
}

func (p *Parser) UpdateToken() {
	if p.TokenIndex < len(p.Tokens) {
		p.CurrToken = p.Tokens[p.TokenIndex]
	}
}

func (p *Parser) Parse() *ParseResult {
	pr := p.Statements()

	if pr.Error == nil && p.CurrToken.Type != TTEOF {
		return pr.Failure(
			NewInvalidSyntaxError(
				"Unexpected token",
				p.CurrToken.StartPos,
				p.CurrToken.EndPos))
	}

	return pr
}

func (p *Parser) SkipNewLines() *ParseResult {
	pr := NewParseResult()

	for p.CurrToken.Type == TTNewLine {
		pr.RegisterAdvance()
		p.Advance()
	}

	return pr
}

func (p *Parser) Statements() *ParseResult {
	pr := NewParseResult()

	stmts := []interface{}{}

	pr.Register(p.SkipNewLines())

	stmt := pr.Register(p.Statement())
	if pr.Error != nil {
		return pr
	}
	stmts = append(stmts, stmt)

	more := true
	for {
		lines := 1

		for p.CurrToken.Type == TTNewLine {
			pr.RegisterAdvance()
			p.Advance()
			lines += 1
		}

		if lines == 0 {
			more = false
		}
		
		if !more {
			break
		}

		stmt := pr.TryRegister(p.Statement())
		if stmt == nil {
			p.Reverse(pr.ToReverseCount)
			more = false
			continue
		}
		stmts = append(stmts, stmt)
	}

	return pr.Success(NewListNode(stmts))
}

func (p *Parser) Statement() *ParseResult {
	pr := NewParseResult()

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "return" {
		pr.RegisterAdvance()
		p.Advance()

		exp := pr.TryRegister(p.Exp())

		if exp == nil {
			pr.Register(pr.ToReverseCount)
		}

		return pr.Success(NewReturnNode(exp))
	}

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "continue" {
		pr.RegisterAdvance()
		p.Advance()
		return pr.Success(NewContinueNode())
	}

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "break" {
		pr.RegisterAdvance()
		p.Advance()
		return pr.Success(NewBreakNode())
	}

	exp := pr.Register(p.Exp())

	if pr.Error != nil {
		return pr
	}

	return pr.Success(exp)
}

func (p *Parser) IfExp() *ParseResult {
	pr := NewParseResult()
	cases := [][2]interface{}{}
	var elseCase interface{}

	if p.CurrToken.Type != TTKeyword || p.CurrToken.Value != "if" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected 'if'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	cond := pr.Register(p.Exp())
	if pr.Error != nil {
		return pr
	}

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "{" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected '{'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	stmts := pr.Register(p.Statements())
	if pr.Error != nil {
		return pr
	}

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "}" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected '}'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	cases = append(cases, [2]interface{}{cond, stmts})

	for p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "elif" {
		pr.RegisterAdvance()
		p.Advance()

		cond := pr.Register(p.Exp())
		if pr.Error != nil {
			return pr
		}

		pr.Register(p.SkipNewLines())

		if p.CurrToken.Type != TTOp || p.CurrToken.Value != "{" {
			return pr.Failure(
				NewInvalidSyntaxError("Expected '{'",
				p.CurrToken.StartPos,
				p.CurrToken.EndPos))
		}

		pr.RegisterAdvance()
		p.Advance()

		stmts := pr.Register(p.Statements())
		if pr.Error != nil {
			return pr
		}

		pr.Register(p.SkipNewLines())

		if p.CurrToken.Type != TTOp || p.CurrToken.Value != "}" {
			return pr.Failure(
				NewInvalidSyntaxError("Expected '}'",
				p.CurrToken.StartPos,
				p.CurrToken.EndPos))
		}

		pr.RegisterAdvance()
		p.Advance()

		cases = append(cases, [2]interface{}{cond, stmts})
	}

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "else" {
		pr.RegisterAdvance()
		p.Advance()

		pr.Register(p.SkipNewLines())

		if p.CurrToken.Type != TTOp || p.CurrToken.Value != "{" {
			return pr.Failure(
				NewInvalidSyntaxError("Expected '{'",
				p.CurrToken.StartPos,
				p.CurrToken.EndPos))
		}

		pr.RegisterAdvance()
		p.Advance()

		stmts := pr.Register(p.Statements())
		if pr.Error != nil {
			return pr
		}

		pr.Register(p.SkipNewLines())
		
		if p.CurrToken.Type != TTOp || p.CurrToken.Value != "}" {
			return pr.Failure(
				NewInvalidSyntaxError("Expected '}'",
				p.CurrToken.StartPos,
				p.CurrToken.EndPos))
		}

		pr.RegisterAdvance()
		p.Advance()

		elseCase = stmts
	}

	return pr.Success(NewIfNode(cases, elseCase))
}

func (p *Parser) WhileExp() *ParseResult {
	pr := NewParseResult()
	
	if p.CurrToken.Type != TTKeyword || p.CurrToken.Value != "while" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected 'while'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	cond := pr.Register(p.Exp())

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "{" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected '{'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	stmts := pr.Register(p.Statements())

	if pr.Error != nil {
		return pr
	}

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "}" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected '}'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	return pr.Success(NewWhileNode(cond, stmts))
}

func (p *Parser) ForExp() *ParseResult {
	pr := NewParseResult()

	if p.CurrToken.Type != TTKeyword || p.CurrToken.Value != "for" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected 'for'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTId {
		return pr.Failure(
			NewInvalidSyntaxError("Expected identifier",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	varName := p.CurrToken

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "=" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected '='",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	from := pr.Register(p.Exp())

	if pr.Error != nil {
		return pr
	}

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != ":" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected ':'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	to := pr.Register(p.Exp())

	if pr.Error != nil {
		return pr
	}

	pr.Register(p.SkipNewLines())

	var by interface{} = nil

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "by" {
		pr.RegisterAdvance()
		p.Advance()

		by = pr.Register(p.Exp())

		if pr.Error != nil {
			return pr
		}

		pr.Register(p.SkipNewLines())
	}


	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "{" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected '{'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	body := pr.Register(p.Statements())
	
	if pr.Error != nil {
		return pr
	}
	
	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "}" {
		return pr.Failure(
			NewInvalidSyntaxError("Expected '}'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	return pr.Success(NewForNode(varName, from, to, by, body))
}

func (p *Parser) FunDef() *ParseResult {
	pr := NewParseResult()

	if p.CurrToken.Type != TTKeyword || p.CurrToken.Value != "fun" {
		return pr.Failure(NewInvalidSyntaxError(
			"Expected 'fun'", p.CurrToken.StartPos, p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	name := ""
	args := []string{}

	if p.CurrToken.Type == TTId {
		name = p.CurrToken.Value.(string)

		pr.RegisterAdvance()
		p.Advance()
	}

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "(" {
		return pr.Failure(NewInvalidSyntaxError(
			"Expected identifier or '('", p.CurrToken.StartPos, p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type == TTId {
		args = append(args, p.CurrToken.Value.(string))

		pr.RegisterAdvance()
		p.Advance()

		pr.Register(p.SkipNewLines())

		for p.CurrToken.Type == TTOp && p.CurrToken.Value == "," {
			pr.RegisterAdvance()
			p.Advance()
			pr.Register(p.SkipNewLines())

			if p.CurrToken.Type != TTId {
				return pr.Failure(NewInvalidSyntaxError(
					"Expected identifer",
					p.CurrToken.StartPos,
					p.CurrToken.EndPos))
			}

			args = append(args, p.CurrToken.Value.(string))

			pr.RegisterAdvance()
			p.Advance()
			pr.Register(p.SkipNewLines())
		}
	}

	if p.CurrToken.Type == TTOp && p.CurrToken.Value == ")" {
		pr.RegisterAdvance()
		p.Advance()
		pr.Register(p.SkipNewLines())

		if p.CurrToken.Type == TTOp && p.CurrToken.Value == "=" {			
			pr.RegisterAdvance()
			p.Advance()
			pr.Register(p.SkipNewLines())
	
			body := pr.Register(p.Exp())
	
			if pr.Error != nil {
				return pr
			}

			return pr.Success(NewFunDefNode(name, args, body, true))
		} else if p.CurrToken.Type == TTOp && p.CurrToken.Value == "{" {
			pr.RegisterAdvance()
			p.Advance()
			pr.Register(p.SkipNewLines())
			
			stmts := pr.Register(p.Statements())

			if pr.Error != nil {
				return pr
			}

			pr.Register(p.SkipNewLines())

			if p.CurrToken.Type != TTOp || p.CurrToken.Value != "}" {
				return pr.Failure(
					NewInvalidSyntaxError("Expected '}'",
					p.CurrToken.StartPos,
					p.CurrToken.EndPos))
			}

			pr.RegisterAdvance()
			p.Advance()

			return pr.Success(NewFunDefNode(name, args, stmts, false))
		}

		return pr.Failure(NewInvalidSyntaxError(
			"Expected '{' or '='", p.CurrToken.StartPos, p.CurrToken.EndPos))
	}

	return pr.Failure(NewInvalidSyntaxError(
		"Expected identifier or ')'", p.CurrToken.StartPos, p.CurrToken.EndPos))
}

func (p *Parser) ListExp() *ParseResult {
	pr := NewParseResult()

	el := []interface{}{}

	if p.CurrToken.Type != TTOp || p.CurrToken.Value != "[" {
		return pr.Failure(NewInvalidSyntaxError(
			"Expected '['", p.CurrToken.StartPos, p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	pr.Register(p.SkipNewLines())

	if p.CurrToken.Type == TTOp && p.CurrToken.Value == "]" {
		pr.RegisterAdvance()
		p.Advance()
		return pr.Success(NewListNode(el))
	}

	el = append(el, pr.Register(p.Exp()))
	if pr.Error != nil {
		return pr
	}
	pr.Register(p.SkipNewLines())

	for p.CurrToken.Type == TTOp && p.CurrToken.Value == "," {
		pr.RegisterAdvance()
		p.Advance()
		pr.Register(p.SkipNewLines())

		el = append(el, pr.Register(p.Exp()))
		if pr.Error != nil {
			return pr
		}
		pr.Register(p.SkipNewLines())
	}

	if p.CurrToken.Type != TTOp && p.CurrToken.Value != "]" {
		return pr.Failure(NewInvalidSyntaxError(
			"Expected ']'",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	pr.RegisterAdvance()
	p.Advance()

	return pr.Success(NewListNode(el))
}

func (p *Parser) Call() *ParseResult {
	pr := NewParseResult()
	atom := pr.Register(p.Atom())

	if pr.Error != nil {
		return pr
	}

	if p.CurrToken.Type == TTOp && p.CurrToken.Value == "(" {
		pr.RegisterAdvance()
		p.Advance()
		pr.Register(p.SkipNewLines())

		args := []interface{}{}

		if p.CurrToken.Type == TTOp && p.CurrToken.Value == ")" {
			pr.RegisterAdvance()
			p.Advance()
			return pr.Success(NewFunCallNode(atom, args))
		} else {
			args = append(args, pr.Register(p.Exp()))
			if pr.Error != nil {
				return pr
			}
			pr.Register(p.SkipNewLines())

			for p.CurrToken.Type == TTOp && p.CurrToken.Value == "," {
				pr.RegisterAdvance()
				p.Advance()
				pr.Register(p.SkipNewLines())

				args = append(args, pr.Register(p.Exp()))
				if pr.Error != nil {
					return pr
				}
				pr.Register(p.SkipNewLines())
			}

			if p.CurrToken.Type != TTOp && p.CurrToken.Value != ")" {
				return pr.Failure(NewInvalidSyntaxError(
					"Expected ')'",
					p.CurrToken.StartPos,
					p.CurrToken.EndPos))
			}

			pr.RegisterAdvance()
			p.Advance()

			return pr.Success(NewFunCallNode(atom, args))
		}
	}

	return pr.Success(atom)
}

func (p *Parser) Atom() *ParseResult {
	pr := NewParseResult()
	t := p.CurrToken

	if t.Type == TTOp && t.Value == "(" {
		pr.RegisterAdvance()
		p.Advance()
		exp := pr.Register(p.Exp())
		if pr.Error != nil {
			return pr
		}
		if p.CurrToken.Type == TTOp && p.CurrToken.Value == ")" {
			pr.RegisterAdvance()
			p.Advance()
			return pr.Success(exp)
		} else {
			return pr.Failure(
				NewInvalidSyntaxError(
					"Expected ')'",
					p.CurrToken.StartPos,
					p.CurrToken.EndPos))
		}
	} else if t.Type == TTNum {
		pr.RegisterAdvance()
		p.Advance()
		return pr.Success(NewNumberNode(t))
	} else if t.Type == TTStr {
		pr.RegisterAdvance()
		p.Advance()
		return pr.Success(NewStringNode(t))
	} else if t.Type == TTNull {
		pr.RegisterAdvance()
		p.Advance()
		return pr.Success(NewNullNode(t))
	} else if t.Type == TTId {
		pr.RegisterAdvance()
		p.Advance()
		return pr.Success(NewVarAccessNode(t))
	} else if t.Type == TTOp && t.Value == "[" {
		list := pr.Register(p.ListExp())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(list)
	} else if t.Type == TTKeyword && t.Value == "if" {
		ifExp := pr.Register(p.IfExp())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(ifExp)
	} else if t.Type == TTKeyword && t.Value == "while" {
		whileExp := pr.Register(p.WhileExp())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(whileExp)
	} else if t.Type == TTKeyword && t.Value == "for" {
		forExp := pr.Register(p.ForExp())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(forExp)
	} else if t.Type == TTKeyword && t.Value == "fun" {
		funDef := pr.Register(p.FunDef())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(funDef)
	} else if t.Type == TTOp && t.Value == "?" {
		pr.RegisterAdvance()
		p.Advance()
	}

	return pr.Failure(NewInvalidSyntaxError(
		"Unexpected token",
		t.StartPos,
		t.EndPos))
}

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Call, p.Factor, TTOp, []string{"^"})
}

func (p *Parser) Factor() *ParseResult {
	pr := NewParseResult()
	t := p.CurrToken

	if t.Type == TTOp && t.Value == "+" || t.Value == "-" {
		pr.RegisterAdvance()
		p.Advance()
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
		pr.RegisterAdvance()
		p.Advance()

		varName := p.CurrToken

		if p.CurrToken.Type == TTId {
			pr.RegisterAdvance()
			p.Advance()
			
			if p.CurrToken.Type != TTOp || p.CurrToken.Value != "=" {
				return pr.Failure(NewInvalidSyntaxError(
					"Expected =", p.CurrToken.StartPos, p.CurrToken.EndPos))
			}

			pr.RegisterAdvance()
			p.Advance()
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
		return pr
	}

	if p.CurrToken.Type == TTOp && p.CurrToken.Value == "?" {
		pr.RegisterAdvance()
		p.Advance()

		left := pr.Register(p.Exp())
		if pr.Error != nil {
			return pr
		}

		if p.CurrToken.Type != TTOp || p.CurrToken.Value != ":" {
			return pr.Failure(
				NewInvalidSyntaxError(
					"Expected ':'",
					p.CurrToken.StartPos,
					p.CurrToken.EndPos))
		}

		pr.RegisterAdvance()
		p.Advance()

		right := pr.Register(p.Exp())
		if pr.Error != nil {
			return pr
		}

		return pr.Success(NewTernOpNode(node, left, right))
	}

	if pr.Error != nil {
		return pr.Failure(NewInvalidSyntaxError(
			"Unexpected token",
			p.CurrToken.StartPos,
			p.CurrToken.EndPos))
	}

	return pr.Success(node)
}

func (p *Parser) CompExp() *ParseResult {
	pr := NewParseResult()

	if p.CurrToken.Type == TTKeyword && p.CurrToken.Value == "not" {
		op := p.CurrToken
		pr.RegisterAdvance()
		p.Advance()
		node := pr.Register(p.CompExp())
		if pr.Error != nil {
			return pr
		}
		return pr.Success(NewUnaryOpNode(op, node))
	}

	node := pr.Register(p.BinOp(p.ArithExp, p.ArithExp, TTOp, []string{"==", "!=", ">", ">=", "<", "<="}))
	if pr.Error != nil {
		return pr
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

	pr.Register(p.SkipNewLines())

	for p.CurrToken.Type == opType && Contains(ops, p.CurrToken.Value) {
		op := p.CurrToken
		pr.RegisterAdvance()
		p.Advance()
		pr.Register(p.SkipNewLines())

		l := pr.Register(lf())

		if pr.Error != nil {
			return pr
		}

		r = NewBinNode(r, l, op)
	}

	return pr.Success(r)
}
