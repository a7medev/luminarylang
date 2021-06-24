package main

import "fmt"

type NumberNode struct {
	Token *Token
}

func NewNumberNode(t *Token) *NumberNode {
	n := &NumberNode{Token: t}
	return n
}

type StringNode struct {
	Token *Token
}

func NewStringNode(t *Token) *StringNode {
	n := &StringNode{Token: t}
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

type IfNode struct {
	Cases [][2]interface{}
	ElseCase interface{}
}

func NewIfNode(cases [][2]interface{}, elseCase interface{}) *IfNode {
	n := &IfNode{
		Cases: cases,
		ElseCase: elseCase,
	}

	return n
}

type WhileNode struct {
	Cond, Exp interface{}
}

func NewWhileNode(c, e interface{}) *WhileNode {
	w := &WhileNode{
		Cond: c,
		Exp: e,
	}

	return w
}

type ForNode struct {
	Var *Token
	From, To, By, Body interface{}
}

func NewForNode(v *Token, f, t, b, bd interface{}) *ForNode {
	fo := &ForNode{
		Var: v,
		From: f,
		To: t,
		By: b,
		Body: bd,
	}

	return fo
}

type TernOpNode struct {
	Cond, Left, Right interface{}
}

func NewTernOpNode(c, l, r interface{}) *TernOpNode {
	t := &TernOpNode{
		Cond: c,
		Left: l,
		Right: r,
	}

	return t
}

type FunDefNode struct {
	Name string
	ArgNames []string
	Body interface{}
}

func NewFunDefNode(n string, a []string, b interface{}) *FunDefNode {
	f := &FunDefNode{
		Name: n,
		ArgNames: a,
		Body: b,
	}

	return f
}

type FunCallNode struct {
	Name interface{}
	Args []interface{}
}

func NewFunCallNode(n interface{}, a []interface{}) *FunCallNode {
	f := &FunCallNode{
		Name: n,
		Args: a,
	}

	return f
}

type ListNode struct {
	Elements []interface{}
}

func NewListNode(el []interface{}) *ListNode {
	l := &ListNode{Elements: el}
	return l
}
