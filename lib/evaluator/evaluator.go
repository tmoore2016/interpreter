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
