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

func (i *Interpretor) Visit(n interface{}, ctx *Context) interface{} {
	if num, ok := n.(*NumberNode); ok {
		return i.VisitNumberNode(num, ctx)
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
	} else {
		panic("no visit method for this node")
	}
}

func (i *Interpretor) VisitNumberNode(n *NumberNode, ctx *Context) *Number {
	if val, ok := n.Token.Value.(float64); ok {
		return NewNumber(val).SetPos(n.Token.StartPos, n.Token.EndPos)
	} else {
		panic("Invalid number node")
	}
}

func (i *Interpretor) VisitBinOpNode(b *BinOpNode, ctx *Context) *Number {
	l := i.Visit(b.Right, ctx).(*Number)
	r := i.Visit(b.Left, ctx).(*Number)

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
	default:
		return nil
	}
}

func (i *Interpretor) VisitUnaryOpNode(u *UnaryOpNode, ctx *Context) *Number {
	n := i.Visit(u.Node, ctx).(*Number)
	if u.Op.Value == "-" {
		res, err := n.MulBy(NewNumber(-1))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return res
	} else {
		return n
	}
}

func (i *Interpretor) VisitVarAssignNode(va *VarAssignNode, ctx *Context) interface{} {
	num := i.Visit(va.ValueNode, ctx)
	return ctx.SymbolTable.Set(va.NameToken.Value.(string), num)
}

func (i *Interpretor) VisitVarAccessNode(va *VarAccessNode, ctx *Context) (interface{}, *Error) {
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

func (i *Interpretor) VisitIfNode(ifN *IfNode, ctx *Context) interface{} {
	for _, cs := range ifN.Cases {
		cond := cs[0]
		condVal := i.Visit(cond, ctx)
		
		if condVal.(*Number).IsTrue() {
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
