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
func Eval(node ast.Node) object.Object {

	// Traverse the AST nodes and act according to type
	switch node := node.(type) {

	// AST Program node is the top AST node, all AST statements below it are evaluated and returned as objects
	case *ast.Program:
		return evalProgram(node)

	// AST block statement evaluates the primary or alternative (else) consequence of an If Expression
	case *ast.BlockStatement:
		return evalBlockStatement(node)

	// AST ExpressionStatement node is the top node for all expression statements and returns expressions
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions:

	// AST IntegerLiteral node returns an Integer Literal expression object with type and value
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	// AST Boolean node returns a Boolean expression object with type and value
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	// AST prefix expression node evaluates the right side of the prefix expression, and then evaluates the prefix expression operator
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	// AST Infix expression evaluates the left and right node expressions, and then evaluates the operator
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	// AST if expression evaluates the If or If/Else expression node
	case *ast.IfExpression:
		return evalIfExpression(node)

	// AST Return statement evaluates the return statement value and creates a Return Value object
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}

	return nil
}

// evalProgram evaluates all AST program statement nodes as objects from evalProgramStatements
func evalProgram(program *ast.Program) object.Object {

	var result object.Object

	// Evaluate all statements in the AST
	for _, statement := range program.Statements {
		result = Eval(statement)

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
func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	// For each block statement in range
	for _, statement := range block.Statements {
		result = Eval(statement)

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
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: -%s", right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIfExpression evaluates the conditions of an If or If/Else expression
func evalIfExpression(ie *ast.IfExpression) object.Object {

	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	// Condition is truthy, not null or false, return primary consequence
	if isTruthy(condition) {
		return Eval(ie.Consequence)

		// If alternative consequence (else) applies, return that instead
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)

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

	default:
		return true
	}
}
