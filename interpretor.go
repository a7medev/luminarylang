package main

import "fmt"

type Interpretor struct {}

func NewInterpretor() *Interpretor {
	i := &Interpretor{}
	return i
}

func (i *Interpretor) Visit(n interface{}) {
	if num, ok := n.(*NumberNode); ok {
		i.VisitNumberNode(num)
	} else if bin, ok := n.(*BinOpNode); ok {
		i.VisitBinOpNode(bin)
	} else if unary, ok := n.(*UnaryOpNode); ok {
		i.VisitUnaryOpNode(unary)
	} else {
		panic("no visit method for this node")
	}
}

func (i *Interpretor) VisitNumberNode(n *NumberNode) {
	fmt.Println("Visiting NumberNode", n)

}

func (i *Interpretor) VisitBinOpNode(b *BinOpNode) {
	fmt.Println("Visiting BinOpNode", b)
	i.Visit(b.Right)
	i.Visit(b.Left)

}

func (i *Interpretor) VisitUnaryOpNode(u *UnaryOpNode) {
	fmt.Println("Visiting UnaryOpNode", u)
	i.Visit(u.Node)

}
