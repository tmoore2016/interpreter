/*
Evaluator package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package evaluator

import (
	"fmt"

	"github.com/tmoore2016/interpreter/lib/ast"
	"github.com/tmoore2016/interpreter/lib/object"
)

// interpreter\evaluator\evaluator.go

// Vars are for types with fixed values, prevents creating new objects for identical references.
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evaluates each AST node by sending the ast.Node interface as input to the object package
func Eval(node ast.Node, env *object.Environment) object.Object {

	// Traverse each AST node and act according to type.
	switch node := node.(type) {

	// AST Program node is the top AST node, all AST statements are evaluated and returned as objects.
	case *ast.Program:
		return evalProgram(node, env)

	// AST block statement evaluates the primary or alternative (else) consequence of an If Expression
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// AST ExpressionStatement node is the top node for all expression statements and returns expressions
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions:

	// AST IntegerLiteral node returns an Integer Literal expression object with type and value
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	// AST StringLiteral node returns a String Literal expression object with type and value
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	// AST Boolean node returns a Boolean expression object with type and value
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	// AST prefix expression node evaluates the right side of the prefix expression, and then evaluates the prefix expression operator
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	// AST Infix expression evaluates the left and right node expressions, and then evaluates the operator
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	// AST if expression evaluates the If or If/Else expression node
	case *ast.IfExpression:
		return evalIfExpression(node, env)

	// AST Return statement evaluates the return statement value and creates a Return Value object
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// LetStatement evaluates an AST let statement identifier and value and sets the environment association.
	case *ast.LetStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		// Let statements can set an environment association
		env.Set(node.Name.Value, val)

	// Identifier evaluates an AST identifier and returns the environment value
	case *ast.Identifier:
		return evalIdentifier(node, env)

	// FunctionLiteral evaluates an AST function literal for params, body, and environment.
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	// CallExpression evaluates a list of expressions from a function as arguments, the process stops if there is an error.
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	}

	return nil
}

// evalProgram evaluates all AST program statement nodes as objects from evalProgramStatements
func evalProgram(program *ast.Program, env *object.Environment) object.Object {

	var result object.Object

	// Evaluate all statements in the AST
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {

		// If the last object evaluated was a ReturnValue, stop and return the unwrapped value
		case *object.ReturnValue:
			return result.Value

		// If the last object evaluated was an Error, stop and return the unwrapped value
		case *object.Error:
			return result
		}
	}

	// Return AST statements as objects
	return result
}

// evalBlockStatement evaluates AST block statements such as the primary and alternative consequences of an If expression
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	// For each block statement in range
	for _, statement := range block.Statements {
		result = Eval(statement, env)

		// If the block statement contains a Return Value Object or an Error object, stop and return
		if result != nil {

			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	// Return the result of the block statement
	return result
}

// newError creates error objects and returns their value (message)
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isError checks Eval() for errors
func isError(obj object.Object) bool {

	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

// nativeBoolToBooleanObject takes bools as input and returns either TRUE or FALSE vars
func nativeBoolToBooleanObject(input bool) *object.Boolean {

	if input {
		return TRUE
	}

	return FALSE
}

// evalPrefixExpression accepts the right side of the prefix expression as an object and applies NULL to the left side operator unless it is a ! or -, which are evaluated.
func evalPrefixExpression(operator string, right object.Object) object.Object {

	switch operator {

	case "!":
		return evalNotOperatorExpression(right)

	case "-":
		return evalMinusPrefixOperatorExpression(right)

	// Create new error object if unkown prefix expression is used
	default:
		return newError("Illegal prefix operator: %s%s", operator, right.Type())
	}
}

// evalNotOperatorExpression evaluates ! prefix expressions and returns the opposite.
func evalNotOperatorExpression(right object.Object) object.Object {

	switch right {

	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

// evalMinusPrefixOperatorExpression evaluates - prefix operators and if the right side of the prefix expression is an integer, returns the negative value.
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {

	// Return error if the right side expression isn't an integer
	if right.Type() != object.INTEGER_OBJ {
		return newError("Illegal prefix operation, expected integer, received: -%s", right.Type())
	}

	value := right.(*object.Integer).Value

	// Apply the negative value to an integer
	return &object.Integer{Value: -value}
}

// evalInfixExpression evaluates the left, right, and operator objects of an infix expression
func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {

	switch {

	// When left and right sides of the infix expression are integers, evaluate the integer infix expression
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	// If infix operator is ==, it will make a pointer comparison between left and right booleans. This works because there are only two Boolean expressions, the vars TRUE and FALSE and they are always in the same memory address. It won't work for integers, but those are compared in the switch statement above.
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	// Works the same as ==
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	// Create new error object if unrelated types are used
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	// When left or right side of the infix expression isn't an integer, return new error
	default:
		return newError("Illegal infix expression, expected integer-operator-integer, received: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression evaluates the operator of an infix expression.
func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {

	case "+":
		return &object.Integer{Value: leftVal + rightVal}

	case "-":
		return &object.Integer{Value: leftVal - rightVal}

	case "*":
		return &object.Integer{Value: leftVal * rightVal}

	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)

	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	// Return new error object if unsupported operator is used
	default:
		return newError("Invalid Infix Expression operator, expected ('+' , '-', '*', '/', '<', '>', '==', '!='),/n received: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIfExpression evaluates the conditions of an If or If/Else expression
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {

	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	// Condition is truthy, not null or false, return primary consequence
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)

		// If alternative consequence (else) applies, return that instead
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)

		// If neither primary or alternative consequence applies, return NULL
	} else {
		return NULL
	}
}

// isTruthy defines what truthy is: not NULL or FALSE
func isTruthy(obj object.Object) bool {

	switch obj {

	case NULL:
		return false

	case TRUE:
		return true

	case FALSE:
		return false

	default: // This isn't working as I'd like. If something isn't NULL or FALSE it should be true, but an identifier assigned a value isn't true or false in Doorkey because its never checked as a Boolean.
		return true
	}
}

// evalIdentifier evaluates an AST identifier node and retrieves its value from the environment association, if it exists.
func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {

	val, ok := env.Get(node.Value)

	if !ok {
		return newError("Identifier not found: " + node.Value)
	}

	return val
}

// evalExpressions evaluates ast.Expressions from a function in the context of the current environment
func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

// applyFunction verifies a function object and converts the function parameter to *object.Function to access the .Env and .Body fields.
func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)

	if !ok {
		return newError("Not a function, received type: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)

	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

// extendFunctionEnv creates a new *object.Environment that's enclosed by the function's environment. This allows the function's arguments to bind to the function's parameter names without overwriting the original environment.
func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {

	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

// unwrapReturnValue unwraps the outer environment for *object.ReturnValues so that evalBlockStatement will evaluate the entire block statement and not just the outer function.
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
