package main

type Value interface {
	SetPos(sp, ep *Position) Value
	String() string
	AddTo(other interface{}) (Value, *Error)
	SubBy(interface{}) (Value, *Error)
	MulBy(interface{}) (Value, *Error)
	DivBy(interface{}) (Value, *Error)
	Mod(interface{}) (Value, *Error)
	Pow(interface{}) (Value, *Error)
	IsEqualTo(interface{}) Value
	IsNotEqualTo(interface{}) Value
	IsGreaterThan(interface{}) (Value, *Error)
	IsGreaterThanOrEqual(interface{}) (Value, *Error)
	IsLessThan(interface{}) (Value, *Error)
	IsLessThanOrEqual(interface{}) (Value, *Error)
	And(interface{}) (Value, *Error)
	Or(interface{}) (Value, *Error)
	Not() Value
	IsTrue() bool
	GetVal() interface{}
	Call(args []interface{}, ctx *Context) *RuntimeResult
}
