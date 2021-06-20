package main

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
		return NewNumber(val).SetPos(n.Token.Pos)
	} else {
		panic("Invalid number node")
	}
}

func (i *Interpretor) VisitBinOpNode(b *BinOpNode) *Number {
	l := i.Visit(b.Right)
	r := i.Visit(b.Left)

	switch b.Op.Value {
	case "+":
		return r.AddTo(l)
	case "-":
		return r.SubBy(l)
	case "*":
		return r.MulBy(l)
	case "/":
		return r.DivBy(l)
	case "%":
		return r.Mod(l)
	case "^":
		return r.Pow(l)
	default:
		return nil
	}
}

func (i *Interpretor) VisitUnaryOpNode(u *UnaryOpNode) *Number {
	n := i.Visit(u.Node)
	if u.Op.Value == "-" {
		return n.MulBy(NewNumber(-1))
	} else {
		return n
	}
}
