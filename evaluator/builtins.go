package evaluator

import (
	"fmt"

	"github.com/mikeraimondi/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported. got %s",
					args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return NULL
			default:
				return newError("argument to `first` not supported. got %s",
					args[0].Type())
			}
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				l := len(arr.Elements)
				if l > 0 {
					return arr.Elements[l-1]
				}
				return NULL
			default:
				return newError("argument to `last` not supported. got %s",
					args[0].Type())
			}
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d. want 1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				l := len(arr.Elements)
				if l > 0 {
					result := make([]object.Object, l-1, l-1)
					copy(result, arr.Elements[1:l])
					return &object.Array{Elements: result}
				}
				return NULL
			default:
				return newError("argument to `rest` not supported. got %s",
					args[0].Type())
			}
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got %d. want 2", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				l := len(arr.Elements)
				result := make([]object.Object, l+1, l+1)
				copy(result, arr.Elements)
				result[l] = args[1]
				return &object.Array{Elements: result}
			default:
				return newError("argument to `push` not supported. got %s",
					args[0].Type())
			}
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
