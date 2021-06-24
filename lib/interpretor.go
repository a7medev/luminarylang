package main

type Interpretor struct {}

func NewInterpretor() *Interpretor {
	i := &Interpretor{}
	return i
}

type RuntimeResult struct {
	Value
	FunReturnValue Value
	ContinueLoop bool
	BreakLoop bool
	Error *Error
}

func NewRuntimeResult() *RuntimeResult {
	r := &RuntimeResult{}
	r.Reset()
	return r
}

func (rr *RuntimeResult) Reset() {
	rr.BreakLoop = false
	rr.ContinueLoop = false
	rr.FunReturnValue = nil
	rr.Error = nil
	rr.Value = nil
}

func (rr *RuntimeResult) Register(res interface{}) Value {
	if r, ok := res.(*RuntimeResult); ok {
		if r.Error != nil {
			rr.Error = r.Error
		}
		rr.ContinueLoop = r.ContinueLoop
		rr.FunReturnValue = r.FunReturnValue
		rr.BreakLoop = r.BreakLoop
		return r.Value
	} else if v, ok := res.(Value); ok {
		return v
	}
	return nil
}

func (rr *RuntimeResult) Success(val Value) *RuntimeResult {
	rr.Reset()
	rr.Value = val
	return rr
}

func (rr *RuntimeResult) SuccessReturn(val Value) *RuntimeResult {
	rr.Reset()
	rr.FunReturnValue = val
	return rr
}

func (rr *RuntimeResult) SuccessContinue() *RuntimeResult {
	rr.Reset()
	rr.ContinueLoop = true
	return rr
}

func (rr *RuntimeResult) SuccessBreak() *RuntimeResult {
	rr.Reset()
	rr.BreakLoop = true
	return rr
}

func (rr *RuntimeResult) Failure(err *Error) *RuntimeResult {
	rr.Reset()
	rr.Error = err
	return rr
}

func (rr *RuntimeResult) ShouldReturn() bool {
	return rr.Error != nil || rr.FunReturnValue != nil || rr.ContinueLoop || rr.BreakLoop
}

func (i *Interpretor) Visit(n interface{}, ctx *Context) *RuntimeResult {
	if num, ok := n.(*NumberNode); ok {
		return i.VisitNumberNode(num, ctx)
	} else if str, ok := n.(*StringNode); ok {
		return i.VisitStringNode(str, ctx)
	} else if null, ok := n.(*NullNode); ok {
		return i.VisitNullNode(null, ctx)
	} else if tern, ok := n.(*TernOpNode); ok {
		return i.VisitTernOpNode(tern, ctx)
	} else if bin, ok := n.(*BinOpNode); ok {
		return i.VisitBinOpNode(bin, ctx)
	} else if unary, ok := n.(*UnaryOpNode); ok {
		return i.VisitUnaryOpNode(unary, ctx)
	} else if list, ok := n.(*ListNode); ok {
		return i.VisitListNode(list, ctx)
	} else if access, ok := n.(*VarAccessNode); ok {
		return i.VisitVarAccessNode(access, ctx)
	} else if elAccess, ok := n.(*ElementAccessNode); ok {
		return i.VisitElementAccessNode(elAccess, ctx)
	} else if assign, ok := n.(*VarAssignNode); ok {
		return i.VisitVarAssignNode(assign, ctx)
	} else if ifN, ok := n.(*IfNode); ok {
		return i.VisitIfNode(ifN, ctx)
	} else if forN, ok := n.(*ForNode); ok {
		return i.VisitForNode(forN, ctx)
	} else if while, ok := n.(*WhileNode); ok {
		return i.VisitWhileNode(while, ctx)
	} else if contin, ok := n.(*ContinueNode); ok {
		return i.VisitContinueNode(contin, ctx)
	} else if brk, ok := n.(*BreakNode); ok {
		return i.VisitBreakNode(brk, ctx)
	} else if funDef, ok := n.(*FunDefNode); ok {
		return i.VisitFunDefNode(funDef, ctx)
	} else if funCall, ok := n.(*FunCallNode); ok {
		return i.VisitFunCallNode(funCall, ctx)
	} else if ret, ok := n.(*ReturnNode); ok {
		return i.VisitReturnNode(ret, ctx)
	} else {
		panic("no visit method for this node")
	}
}

func (i *Interpretor) VisitStringNode(s *StringNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	if val, ok := s.Token.Value.(string); ok {
		str := NewString(val).SetPos(s.Token.StartPos, s.Token.EndPos)
		return rr.Success(str)
	} else {
		return rr.Failure(NewRuntimeError("Invalid string node", s.Token.StartPos, s.Token.EndPos))
	}
}

func (i *Interpretor) VisitNullNode(s *NullNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	str := NewNull().SetPos(s.Token.StartPos, s.Token.EndPos)
	return rr.Success(str)
}

func (i *Interpretor) VisitNumberNode(n *NumberNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	if val, ok := n.Token.Value.(float64); ok {
		num := NewNumber(val).SetPos(n.Token.StartPos, n.Token.EndPos)
		return rr.Success(num)
	} else {
		return rr.Failure(NewRuntimeError("Invalid number node", n.Token.StartPos, n.Token.EndPos))
	}
}

func (i *Interpretor) VisitTernOpNode(t *TernOpNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	c := rr.Register(i.Visit(t.Cond, ctx))
	if rr.ShouldReturn() {
		return rr
	}

	if c.IsTrue() {
		l := rr.Register(i.Visit(t.Left, ctx))
		if rr.ShouldReturn() {
			return rr
		}
		return rr.Success(l)
	}
	r := rr.Register(i.Visit(t.Right, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	return rr.Success(r)
}

func (i *Interpretor) VisitBinOpNode(b *BinOpNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	l := rr.Register(i.Visit(b.Right, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	r := rr.Register(i.Visit(b.Left, ctx))
	if rr.ShouldReturn() {
		return rr
	}

	switch b.Op.Value {
	case "+":
		res, err := r.AddTo(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "-":
		res, err := r.SubBy(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "*":
		res, err := r.MulBy(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "/":
		res, err := r.DivBy(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "%":
		res, err := r.Mod(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "^":
		res, err := r.Pow(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "==":
		return rr.Success(r.IsEqualTo(l))
	case "!=":
		return rr.Success(r.IsNotEqualTo(l))
	case ">":
		res, err := r.IsGreaterThan(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case ">=":
		res, err := r.IsGreaterThanOrEqual(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "<":
		res, err := r.IsLessThan(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "<=":
		res, err := r.IsLessThanOrEqual(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "and":
		res, err := r.And(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	case "or":
		res, err := r.Or(l)
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	default:
		return rr.Failure(NewInvalidSyntaxError("Unexpected operator", nil, nil))
	}
}

func (i *Interpretor) VisitUnaryOpNode(u *UnaryOpNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	n := rr.Register(i.Visit(u.Node, ctx))
	if rr.ShouldReturn() {
		return rr
	}

	if u.Op.Value == "-" {
		res, err := n.MulBy(NewNumber(-1))
		if err != nil {
			return rr.Failure(err)
		}
		return rr.Success(res)
	} else if u.Op.Value == "not" {
		return rr.Success(n.Not())
	} else {
		return rr.Success(n)
	}
}

func (i *Interpretor) VisitVarAssignNode(va *VarAssignNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()
	num := rr.Register(i.Visit(va.ValueNode, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	return rr.Success(ctx.SymbolTable.Set(va.NameToken.Value.(string), num))
}

func (i *Interpretor) VisitVarAccessNode(va *VarAccessNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	name := va.NameToken.Value.(string)
	val := ctx.SymbolTable.Get(name)

	if val == nil {
		return rr.Success(NewNull())
	}
	return rr.Success(val)
}

func (i *Interpretor) VisitIfNode(ifN *IfNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	for _, cs := range ifN.Cases {
		cond := cs[0]
		condVal := rr.Register(i.Visit(cond, ctx))
		if rr.ShouldReturn() {
			return rr
		}

		if condVal.IsTrue() {
			exp := cs[1]
			expVal := rr.Register(i.Visit(exp, ctx))
			if rr.ShouldReturn() {
				return rr
			}
			return rr.Success(expVal)
		}
	}

	if ifN.ElseCase != nil {
		exp := ifN.ElseCase
		expVal := rr.Register(i.Visit(exp, ctx))
		if rr.ShouldReturn() {
			return rr
		}
		return rr.Success(expVal)
	}

	return rr
}

func (i *Interpretor) VisitWhileNode(w *WhileNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	cond := func() *RuntimeResult {
		val := rr.Register(i.Visit(w.Cond, ctx))
		if rr.ShouldReturn() {
			return rr
		}
		return rr.Success(val)
	}

	for {
		res := cond()
		if rr.ShouldReturn() {
			return rr
		}
		if !res.IsTrue() {
			break
		}
		rr.Register(i.Visit(w.Exp, ctx))
		if rr.ShouldReturn() && !rr.BreakLoop && !rr.ContinueLoop {
			return rr
		}
		if rr.BreakLoop {
			break
		}
		if rr.ContinueLoop {
			continue
		}
	}

	return rr
}

func (i *Interpretor) VisitForNode(f *ForNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	fromVal := rr.Register(i.Visit(f.From, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	from := fromVal.GetVal().(float64)
	toVal := rr.Register(i.Visit(f.To, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	to := toVal.GetVal().(float64)
	byVal := NewNumber(1)

	varName := f.Var.Value.(string)

	if f.By != nil {
		byVal = rr.Register(i.Visit(f.By, ctx))
		if rr.ShouldReturn() {
			return rr
		}
	}

	by := byVal.GetVal().(float64)

	cond := func() bool {
		if by > 0 {
			return from <= to
		} else {
			return from >= to
		}
	}

	for {
		if cond() {
			ctx.SymbolTable.Set(varName, NewNumber(from))
			from += by
			rr.Register(i.Visit(f.Body, ctx))
			if rr.ShouldReturn() && !rr.BreakLoop && !rr.ContinueLoop {
				return rr
			}
			if rr.BreakLoop {
				break
			}
			if rr.ContinueLoop {
				continue
			}
		} else {
			ctx.SymbolTable.Del(varName)
			break
		}
	}

	return rr
}

func (i *Interpretor) VisitContinueNode(r *ContinueNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()
	return rr.SuccessContinue()
}

func (i *Interpretor) VisitBreakNode(r *BreakNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()
	return rr.SuccessBreak()
}

func (i *Interpretor) VisitFunDefNode(f *FunDefNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	fun := NewFunction(f.Name, f.ArgNames, f.Body, f.ReturnBody)

	if f.Name != "" {
		ctx.SymbolTable.Set(f.Name, fun)
	}

	return rr.Success(fun)
}

func (i *Interpretor) VisitFunCallNode(f *FunCallNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	fun := rr.Register(i.Visit(f.Name, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	args := []interface{}{}

	for _, val := range f.Args {
		item := rr.Register(i.Visit(val, ctx))
		if rr.ShouldReturn() {
			return rr
		}
		args = append(args, item)
	}

	val := rr.Register(fun.Call(args, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	return rr.Success(val)
}

func (i *Interpretor) VisitReturnNode(r *ReturnNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	if r.Value != nil {
		val := rr.Register(i.Visit(r.Value, ctx))
		if rr.ShouldReturn() {
			return rr
		}
		return rr.SuccessReturn(val)
	}

	return rr.SuccessReturn(NewNull())
}

func (i *Interpretor) VisitListNode(f *ListNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()

	elements := []interface{}{}

	for _, el := range f.Elements {
		element := rr.Register(i.Visit(el, ctx))
		if rr.ShouldReturn() {
			return rr
		}
		elements = append(elements, element)
	}

	return rr.Success(NewList(elements))
}

func (i *Interpretor) VisitElementAccessNode(a *ElementAccessNode, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()
	list := rr.Register(i.Visit(a.Node, ctx))
	if rr.ShouldReturn() {
		return rr
	}
	index := rr.Register(i.Visit(a.Index, ctx))
	switch i := index.(type) {
	case *Number:
		el := rr.Register(list.AccessElement(int(i.Value), ctx))
		if rr.ShouldReturn() {
			return rr
		}
		return rr.Success(el)
	default:
		return rr.Failure(NewRuntimeError("Expected a number for the index", nil, nil))
	}
}
