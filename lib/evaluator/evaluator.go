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

	// AST IntegerLiteral node returns an Integer Literal expression object with type and value
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	// AST StringLiteral node returns a String Literal expression object with type and value
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	// AST ArrayLiteral node returns an array literal expression object with element and index number
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	// AST IndexExpression node returns an array's index expression object from the running environment
	case *ast.IndexExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	// AST Boolean node returns a Boolean expression object with type and value
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	// AST HashLiteral node evaluates HashLiterals
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

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

	// When left and right sides are strings, evaluate a string infix expression
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	// If infix operator is ==, it will make a pointer comparison between left and right booleans. This works because there are only two Boolean expressions, the vars TRUE and FALSE and they are always in the same memory address. It won't work for integers, but those are compared in the switch statement above.
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	// Works the same as ==
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	// Create new error object if unrelated types are compared
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

// evalStringInfixExpression evaluates string operations. Currently only concatenation.
// To add == and != String comparisons, put here and use values rather than pointers.
func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator != "+" {
		return newError("Invalid operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
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

// evalHashLiteral evaluates the key node to determine it is a hashable type, then evaluates the value node and adds the key-value pair to the pairs map by calling HashKey(). A new HashPair object is created by pointing to key and value and added to pairs.
func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)

		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)

		if !ok {
			return newError("Unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)

		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()

		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

// evalIdentifier evaluates an AST identifier node and retrieves its value from the environment association, if it exists.
func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	// Fallback when identifier is not bound to value in current environment, checks builtin functions (builtins.go)
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	// Failure mode
	return newError("Identifier not found: " + node.Value)
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

// evalIndexExpression accepts an array object and the array's index, if both are valid it calls evalArrayIndexExpression
func evalIndexExpression(left, index object.Object) object.Object {
	switch {

	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)

	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)

	default:
		return newError("Index operator not supported: %s", left.Type())
	}
}

// evalArrayIndexExpression matches an element of an array with its index, and returns an object containing the element and index number
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

// evalHashIndexExpression matches a hash key to its value, if the hash key doesn't exist, returns null
func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)

	if !ok {
		return newError("Unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
}

// applyFunction verifies a function object and converts the function parameter to *object.Function to access the .Env and .Body fields.
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	// Standard object.Function types
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	// Builtin function types
	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("Not a function, received type: %s", fn.Type())
	}
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
