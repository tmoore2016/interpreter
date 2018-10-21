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

// testEval sends input to the lexer, parses it, assigns it to an AST program node, and returns the evaluated node.
func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

// TestEvalIntegerExpressions checks the type and value of integer input
func TestEvalIntegerExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		// Test input
		{"8", 8},
		{"32", 32},
		{"-8", -8},
		{"-32", -32},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2 * 2", 64},
		{"-64 + 128 + -64", 0},
		{"8 * 8 + 6 - 75", -5},
		{"100 / 10 * 4 - 40 + 5", 5},
		{"(8 - 6) / 2 - 1", 0},
		{"-(2 + 2) - 10", -14},
		{"(6 + 5 - 2 + 1) * 4 / 8 + -9", -4},
	}

	// For each test input, send to testEval() and confirm that the evaluated output is equal to expected output
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
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

// TestEvalBooleanExpression tests the evaluation of Boolean expressions
func TestEvalBooleanExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"5 < 10", true},
		{"5 > 10", false},
		{"10 < 5", false},
		{"10 > 5", true},
		{"1 < 1", false},
		{"1 > 1", false},
		{"4 == 4", true},
		{"4 != 4", false},
		{"4 == 5", false},
		{"4 != 5", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(3 < 6) == true", true},
		{"(3 > 6) == false", true},
		{"(3 > 6) == true", false},
		{"(3 < 6) == false", false},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)

		testBooleanObject(t, evaluated, tt.expected)
	}
}

// testBooleanObject tests Boolean objects for type and value
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {

	result, ok := obj.(*object.Boolean)

	if !ok {

		t.Errorf("Object is not a Boolean. got=%T (%+v)", obj, obj)

		return false
	}

	if result.Value != expected {

		t.Errorf("Object has wrong value. Got=%t, want=%t", result.Value, expected)

		return false
	}

	return true
}

// TestNotOperator tests the evaluation of the ! prefix expression operator
func TestNotOperator(t *testing.T) {

	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)

		testBooleanObject(t, evaluated, tt.expected)
	}
}

// TestIfElseExpressions tests the evaluation of If/Else conditionals
func TestIfElseExpressions(t *testing.T) {

	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 == 2) { 5 } else { 10 }", 10},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)

		// If actual integer is expected integer, test the integer object
		if ok {
			testIntegerObject(t, evaluated, int64(integer))

			// If a conditional doesn't evaluate to a value, it should return NULL.
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// testNullObject confirms that an object is NULL, returns true.
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {

		t.Errorf("Expected NULL, Object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 8 * 6; 9;", 48},
		{"12; return 3 + 3; 7;", 6},
		// A block statement with 2 returns, only the first should return (10)
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`,
			10,
		},
	}

	for _, tt := range tests {

		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// TestErrorHandling tests the evaluation of error objects and error message handling
func TestErrorHandling(t *testing.T) {

	tests := []struct {
		input           string
		expectedMessage string
	}{
		// input, expected error message
		{
			"8 + false;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"8 + true; 8;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"false + true",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"8; true + false; 8",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (8 > 6) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}

				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		// Error! No error detected!
		if !ok {
			t.Errorf("No error object returned. got=%T(%+v", evaluated, evaluated)
			continue
		}

		// Actual error message doesn't match expected error message.
		if errObj.Message != tt.expectedMessage {
			t.Errorf("Wrong error message. Expected=%q, got =%q", tt.expectedMessage, errObj.Message)
		}
	}
}

// TestLetStatements tests let statement evaluation
func TestLetStatements(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 8; a;", 8},
		{"let a = 8 * 8; a;", 64},
		{"let a = 8; let b = a; b;", 8},
		{"let a = 8; let b = a; let c = a + b + 8; c;", 24},
	}

	for _, tt := range tests {

		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
