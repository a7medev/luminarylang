package main

import "fmt"

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
		return i.VisitVarAccessNode(access, ctx)
	} else if assign, ok := n.(*VarAssignNode); ok {
		return i.VisitVarAssignNode(assign, ctx)
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

func (i *Interpretor) VisitVarAccessNode(va *VarAccessNode, ctx *Context) interface{} {
	name := va.NameToken.Value.(string)
	val := ctx.SymbolTable.Get(name)
	if val == nil {
		return NewRuntimeError(
			"undefined variable '" + name + "'",
			va.NameToken.StartPos,
			va.NameToken.EndPos)
	}
	return val
}
