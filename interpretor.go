package main

import (
	"fmt"
	"os"
)

type Interpretor struct {}

func NewInterpretor() *Interpretor {
	i := &Interpretor{}
	return i
}

func (i *Interpretor) Visit(n interface{}, ctx *Context) Value {
	if num, ok := n.(*NumberNode); ok {
		return i.VisitNumberNode(num, ctx)
	} else if str, ok := n.(*StringNode); ok {
		return i.VisitStringNode(str, ctx)
	} else if tern, ok := n.(*TernOpNode); ok {
		return i.VisitTernOpNode(tern, ctx)
	} else if bin, ok := n.(*BinOpNode); ok {
		return i.VisitBinOpNode(bin, ctx)
	} else if unary, ok := n.(*UnaryOpNode); ok {
		return i.VisitUnaryOpNode(unary, ctx)
	} else if access, ok := n.(*VarAccessNode); ok {
		val, err := i.VisitVarAccessNode(access, ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
			return nil
		}
		return val
	} else if assign, ok := n.(*VarAssignNode); ok {
		return i.VisitVarAssignNode(assign, ctx)
	} else if ifN, ok := n.(*IfNode); ok {
		return i.VisitIfNode(ifN, ctx)
	} else if forN, ok := n.(*ForNode); ok {
		i.VisitForNode(forN, ctx)
		return nil
	} else if while, ok := n.(*WhileNode); ok {
		i.VisitWhileNode(while, ctx)
		return nil
	} else {
		panic("no visit method for this node")
	}
}

func (i *Interpretor) VisitStringNode(s *StringNode, ctx *Context) Value {
	if val, ok := s.Token.Value.(string); ok {
		return NewString(val).SetPos(s.Token.StartPos, s.Token.EndPos)
	} else {
		panic("Invalid string node")
	}
}

func (i *Interpretor) VisitNumberNode(n *NumberNode, ctx *Context) Value {
	if val, ok := n.Token.Value.(float64); ok {
		return NewNumber(val).SetPos(n.Token.StartPos, n.Token.EndPos)
	} else {
		panic("Invalid number node")
	}
}

func (i *Interpretor) VisitTernOpNode(t *TernOpNode, ctx *Context) Value {
	c := i.Visit(t.Cond, ctx)
	
	if c.IsTrue() {
		return i.Visit(t.Left, ctx)	
	}
	return i.Visit(t.Right, ctx)
}

func (i *Interpretor) VisitBinOpNode(b *BinOpNode, ctx *Context) Value {
	l := i.Visit(b.Right, ctx)
	r := i.Visit(b.Left, ctx)

	switch b.Op.Value {
	case "+":
		res, err := r.AddTo(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "-":
		res, err := r.SubBy(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "*":
		res, err := r.MulBy(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "/":
		res, err := r.DivBy(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "%":
		res, err := r.Mod(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "^":
		res, err := r.Pow(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "==":
		return r.IsEqualTo(l)
	case "!=":
		return r.IsNotEqualTo(l)
	case ">":
		res, err := r.IsGreaterThan(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case ">=":
		res, err := r.IsGreaterThanOrEqual(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "<":
		res, err := r.IsLessThan(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "<=":
		res, err := r.IsLessThanOrEqual(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "and":
		res, err := r.And(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	case "or":
		res, err := r.Or(l)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	default:
		return nil
	}
}

func (i *Interpretor) VisitUnaryOpNode(u *UnaryOpNode, ctx *Context) Value {
	n := i.Visit(u.Node, ctx)
	if u.Op.Value == "-" {
		res, err := n.MulBy(NewNumber(-1))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	} else if u.Op.Value == "not" {
		return n.Not()
	} else {
		return n
	}
}

func (i *Interpretor) VisitVarAssignNode(va *VarAssignNode, ctx *Context) Value {
	num := i.Visit(va.ValueNode, ctx)
	return ctx.SymbolTable.Set(va.NameToken.Value.(string), num)
}

func (i *Interpretor) VisitVarAccessNode(va *VarAccessNode, ctx *Context) (Value, *Error) {
	name := va.NameToken.Value.(string)
	val := ctx.SymbolTable.Get(name)
	if val == nil {
		return nil, NewRuntimeError(
			"undefined variable '" + name + "'",
			va.NameToken.StartPos,
			va.NameToken.EndPos)
	}
	return val, nil
}

func (i *Interpretor) VisitIfNode(ifN *IfNode, ctx *Context) Value {
	for _, cs := range ifN.Cases {
		cond := cs[0]
		condVal := i.Visit(cond, ctx)

		if condVal.IsTrue() {
			exp := cs[1]
			expVal := i.Visit(exp, ctx)
			return expVal
		}
	}

	if ifN.ElseCase != nil {
		exp := ifN.ElseCase
		expVal := i.Visit(exp, ctx)
		return expVal
	}

	return nil
}

func (i *Interpretor) VisitWhileNode(w *WhileNode, ctx *Context) {
	cond := func() bool {
		return i.Visit(w.Cond, ctx).IsTrue()
	}

	for {
		if !cond() {
			break
		}
		i.Visit(w.Exp, ctx)
	}
}

func (i *Interpretor) VisitForNode(f *ForNode, ctx *Context) {
	from := i.Visit(f.From, ctx).GetVal().(float64)
	to := i.Visit(f.To, ctx).GetVal().(float64)
	by := NewNumber(1)

	varName := f.Var.Value.(string)

	if f.By != nil {
		by = i.Visit(f.By, ctx)
	}

	byVal := by.GetVal().(float64)

	cond := func() bool {
		if byVal > 0 {
			return from <= to
		} else {
			return from >= to
		}
	}

	for {
		if cond() {
			ctx.SymbolTable.Set(varName, NewNumber(from))
			from += byVal
			i.Visit(f.Body, ctx)
			} else {
			ctx.SymbolTable.Del(varName)
			break
		}
	}
}
