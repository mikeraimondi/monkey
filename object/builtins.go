package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported. got %s",
					args[0].Type())
			}
		},
		},
	},
	{
		"puts",
		&Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	},
	{
		"first",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch args[0].(type) {
			case *Array:
				arr := args[0].(*Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return nil
			default:
				return newError("argument to `first` not supported. got %s",
					args[0].Type())
			}
		},
		},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch args[0].(type) {
			case *Array:
				arr := args[0].(*Array)
				l := len(arr.Elements)
				if l > 0 {
					return arr.Elements[l-1]
				}
				return nil
			default:
				return newError("argument to `last` not supported. got %s",
					args[0].Type())
			}
		},
		},
	},
	{
		"rest",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch args[0].(type) {
			case *Array:
				arr := args[0].(*Array)
				l := len(arr.Elements)
				if l > 0 {
					result := make([]Object, l-1, l-1)
					copy(result, arr.Elements[1:l])
					return &Array{Elements: result}
				}
				return nil
			default:
				return newError("argument to `rest` not supported. got %s",
					args[0].Type())
			}
		},
		},
	},
	{
		"push",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got %d. want 2", len(args))
			}

			switch args[0].(type) {
			case *Array:
				arr := args[0].(*Array)
				l := len(arr.Elements)
				result := make([]Object, l+1, l+1)
				copy(result, arr.Elements)
				result[l] = args[1]
				return &Array{Elements: result}
			default:
				return newError("argument to `push` not supported. got %s",
					args[0].Type())
			}
		},
		},
	},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}

	return nil
}
