package main

import (
	"fmt"
)

type String struct {
	Value string
	StartPos, EndPos *Position
}

func NewString(v string) Value {
	s := &String{Value: v}

	return s
}

func (s *String) String() string {
	return fmt.Sprintf("%v", s.Value);
}

func (s *String) SetPos(sp, ep *Position) Value {
	s.StartPos = sp
	s.EndPos = ep
	if ep == nil {
		endPos := *sp
		endPos.Advance("")
		s.EndPos = &endPos
	}
	return s
}

func (s *String) AddTo(other interface{}) (Value, *Error) {
	switch o := other.(type) {
	case *String:
		return NewString(s.Value + o.Value), nil
	case *Number:
		return NewString(fmt.Sprintf("%v%v", s.Value, o.Value)), nil
	default:
		return nil, NewInvalidSyntaxError("Expected a number", s.StartPos, s.EndPos)
	}
}

func (s *String) SubBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '-' operation on a string", s.StartPos, s.EndPos)
}

func (s *String) MulBy(other interface{}) (Value, *Error) {
	if o, ok := other.(*Number); ok {
		str := ""
		for i := .0; i < o.Value; i++ {
			str += s.Value
		}
		return NewString(str), nil
	}
	return nil, NewInvalidSyntaxError("Expected a number", s.StartPos, s.EndPos)
}

func (s *String) DivBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '/' operation on a string", s.StartPos, s.EndPos)
}

func (s *String) Mod(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '%' operation on a string", s.StartPos, s.EndPos)
}

func (s *String) Pow(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '^' operation on a string", s.StartPos, s.EndPos)
}

func (s *String) IsEqualTo(other interface{}) Value {
	if o, ok := other.(*String); ok {
		var val float64 = 0
		if s.Value == o.Value {
			val = 1
		}
		return NewNumber(val)
	}

	return NewNumber(0)
}

func (s *String) IsNotEqualTo(other interface{}) Value {
	isEq := s.IsEqualTo(other)

	if isEq.GetVal() == 1 {
		return NewNumber(0)
	}

	return NewNumber(1)
}

func (s *String) IsGreaterThan(other interface{}) (Value, *Error) {
	if o, ok := other.(*String); ok {
		var val float64 = 0
		if s.Value > o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", s.StartPos, nil)
}

func (s *String) IsGreaterThanOrEqual(other interface{}) (Value, *Error) {
	if o, ok := other.(*String); ok {
		var val float64 = 0
		if s.Value >= o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", s.StartPos, nil)
}

func (s *String) IsLessThan(other interface{}) (Value, *Error) {
	if o, ok := other.(*String); ok {
		var val float64 = 0
		if s.Value < o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", s.StartPos, nil)
}

func (s *String) IsLessThanOrEqual(other interface{}) (Value, *Error) {
	if o, ok := other.(*String); ok {
		var val float64 = 0
		if s.Value <= o.Value {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", s.StartPos, nil)
}

func (s *String) And(other interface{}) (Value, *Error) {
	if o, ok := other.(*String); ok {
		var val float64 = 0
		if s.IsTrue() && o.IsTrue() {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", s.StartPos, nil)
}

func (s *String) Or(other interface{}) (Value, *Error) {
	if o, ok := other.(*String); ok {
		var val float64 = 0
		if s.IsTrue() || o.IsTrue() {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", s.StartPos, nil)
}

func (s *String) Not() Value {
	if s.IsTrue() {
		return NewNumber(0)
	}
	return NewNumber(1)
}

func (s *String) IsTrue() bool {
	return s.Value != ""
}

func (s *String) GetVal() interface{} {
	return s.Value
}
