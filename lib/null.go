package main

type Null struct {
	StartPos, EndPos *Position
}

func NewNull() Value {
	n := &Null{}
	return n
}

func (n *Null) String() string {
	return "(null)"
}

func (n *Null) SetPos(sp, ep *Position) Value {
	n.StartPos = sp
	n.EndPos = ep
	if ep == nil {
		endPos := *sp
		endPos.Advance("")
		n.EndPos = &endPos
	}
	return n
}

func (n *Null) AddTo(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '+' operation on a null value", n.StartPos, n.EndPos)
}

func (n *Null) SubBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '-' operation on a null value", n.StartPos, n.EndPos)
}

func (n *Null) MulBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '*' operation on a null value", n.StartPos, n.EndPos)
}

func (n *Null) DivBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '/' operation on a null value", n.StartPos, n.EndPos)
}

func (n *Null) Mod(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '%' operation on a null value", n.StartPos, n.EndPos)
}

func (n *Null) Pow(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '^' operation on a null value", n.StartPos, n.EndPos)
}

func (n *Null) IsEqualTo(other interface{}) Value {
	if _, ok := other.(*Null); ok {
		return NewNumber(1)
	}
	return NewNumber(0)
}

func (n *Null) IsNotEqualTo(other interface{}) Value {
	if _, ok := other.(*Null); ok {
		return NewNumber(0)
	}
	return NewNumber(1)
}

func (n *Null) IsGreaterThan(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Can't compare null values", n.StartPos, n.EndPos)
}

func (n *Null) IsGreaterThanOrEqual(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare null values", n.StartPos, nil)
}

func (n *Null) IsLessThan(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare null values", n.StartPos, nil)
}

func (n *Null) IsLessThanOrEqual(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare null values", n.StartPos, nil)
}

func (n *Null) And(other interface{}) (Value, *Error) {
	return NewNumber(0), nil
}

func (n *Null) Or(other interface{}) (Value, *Error) {
	if o, ok := other.(Value); ok {
		if o.IsTrue() {
			return o, nil
		}
		return NewNumber(0), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", n.StartPos, nil)
}

func (n *Null) Not() Value {
	return NewNumber(1)
}

func (n *Null) IsTrue() bool {
	return false
}

func (n *Null) GetVal() interface{} {
	return nil
}

func (n *Null) Call(args []interface{}, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()
	return rr.Failure(NewRuntimeError("Can't call null values", n.StartPos, n.EndPos))
}

func (n *Null) AccessElement(index int, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()
	return rr.Failure(NewRuntimeError("Can't access element from a null value", n.StartPos, n.EndPos))
}
