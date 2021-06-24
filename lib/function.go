package main

type Function struct {
	Name string
	ArgNames []string
	Body interface{}
	StartPos, EndPos *Position
}

func NewFunction(n string, a []string, b interface{}) Value {
	if n == "" {
		n = "anonymous"
	}

	f := &Function{
		Name: n,
		ArgNames: a,
		Body: b,
	}

	return f
}

func (f *Function) String() string {
	str := f.Name + "("
	for i, arg := range f.ArgNames {
		if i != 0 {
			str += ", "
		}
		str += arg
	}
	str += ")"
	return str
}

func (f *Function) SetPos(sp, ep *Position) Value {
	f.StartPos = sp
	f.EndPos = ep
	if ep == nil {
		endPos := *sp
		endPos.Advance("")
		f.EndPos = &endPos
	}
	return f
}

func (f *Function) AddTo(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '+' operation on a function", f.StartPos, f.EndPos)
}

func (f *Function) SubBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '-' operation on a function", f.StartPos, f.EndPos)
}

func (f *Function) MulBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '*' operation on a function", f.StartPos, f.EndPos)
}

func (f *Function) DivBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '/' operation on a function", f.StartPos, f.EndPos)
}

func (f *Function) Mod(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '%' operation on a function", f.StartPos, f.EndPos)
}

func (f *Function) Pow(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '^' operation on a function", f.StartPos, f.EndPos)
}

func (f *Function) IsEqualTo(other interface{}) Value {
	return NewNumber(0)
}

func (f *Function) IsNotEqualTo(other interface{}) Value {
	return NewNumber(1)
}

func (f *Function) IsGreaterThan(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Can't compare functions", f.StartPos, f.EndPos)
}

func (f *Function) IsGreaterThanOrEqual(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare functions", f.StartPos, nil)
}

func (f *Function) IsLessThan(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare functions", f.StartPos, nil)
}

func (f *Function) IsLessThanOrEqual(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare functions", f.StartPos, nil)
}

func (f *Function) And(other interface{}) (Value, *Error) {
	if o, ok := other.(Value); ok {
		if o.IsTrue() {
			return o, nil
		}
		return NewNumber(0), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", f.StartPos, nil)
}

func (f *Function) Or(other interface{}) (Value, *Error) {
	if o, ok := other.(Value); ok {
		if o.IsTrue() {
			return o, nil
		}
		return NewNumber(0), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", f.StartPos, nil)
}

func (f *Function) Not() Value {
	return NewNumber(0)
}

func (f *Function) IsTrue() bool {
	return true
}

func (f *Function) GetVal() interface{} {
	return nil
}

func (f *Function) Call(args []interface{}, ctx *Context) (Value, *Error) {
	i := NewInterpretor()
	newCtx := NewContext(f.Name)
	newCtx.Parent = ctx
	newCtx.SymbolTable = NewSymbolTable()
	newCtx.SymbolTable.Parent = newCtx.Parent.SymbolTable

	for key, argVal := range args {
		argName := f.ArgNames[key]
		newCtx.SymbolTable.Set(argName, argVal.(Value))
	}

	val := i.Visit(f.Body, newCtx)

	return val, nil
}
