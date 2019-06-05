/*
Builtins environment for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package evaluator

import (
	"github.com/tmoore2016/interpreter/lib/object"
)

// Separate Builtins environment, allowing builtin Go functions to be called through Doorkey.
var builtins = map[string]*object.Builtin{

	// length (len) function for counting characters in a string
	"len": &object.Builtin{
		// Fail if number of evals isn't 1
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			// If object evaluated is type string, return its index length as an integer
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			// In all other cases return an error
			default:
				return newError("argument to 'len' not supported, got %s", args[0].Type())
			}
		},
	},
}
