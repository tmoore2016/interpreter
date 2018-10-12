/*
Evaluator_test package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package evaluator

import (
	"testing"

	"github.com/tmoore2016/interpreter/lib/lexer"
	"github.com/tmoore2016/interpreter/lib/object"
	"github.com/tmoore2016/interpreter/lib/parser"
)

// TestEvalIntegerExpressions checks the type and value of integer input
func TestEvalIntegerExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		// Test input
		{"8", 8},
		{"32", 32},
	}

	// For each test input, send to testEval() and confirm that the evaluated output is equal to expected output
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// testEval sends input to the lexer, parses it, assigns it to an AST program node, and returns the evaluated node.
func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

// testIntegerObject fails if the expected type or value of the evaluated object isn't the actual type or value
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {

	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("Object is not an Integer. got=%T (%+v", obj, obj)

		return false
	}

	if result.Value != expected {
		t.Errorf("Object has the wrong value. got=%d, want=%d", result.Value, expected)

		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {

	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("Object is not a Boolean. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}
