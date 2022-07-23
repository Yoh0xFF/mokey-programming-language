package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.BuiltinObject{
	"len": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newErrorObject("Wrong number of arguments. got=%d, want=1",
				len(args))
		}

		switch arg := args[0].(type) {
		case *object.ArrayObject:
			return &object.IntObject{Value: int64(len(arg.Elements))}
		case *object.StringObject:
			return &object.IntObject{Value: int64(len(arg.Value))}
		default:
			return newErrorObject("Argument to `len` not supported, got %s",
				args[0].Type())
		}
	},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObject("Wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newErrorObject("Argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.ArrayObject)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObject("Wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newErrorObject("Argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.ArrayObject)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObject("Wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newErrorObject("Argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.ArrayObject)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.ArrayObject{Elements: newElements}
			}

			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorObject("Wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newErrorObject("Argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.ArrayObject)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.ArrayObject{Elements: newElements}
		},
	},
}
