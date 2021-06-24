package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type BuiltinFunction struct {
	Name string
	ArgNames []string
	OnCall func([]interface{}) *RuntimeResult
	StartPos, EndPos *Position
}

func NewBuiltinFunction(n string, a []string, oc func([]interface{}) *RuntimeResult) Value {
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
		if o.IsTrue() {
			return o, nil
		}
		return NewNumber(0), nil
	}

	return nil, NewRuntimeError("Can't compare values of different types", f.StartPos, nil)
}

func (f *BuiltinFunction) Or(other interface{}) (Value, *Error) {
	if o, ok := other.(Value); ok {
		if o.IsTrue() {
			return o, nil
		}
		return NewNumber(0), nil
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

func (f *BuiltinFunction) Call(args []interface{}, ctx *Context) *RuntimeResult {
	rr := NewRuntimeResult()
	val := rr.Register(f.OnCall(args))
	if rr.ShouldReturn() {
		return rr
	}
	return rr.Success(val)
}

var BuiltinPrint = NewBuiltinFunction(
	"print",
	[]string{"...values"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()
		fmt.Print(args...)
		return rr.Success(NewNull())
	},
)

var BuiltinPrintln = NewBuiltinFunction(
	"println",
	[]string{"...values"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()
		fmt.Println(args...)
		return rr.Success(NewNull())
	},
)

var BuiltinScan = NewBuiltinFunction(
	"scan",
	[]string{"prompt"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		prompt := "> "
		if len(args) > 0 {
			prompt = args[0].(*String).GetVal().(string)
		}

		text, err := GetInput(prompt)

		if err != nil {
			return rr.Failure(NewRuntimeError("Failed to get scan from the stdin", nil, nil))
		}

		return rr.Success(NewString(text))
	},
)

var BuiltinLen = NewBuiltinFunction(
	"len",
	[]string{"list|string"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			arg := args[0]
			switch val := arg.(type) {
				case *List:
					return rr.Success(val.Length)
				case *String:
					return rr.Success(NewNumber(float64(len(val.Value))))
			}

			return rr.Failure(NewRuntimeError("len() only works for strings or lists", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected one argument to be passed to len()", nil, nil))
	},
)

var BuiltinTrim = NewBuiltinFunction(
	"trim",
	[]string{"string"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			arg := args[0]
			if str, ok := arg.(*String); ok {
				return rr.Success(NewString(strings.TrimSpace(str.Value)))
			}

			return rr.Failure(NewRuntimeError("trim() only works for strings", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected one argument to be passed to trim()", nil, nil))
	},
)

var BuiltinUpper = NewBuiltinFunction(
	"upper",
	[]string{"string"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			arg := args[0]
			if str, ok := arg.(*String); ok {
				return rr.Success(NewString(strings.ToUpper(str.Value)))
			}

			rr.Failure(NewRuntimeError("upper() only works for strings", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected one argument to be passed to upper()", nil, nil))
	},
)

var BuiltinLower = NewBuiltinFunction(
	"lower",
	[]string{"string"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			arg := args[0]
			if str, ok := arg.(*String); ok {
				return rr.Success(NewString(strings.ToLower(str.Value)))
			}

			return rr.Failure(NewRuntimeError("lower() only works for strings", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected one argument to be passed to lower()", nil, nil))
	},
)

var BuiltinReplace = NewBuiltinFunction(
	"replace",
	[]string{"string"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) == 3 {
			if str, ok := args[0].(*String); ok {
				if old, ok := args[1].(*String); ok {
					if new, ok := args[2].(*String); ok {
						return rr.Success(NewString(strings.ReplaceAll(str.Value, old.Value, new.Value)))
					}
				}
			}

			return rr.Failure(NewRuntimeError("replace() only works for strings", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected 3 arguments to be passed to replace()", nil, nil))
	},
)

var BuiltinAppend = NewBuiltinFunction(
	"append",
	[]string{"list", "...elements"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 1 {
			arg := args[0]
			newEl := args[1:]
			if list, ok := arg.(*List); ok {
				el := append(list.Elements, newEl...)
				return rr.Success(NewList(el))
			}

			return rr.Failure(NewRuntimeError("append() only works for lists", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected at least 2 argument to be passed to append()", nil, nil))
	},
)

var BuiltinPrepend = NewBuiltinFunction(
	"prepend",
	[]string{"list", "...elements"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 1 {
			arg := args[0]
			newEl := args[1:]
			if list, ok := arg.(*List); ok {
				el := append(newEl, list.Elements...)
				return rr.Success(NewList(el))
			}

			rr.Failure(NewRuntimeError("prepend() only works for lists", nil, nil))
		}

		rr.Failure(NewRuntimeError("Expected at least 2 argument to be passed to prepend()", nil, nil))

		return nil
	},
)

var BuiltinShift = NewBuiltinFunction(
	"shift",
	[]string{"list"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			arg := args[0]
			if list, ok := arg.(*List); ok {
				return rr.Success(NewList(list.Elements[1:]))
			}

			rr.Failure(NewRuntimeError("shift() only works for lists", nil, nil))
		}

		rr.Failure(NewRuntimeError("Expected one argument to be passed to shift()", nil, nil))

		return nil
	},
)

var BuiltinPop = NewBuiltinFunction(
	"pop",
	[]string{"list"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			arg := args[0]
			if list, ok := arg.(*List); ok {
				return rr.Success(NewList(list.Elements[:len(list.Elements) - 1]))
			}

			rr.Failure(NewRuntimeError("pop() only works for lists", nil, nil))
		}

		rr.Failure(NewRuntimeError("Expected one argument to be passed to pop()", nil, nil))

		return nil
	},
)

var BuiltinExit = NewBuiltinFunction(
	"exit",
	[]string{"code"},
	func(args []interface{}) *RuntimeResult {
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

var BuiltinNum = NewBuiltinFunction(
	"num",
	[]string{"value"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			if val, ok := args[0].(*String); ok {
				num, err := strconv.ParseFloat(val.GetVal().(string), 64)
				if err != nil {
					return rr.Failure(NewRuntimeError("num() only converts string numbers into raw numbers", nil, nil))
				}
				return rr.Success(NewNumber(num))
			}
		}

		return rr.Failure(NewRuntimeError("Expected one argument to be passed to num()", nil, nil))
	},
)

var BuiltinStr = NewBuiltinFunction(
	"str",
	[]string{"value"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) > 0 {
			if val, ok := args[0].(Value); ok {
				return rr.Success(NewString(val.String()))
			}
			return rr.Failure(NewRuntimeError("Expected a value to be passed to str()", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected one argument to be passed to str()", nil, nil))
	},
)

var BuiltinMap = NewBuiltinFunction(
	"map",
	[]string{"list, fun"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) == 2 {
			if list, ok := args[0].(*List); ok {
				if fun, ok := args[1].(Value); ok {
					newList := []interface{}{}

					for _, val := range list.Elements {
						ctx := NewContext("map")
						res := rr.Register(fun.Call([]interface{}{val}, ctx))
						if rr.ShouldReturn() {
							return rr
						}
						newList = append(newList, res)
					}

					return rr.Success(NewList(newList))
				}
			}
			return rr.Failure(NewRuntimeError("Expected first argument of map() to be a list and second to be a value", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected 2 arguments to be passed to map()", nil, nil))
	},
)

var BuiltinReduce = NewBuiltinFunction(
	"reduce",
	[]string{"list, fun, initialValue"},
	func(args []interface{}) *RuntimeResult {
		rr := NewRuntimeResult()

		if len(args) == 3 {
			if list, ok := args[0].(*List); ok {
				if fun, ok := args[1].(Value); ok {
					if initial, ok := args[2].(Value); ok {
						accum := initial

						for _, curr := range list.Elements {
							ctx := NewContext("reduce")
							accum = rr.Register(fun.Call([]interface{}{accum, curr}, ctx))
							if rr.ShouldReturn() {
								return rr
							}
						}

						return rr.Success(accum)
					}
				}
			}
			return rr.Failure(NewRuntimeError("Expected first argument of reduce() to be a list and second to be a value", nil, nil))
		}

		return rr.Failure(NewRuntimeError("Expected 2 arguments to be passed to reduce()", nil, nil))
	},
)
