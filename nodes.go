package main

type NumberNode struct {
	Token *Token
}

func NewNumberNode(t *Token) *NumberNode {
	n := &NumberNode{Token: t}
	return n
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
