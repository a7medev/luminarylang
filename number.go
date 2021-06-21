package main

import (
	"fmt"
	"math"
)

type Number struct {
	Value float64
	StartPos, EndPos *Position
}

func NewNumber(v float64) *Number {
	n := &Number{Value: v}

	return n
}

func (n *Number) String() string {
	return fmt.Sprintf("%v", n.Value);
}

func (n *Number) SetPos(sp, ep *Position) *Number {
	n.StartPos = sp
	n.EndPos = ep
	if ep == nil {
		endPos := *sp
		endPos.Advance("")
		n.EndPos = &endPos
	}
	return n
}

func (n *Number) AddTo(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		return NewNumber(n.Value + o.Value), nil
	}
	return nil, NewInvalidSyntaxError("Expected a number", n.StartPos, nil)
}

func (n *Number) SubBy(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		return NewNumber(n.Value - o.Value), nil
	}
	return nil, NewInvalidSyntaxError("Expected a number", n.StartPos, nil)
}

func (n *Number) MulBy(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		return NewNumber(n.Value * o.Value), nil
	}
	return nil, NewInvalidSyntaxError("Expected a number", n.StartPos, nil)
}

func (n *Number) DivBy(other interface{}) (*Number, *Error) {
	switch o := other.(type) {
	case *Number:
		if o.Value != 0 {
			return NewNumber(n.Value / o.Value), nil
		} else {
			return nil, NewRuntimeError("Can't divide by zero", n.StartPos, o.EndPos)
		}
	}
	return nil, NewInvalidSyntaxError("Expected a number", n.StartPos, nil)
}

func (n *Number) Mod(other interface{}) (*Number, *Error) {
	switch o := other.(type) {
	case *Number:
		if o.Value != 0 {
			return NewNumber(math.Mod(n.Value, o.Value)), nil
		} else {
			return nil, NewRuntimeError("Can't divide by zero", n.StartPos, o.EndPos)
		}
	}
	return nil, NewInvalidSyntaxError("Expected a number", n.StartPos, nil)
}

func (n *Number) Pow(other interface{}) (*Number, *Error) {
	switch o := other.(type) {
	case *Number:
		return NewNumber(math.Pow(n.Value, o.Value)), nil
	}
	return nil, NewInvalidSyntaxError("Expected a number", n.StartPos, nil)
}

func (n *Number) IsEqualTo(other interface{}) *Number {
	if o, ok := other.(*Number); ok {
		var val float64 = 0
		if n.Value == o.Value {
			val = 1
		}
		return NewNumber(val)
	}

	return NewNumber(0)
}

func (n *Number) IsNotEqualTo(other interface{}) *Number {
	isEq := n.IsEqualTo(other)

	if isEq.Value == 1 {
		return NewNumber(0)
	}

	return NewNumber(1)
}

func (n *Number) IsGreaterThan(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		var val float64 = 0
		if n.Value > o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", n.StartPos, nil)
}

func (n *Number) IsGreaterThanOrEqual(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		var val float64 = 0
		if n.Value >= o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", n.StartPos, nil)
}

func (n *Number) IsLessThan(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		var val float64 = 0
		if n.Value < o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", n.StartPos, nil)
}

func (n *Number) IsLessThanOrEqual(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		var val float64 = 0
		if n.Value <= o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", n.StartPos, nil)
}

func (n *Number) And(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		var val float64 = 0
		if n.Value >= 1 && o.Value >= 1 {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", n.StartPos, nil)
}

func (n *Number) Or(other interface{}) (*Number, *Error) {
	if o, ok := other.(*Number); ok {
		var val float64 = 0
		if n.Value >= 1 || o.Value >= 1 {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", n.StartPos, nil)
}

func (n *Number) IsTrue() bool {
	return n.Value != 0
}
