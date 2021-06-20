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

func NewUnaryOpNode(o *Token, n interface{}) *UnaryOpNode {
	u := &UnaryOpNode{
		Op: o,
		Node: n,
	}

	return u
}
