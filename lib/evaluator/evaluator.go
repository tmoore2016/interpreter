/*
Evaluator package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package evaluator

import (
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
		return evalStatements(node.Statements)

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
		return evalPrefixExpression(node.Operator, right)

	// AST Infix expression evaluates the left and right node expressions, and then evaluates the operator
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	// AST block statement evaluates the primary or alternative (else) consequence of an If Expression
	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	// AST if expression evaluates the If or If/Else expression node
	case *ast.IfExpression:
		return evalIfExpression(node)

	}

	return nil
}

// evalStatements evaluates all AST statement nodes as objects
func evalStatements(stmts []ast.Statement) object.Object {

	var result object.Object

	// Evaluate all statements in the AST
	for _, statement := range stmts {
		result = Eval(statement)
	}

	// Return AST statements as objects
	return result
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

	default:
		return NULL
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

	// Return null if the right side expression isn't an integer
	if right.Type() != object.INTEGER_OBJ {
		return NULL
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

	// When left or right side of the infix expression isn't an integer, return null
	default:
		return NULL
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

	default:
		return NULL
	}
}

// evalIfExpression evaluates the conditions of an If or If/Else expression
func evalIfExpression(ie *ast.IfExpression) object.Object {

	condition := Eval(ie.Condition)

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
