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

// TestStringObject fails if the expected object type and value aren't the actual type or value.
func TestStringObject(t *testing.T) {
	input := `"Doorkey has strings!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Doorkey has strings!" {
		t.Errorf("string has wrong value. got=%q", str.Value)
	}
}

// TestStringConcatenation tests if strings can be added via + infix operator
func TestStringConcatenation(t *testing.T) {
	input := `"Peanut" + " " + "Butter"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String) // Ok if object type is string

	// Fail if object type is not string
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	// Fail if string value isn't input value
	if str.Value != "Peanut Butter" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
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
			"Illegal prefix operation, expected integer, received: -BOOLEAN",
		},
		{
			"false + true",
			"Illegal infix expression, expected integer-operator-integer, received: BOOLEAN + BOOLEAN",
		},
		{
			"8; true + false; 8",
			"Illegal infix expression, expected integer-operator-integer, received: BOOLEAN + BOOLEAN",
		},
		{
			"if (8 > 6) { true + false; }",
			"Illegal infix expression, expected integer-operator-integer, received: BOOLEAN + BOOLEAN",
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
			"Illegal infix expression, expected integer-operator-integer, received: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Identifier not found: foobar",
		},
		{
			`"Hulk" - "Smash"`,
			"Invalid operator: STRING - STRING",
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

func TestFunctionObject(t *testing.T) {

	// Test input function, () is parameters, {} is function statement
	input := "fn(x) { x + 2; };"

	expectedBody := "(x + 2)"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)

	// Object isn't a function, return an error with actual type and value.
	if !ok {
		t.Fatalf("Object is not a function. Got=%T (%+v)", evaluated, evaluated)
	}

	// Only 1 parameter in input, return error with actual parameters if more than 1.
	if len(fn.Parameters) != 1 {
		t.Fatalf("Function has the wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("Parameter is not 'x'. Got%q", fn.Parameters[0])
	}

	if fn.Body.String() != expectedBody {

		t.Fatalf("Body is not %q. Got = %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let myNumber = fn(x) { x; }; myNumber(15);", 15},
		{"let yourNumber = fn(x) { return x; }; yourNumber(16);", 16},
		{"let double = fn(x) { x * 2; }; double(18);", 36},
		{"let add = fn(x, y) { x + y; }; add(16, 14);", 30},
		{"let addTwice = fn(x, y) { x + y; }; addTwice(5 + 5, addTwice(5, 5));", 20},
		{"fn(x) { x; }(8)", 8},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}
