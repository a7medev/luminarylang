package main

import (
	"fmt"
	"os"
)

type BuiltinFunction struct {
	Name string
	ArgNames []string
	OnCall func([]interface{}) Value
	StartPos, EndPos *Position
}

func NewBuiltinFunction(n string, a []string, oc func([]interface{}) Value) Value {
	f := &BuiltinFunction{
		Name: n,
		ArgNames: a,
		OnCall: oc,
	}

	return f
}

func (f *BuiltinFunction) String() string {
	str := "builtin:" + f.Name + "("
	for i, arg := range f.ArgNames {
		if i != 0 {
			str += ", "
		}
		str += arg
	}
	str += ")"
	return str
}

func (f *BuiltinFunction) SetPos(sp, ep *Position) Value {
	f.StartPos = sp
	f.EndPos = ep
	if ep == nil {
		endPos := *sp
		endPos.Advance("")
		f.EndPos = &endPos
	}
	return f
}

func (f *BuiltinFunction) AddTo(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '+' operation on a function", f.StartPos, f.EndPos)
}

func (f *BuiltinFunction) SubBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '-' operation on a function", f.StartPos, f.EndPos)
}

func (f *BuiltinFunction) MulBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '*' operation on a function", f.StartPos, f.EndPos)
}

func (f *BuiltinFunction) DivBy(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '/' operation on a function", f.StartPos, f.EndPos)
}

func (f *BuiltinFunction) Mod(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '%' operation on a function", f.StartPos, f.EndPos)
}

func (f *BuiltinFunction) Pow(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Invalid '^' operation on a function", f.StartPos, f.EndPos)
}

func (f *BuiltinFunction) IsEqualTo(other interface{}) Value {
	return NewNumber(0)
}

func (f *BuiltinFunction) IsNotEqualTo(other interface{}) Value {
	return NewNumber(1)
}

func (f *BuiltinFunction) IsGreaterThan(other interface{}) (Value, *Error) {
	return nil, NewInvalidSyntaxError("Can't compare Builtinfunctions", f.StartPos, f.EndPos)
}

func (f *BuiltinFunction) IsGreaterThanOrEqual(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare Builtinfunctions", f.StartPos, nil)
}

func (f *BuiltinFunction) IsLessThan(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare Builtinfunctions", f.StartPos, nil)
}

func (f *BuiltinFunction) IsLessThanOrEqual(other interface{}) (Value, *Error) {
	return nil, NewRuntimeError("Can't compare Builtinfunctions", f.StartPos, nil)
}

func (f *BuiltinFunction) And(other interface{}) (Value, *Error) {
	if o, ok := other.(Value); ok {
		var val float64 = 0
		if o.IsTrue() {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", f.StartPos, nil)
}

func (f *BuiltinFunction) Or(other interface{}) (Value, *Error) {
	if o, ok := other.(Value); ok {
		var val float64 = 0
		if o.IsTrue() {
			val = 1
		}
		return NewNumber(val), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", f.StartPos, nil)
}

func (f *BuiltinFunction) Not() Value {
	return NewNumber(0)
}

func (f *BuiltinFunction) IsTrue() bool {
	return true
}

func (f *BuiltinFunction) GetVal() interface{} {
	return nil
}

func (f *BuiltinFunction) Call(args []interface{}, ctx *Context) (Value, *Error) {
	val := f.OnCall(args)
	return val, nil
}

var BuiltinPrint = NewBuiltinFunction(
	"print",
	[]string{"...values"},
	func(args []interface{}) Value {
		fmt.Print(args...)
		return nil
	},
)

var BuiltinPrintln = NewBuiltinFunction(
	"println",
	[]string{"...values"},
	func(args []interface{}) Value {
		fmt.Println(args...)
		return nil
	},
)

var BuiltinScan = NewBuiltinFunction(
	"scan",
	[]string{"prompt"},
	func(args []interface{}) Value {
		prompt := "> "
		if len(args) > 0 {
			prompt = args[0].(*String).GetVal().(string)
		}

		text, _ := GetInput(prompt)

		return NewString(text)
	},
)

var BuiltinLen = NewBuiltinFunction(
	"len",
	[]string{"list|string"},
	func(args []interface{}) Value {
		if len(args) > 0 {
			arg := args[0]
			if val, ok := arg.(Value); ok {
				switch raw := val.GetVal().(type) {
				case string:
					return NewNumber(float64(len(raw)))
				case []interface{}:
					return NewNumber(float64(len(raw)))
				}
			}

			err := NewRuntimeError("len() only works for strings or lists", nil, nil)
			fmt.Println(err)
			os.Exit(0)
		}

		err := NewRuntimeError("Expected one argument to be passed to len()", nil, nil)
		fmt.Println(err)
		os.Exit(0)

		return nil
	},
)

var BuiltinAppend = NewBuiltinFunction(
	"append",
	[]string{"list", "...elements"},
	func(args []interface{}) Value {
		if len(args) > 1 {
			arg := args[0]
			newEl := args[1:]
			if list, ok := arg.(*List); ok {
				el := append(list.Elements, newEl...)
				return NewList(el)
			}

			err := NewRuntimeError("append() only works for lists", nil, nil)
			fmt.Println(err)
			os.Exit(0)
		}

		err := NewRuntimeError("Expected at least 2 argument to be passed to append()", nil, nil)
		fmt.Println(err)
		os.Exit(0)

		return nil
	},
)

var BuiltinPrepend = NewBuiltinFunction(
	"prepend",
	[]string{"list", "...elements"},
	func(args []interface{}) Value {
		if len(args) > 1 {
			arg := args[0]
			newEl := args[1:]
			if list, ok := arg.(*List); ok {
				el := append(newEl, list.Elements...)
				return NewList(el)
			}

			err := NewRuntimeError("prepend() only works for lists", nil, nil)
			fmt.Println(err)
			os.Exit(0)
		}

		err := NewRuntimeError("Expected at least 2 argument to be passed to prepend()", nil, nil)
		fmt.Println(err)
		os.Exit(0)

		return nil
	},
)

var BuiltinShift = NewBuiltinFunction(
	"shift",
	[]string{"list"},
	func(args []interface{}) Value {
		if len(args) > 0 {
			arg := args[0]
			if list, ok := arg.(*List); ok {
				return NewList(list.Elements[1:])
			}

			err := NewRuntimeError("shift() only works for lists", nil, nil)
			fmt.Println(err)
			os.Exit(0)
		}

		err := NewRuntimeError("Expected one argument to be passed to shift()", nil, nil)
		fmt.Println(err)
		os.Exit(0)

		return nil
	},
)

var BuiltinPop = NewBuiltinFunction(
	"pop",
	[]string{"list"},
	func(args []interface{}) Value {
		if len(args) > 0 {
			arg := args[0]
			if list, ok := arg.(*List); ok {
				return NewList(list.Elements[:len(list.Elements) - 1])
			}

			err := NewRuntimeError("pop() only works for lists", nil, nil)
			fmt.Println(err)
			os.Exit(0)
		}

		err := NewRuntimeError("Expected one argument to be passed to pop()", nil, nil)
		fmt.Println(err)
		os.Exit(0)

		return nil
	},
)

var BuiltinExit = NewBuiltinFunction(
	"exit",
	[]string{"code"},
	func(args []interface{}) Value {
		var code interface{} = 0
		if len(args) > 0 {
			code = args[0]
		}
		switch c := code.(type) {
		case *Number:
			os.Exit(int(c.GetVal().(float64)))
		case int:
			os.Exit(c)
		default:
			os.Exit(0)
		}
		return nil
	},
)
