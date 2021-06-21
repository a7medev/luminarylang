package main

import "fmt"

type Interpretor struct {}

func NewInterpretor() *Interpretor {
	i := &Interpretor{}
	return i
}

func (i *Interpretor) Visit(n interface{}) *Number {
	if num, ok := n.(*NumberNode); ok {
		return i.VisitNumberNode(num)
	} else if bin, ok := n.(*BinOpNode); ok {
		return i.VisitBinOpNode(bin)
	} else if unary, ok := n.(*UnaryOpNode); ok {
		return i.VisitUnaryOpNode(unary)
	} else {
		panic("no visit method for this node")
	}
}

func (i *Interpretor) VisitNumberNode(n *NumberNode) *Number {
	if val, ok := n.Token.Value.(float64); ok {
		return NewNumber(val).SetPos(n.Token.StartPos, n.Token.EndPos)
	} else {
		panic("Invalid number node")
	}
}

func (i *Interpretor) VisitBinOpNode(b *BinOpNode) *Number {
	l := i.Visit(b.Right)
	r := i.Visit(b.Left)

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

func (i *Interpretor) VisitUnaryOpNode(u *UnaryOpNode) *Number {
	n := i.Visit(u.Node)
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
