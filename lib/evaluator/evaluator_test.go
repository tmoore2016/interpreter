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
		{
			`{"Hulk": "Smash"}[fn(x) {x}];`,
			"Unusable as hash key: FUNCTION",
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

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("five")`, 4},
		{`len("Hulk Smash!")`, 11},
		{`len(8)`, "argument to 'len' not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`let arr = [4, 5 * 5, 32]; len(arr)`, 3},
		{`let arr = ["thursday", "friday", "saturday"]; len(arr[1])`, 6},
		{`let arr = [2, 4, 6]; first(arr)`, 2},
		{`let arr = []; first(arr)`, nil},
		{`let arr = [10, 100, 1000, 10000]; last(arr)`, 10000},
		{`let arr = []; last(arr)`, nil},
		{`let arr = [20, 40, 60, 80, 100]; tail(arr)`, []int{40, 60, 80, 100}},
		{`let arr = []; tail(arr)`, nil},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, "argument to 'push' must be an ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {

		// If Go returns an int, len function worked, return expected and value
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		// If Go returns a string, return error object
		case string:
			errObj, ok := evaluated.(*object.Error)

			// If string isn't an error object, return evals
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			// If error message received isn't expected, show expected and actual error messages
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}

		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not an array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("Wrong number of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

// TestArrayLiterals is a test for array elements and indexing
func TestArrayLiterals(t *testing.T) {
	input := "[1, 8 * 8, 4 + 4]"

	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("Object is not an Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("Array has wrong number of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 64)
	testIntegerObject(t, result.Elements[2], 8)
}

// TestArrayIndexExpressions tests calling array elements by index number
func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[7, 8, 9][1]",
			8,
		},
		{
			"[6, 14, 0][2]",
			0,
		},
		{
			"let i = 0; [100][i];",
			100,
		},
		{
			"let myArray = [10, 100 * 10, 20]; myArray[1];",
			1000,
		},
		{
			"let myArray = [2 + 2, 12, 16, 34, 88, 23]; myArray[4];",
			88,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// TestHashLiterals tests that when an ast.HashLiteral is encountered, a new object.Hash with HashPairs is mapped to the matching HashKey using the Pairs attribute.
func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10-9,
		two : 1 + 1,
		"thr" + "ee": 6/2,
		4: 4,
		true: 5,
		false: 6	
	}`

	evaluated := testEval(input)

	result, ok := evaluated.(*object.Hash)

	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]

		if !ok {
			t.Errorf("No pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

// TestHashIndexExpressions tests calling hash index expressions
func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"age": 5}["age"]`,
			5,
		},
		{
			`{"age": 5}["height"]`,
			nil,
		},
		{
			`let key = "age"; {"age": 5}[key]`,
			5,
		},
		{
			`{}["age"]`,
			nil,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
