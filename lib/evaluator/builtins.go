/*
Builtins environment for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package evaluator

import (
	"fmt"

	"github.com/tmoore2016/interpreter/lib/object"
)

// Separate Builtins environment, allowing builtin Go functions to be called through Doorkey.
var builtins = map[string]*object.Builtin{

	// puts function allows Doorkey to print to terminal
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},

	// length (len) function for counting characters in a string
	"len": &object.Builtin{
		// Fail if number of evals isn't 1
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			// If object type is array, length will return the number of elements as an integer
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			// If object evaluated is type string, length will return the number of characters
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			// In all other cases return an error
			default:
				return newError("argument to 'len' not supported, got %s", args[0].Type())
			}
		},
	},

	// first() retrieves the first element in an array
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {

			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'first' must be an ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},

	// last() retrieves the last element in an array
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'last' must be an ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)

			length := len(arr.Elements)

			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},

	// tail() returns a new array containing all of the elements in the input array, except the first.
	"tail": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'tail' must be an ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},

	// push() returns a new array containing all of the elements of the input array, plus the new element
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to 'push' must be an ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
}
