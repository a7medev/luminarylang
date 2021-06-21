package main

import "fmt"

type NumberNode struct {
	Token *Token
}

func NewNumberNode(t *Token) *NumberNode {
	n := &NumberNode{Token: t}
	return n
}

func (n *NumberNode) String() string {
	return n.Token.String()
}

type BinOpNode struct {
	Left interface{}
	Op *Token
	Right interface{}
}

func NewBinNode(l, r interface{}, o *Token) *BinOpNode {
	b := &BinOpNode{
		Left: l,
		Op: o,
		Right: r,
	}

	return b
}

func (b *BinOpNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", b.Left, b.Op, b.Right)
}

type UnaryOpNode struct {
	Op *Token
	Node interface{}
}

func (u *UnaryOpNode) String() string {
	return fmt.Sprintf("(%v, %v)", u.Op, u.Node)
}

func NewUnaryOpNode(o *Token, n interface{}) *UnaryOpNode {
	u := &UnaryOpNode{
		Op: o,
		Node: n,
	}

	return u
}

type VarAssignNode struct {
	NameToken *Token
	ValueNode interface{}
}

func NewVarAssignNode(n *Token, v interface{}) *VarAssignNode {
	va := &VarAssignNode{
		NameToken: n,
		ValueNode: v,
	}

	return va
}


type VarAccessNode struct {
	NameToken *Token
}

func NewVarAccessNode(n *Token) *VarAccessNode {
	va := &VarAccessNode{
		NameToken: n,
	}

	return va
}
