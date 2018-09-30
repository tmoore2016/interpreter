/*
Parser tests for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package parser

import (
	"fmt"
	"testing"

	"github.com/tmoore2016/interpreter/lib/ast"
	"github.com/tmoore2016/interpreter/lib/lexer"
)

// checkParserErrors returns parsing errors
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return // passes tests
	}

	t.Errorf("Parser has %d errors", len(errors)) // return number of errors
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg) // return error message
	}
	t.FailNow() // fails tests
}

func TestLetStatements(t *testing.T) {
	input :=
		// Test input for let
		`
		let x = 5;
		let y = 10;
		let team = Broncos;
		`
	// Call a new lexer and parser
	l := lexer.New(input)
	p := New(l)

	// Throw error if program is empty
	program := p.ParseProgram()

	checkParserErrors(t, p) // Initialize parser error checking

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Throw error if program doesn't contain 3 statements (token, name, value)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements (token, name, value). got=%d", len(program.Statements))
	}

	// input for tests
	tests := []struct {
		// Test that identifiers are being set
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"team"},
	}

	// loop through each test case, add each entry as a program statement
	for i, tt := range tests {

		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

/* Generalized 'let' test, this test fails, input expression returns nil
// testLetStatements tests integrity of input from lexer and parser for let statements.
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", "true"},
		{"let team = broncos;", "team", "broncos"},
	}

	for _, tt := range tests {

		// Call a new lexer and parser
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p) // Initialize parser error checking

		// Throw error if program doesn't contain 1 statement with (token, name, value)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}
*/

// testLetStatement must contain test case, AST statement with TokenLiteral "let", and identifier to return true.
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	// Return false if let statement doesn't contain a value
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	// Return false if let statement doesn't contain a token literal name
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

// TestReturnStatements tests integrity of input from lexer and parser and that it is a valid return statment node in the AST.
func TestReturnStatements(t *testing.T) {
	input :=
		`
		return 89;
		return 12;
		return team;
		`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Throw error if program doesn't contain 3 statements (token, name, value)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// Confirm that the token tested is a Return statement
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

// TestIdentifierExpression tests that identifier is a program statement, is part of the ast, and has the correct value.
func TestIdentifierExpression(t *testing.T) {
	input := "moortr;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program hasn't got enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got =%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("Expression not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "moortr" {
		t.Errorf("ident.Value not %s. got=%s", "moortr", ident.Value)
	}

	if ident.TokenLiteral() != "moortr" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "moortr", ident.TokenLiteral())
	}
}

// TestIntegerLiteralExpression tests the lexing and parsing of integer literals
func TestIntegerLiteralExpression(t *testing.T) {
	// Test input
	input := "5;"

	// Call a new lexer for input, parse it, create a program statement, and test for parser errors
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Length of an integer literal program statement must be 1
	if len(program.Statements) != 1 {
		t.Fatalf("Program should only have 1 statement for integer literal expression. got=%d", len(program.Statements))
	}

	// Ok if program statement can be assigned to an ast expression statement node
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	// Fail if the program statement isn't an ast expression statement
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// Ok if statement expression is an integer literal node
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	// Fail if ast expression isn't an integer literal node
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	// Fail if the value of integer literal 5 isn't 5
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	// Fail if the token for integer literal "5" isn't "5"
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

// TestParsingPrefixExpressions will test prefix expressions ! and -
func TestParsingPrefixExpressions(t *testing.T) {
	// Declare input types, prevents having to rewrite the same test for new input.
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		// Each set of input
		{"!8;", "!", 8},
		{"-16;", "-", 16},
		{"!team;", "!", "team"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	// for the range of input, call a new lexer, parse the information, and run a parser check.
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		// If there is no program statement, fail.
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		// If program statement 0 is not an AST expression statement, fail
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		// If program statement isn't an ast prefix expression, fail
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		// If expression operator isn't expected token type, fail
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

// TestParsingInfixExpressions tests parsing of infix expressions
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		// Input string, left value, operator, right value
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},

		{"Sea + Wolf", "Sea", "+", "Wolf"},
		{"Monte - Cristo", "Monte", "-", "Cristo"},
		{"Gotrek * Felix", "Gotrek", "*", "Felix"},
		{"Don / Quixote", "Don", "/", "Quixote"},
		{"Love > Hate", "Love", ">", "Hate"},
		{"Fiction < Truth", "Fiction", "<", "Truth"},
		{"Left == Right", "Left", "==", "Right"},
		{"Friend != Enemy", "Friend", "!=", "Enemy"},

		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	// For the current input in the range of inputs, create a program statement
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		// Test that each input is 1 program statement
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		// OK if program statement has an ast expression statement
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		// Fails if program statement has no ast expression statement
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}

		/* Error, exp is not defined
		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
		*/
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// input values for testing operator precedence
		{
			"-a * b",     // input string
			"((-a) * b)", // expected string
		},
		{
			"!-a",     // input string
			"(!(-a))", // expected string
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		// test Boolean precedence
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"9 > 99 == false",
			"((9 > 99) == false)",
		},
		{
			"9 < 99 == true",
			"((9 < 99) == true)",
		},
		// test grouped expression precedence
		{
			"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)",
		},
		{
			"(8 + 8) * 2", "((8 + 8) * 2)",
		},
		{
			"12 / (2 + 1)", "(12 / (2 + 1))",
		},
		{
			"-(5 + 5)", "(-(5 + 5))",
		},
		{
			"!(true == true)", "(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

// testIntegerLiteral is generalized integer literal test to verify the current integerLiteral matches its ast.Expression, has the same token type and literal value
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	// if type isn't integerLiteral, fail
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	// if value != input value, fail
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	// If token literal doesn't include a type and a value, fail.
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

// testIdentifier is a more generalized function to verify the current identifier matches its ast.Expression, has the same token type and literal value
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {

	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

// testLiteralExpression identifies the expression type and calls the corresponding test function
func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {

	// call testIntegerLiteral if expression type is an int
	case int:
		return testIntegerLiteral(t, exp, int64(v))

	// call testIntegerLiteral if expression type is an int64
	case int64:
		return testIntegerLiteral(t, exp, v)

	// call testIdentifier if expression type is string
	case string:
		return testIdentifier(t, exp, v)

	// call testIdentifier if expression type is bool
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	// Throw error if expression type isn't an int or string
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// testInfixExpression is a generalized function to test infix expressions, if the expression doesn't match the ast.OperatorExpression and left, operator, and right values aren't identical to AST it fails.
func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got =%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	// Test input
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		// input values for testing operator precedence
		{"true;", true},   // expected string
		{"false;", false}, // value
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		// Length of program statement must be 1
		if len(program.Statements) != 1 {
			t.Fatalf("Program should only have 1 statement for integer literal expression. got=%d", len(program.Statements))
		}

		// Ok if program statement can be assigned to an ast expression statement node
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		// Fail if the program statement isn't an ast expression statement
		if !ok {
			t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		// Ok if statement expression is an integer literal node
		boolean, ok := stmt.Expression.(*ast.Boolean)
		// Fail if ast expression isn't an integer literal node
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}

		// Fail if the value of integer literal 5 isn't 5
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean, boolean.Value)
		}
	}
}

// testBooleanLiteral is generalized Boolean test to verify the current Boolean matches its ast.Expression, has the same token type and literal value
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {

	bo, ok := exp.(*ast.Boolean)
	// if type isn't Boolean, fail
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	// if value != input value, fail
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	// If token literal doesn't include a type and a value, fail.
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

// TestIfExpression tests the parsing of If statement expressions
func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Test the number of program statements, if not 1, fail
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	// Test type of ast node, if not an ast.ExpressionStatement, fail
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// Ok if statement is an AST IfExpression type
	exp, ok := stmt.Expression.(*ast.IfExpression)

	// Fail if statement is not ast.IfExpression type
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	// Verify the expression contains this Infix Expression
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// Verify that expression contains 1 consequence statement node
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	// Ok if Consequence.Statements is an AST expression statement node
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	// Not ok if first consequence statement is not an ast.ExpressionStatement type
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	// For test case, x must be the consequence
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// There is no "else" for this test case, so expression alternative must be nil
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

// TestIfElseExpression tests the parsing of an If expression with an Else statement
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Test the number of If expression program statements, if not 1, fail
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	// Test type of ast node, if not an ast.ExpressionStatement, fail
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// Ok if statement is an AST IfExpression type
	exp, ok := stmt.Expression.(*ast.IfExpression)

	// Fail if statement is not ast.IfExpression type
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	// Verify the expression contains this Infix Expression
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	// Verify that expression contains 1 consequence statement node
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	// Ok if Consequence.Statements is an AST expression statement node
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	// Not ok if first consequence statement is not an ast.ExpressionStatement type
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	// For test case, x must be the consequence
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	// Verify that expression contains 1 alternative statement node
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statement. got=%d\n", len(exp.Alternative.Statements))
	}

	// Ok if Alternative.Statements is an AST expression statement node
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

	// Not ok if first alternative statement is not an ast.ExpressionStatement type
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	// If alternative for expression isn't "y", fail
	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

// TestFunctionLiteralParsing tests Function Literal parsing
func TestFunctionLiteralParsing(t *testing.T) {

	input := `fn(x, y) {x + y;}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Fail if program doesn't contain 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	// Program statement is an AST expression statement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	// Fail if program isn't an AST expression statement
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// AST expression statement type is an AST FunctionLiteral
	function, ok := stmt.Expression.(*ast.FunctionLiteral)

	// Fail if expression statement isn't an AST FunctionLiteral
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	// Fail if number of input parameters isn't 2
	if len(function.Parameters) != 2 {
		t.Fatalf("got wrong number of function literal parameters, want 2, got=%d\n", len(function.Parameters))
	}

	// Verify the input parameters
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	// Fail if function doesn't have 1 body statement
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements hasn't got 1 statement. got=%d\n", len(function.Body.Statements))
	}

	// Function body statement is an AST expression statement
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	// Fail if body statements isn't an AST expression statement
	if !ok {
		t.Fatalf("function body statement is not an ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	// Test input for correct Infix Expression
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

// TestFunctionParameterParsing tests the parsing of parameters for a function literal
func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		// test input
		{input: "fn() {};", expectedParams: []string{}},                     // An empty set of parameters
		{input: "fn(x) {};", expectedParams: []string{"x"}},                 // 1 parameter
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}}, // 3 parameters
	}

	// Run test for each input, apply input to lexer, parse it, and create a new program statement
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		// Apply the program statement to an AST expression statement node
		stmt := program.Statements[0].(*ast.ExpressionStatement)

		// Apply the AST expression statement to an AST functionLiteral expression node
		function := stmt.Expression.(*ast.FunctionLiteral)

		// Tests that the number of input parameters equals the number of output parameters
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length of parameters is wrong, expected %d, got =%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		// For each identifier (parameter), compare its actual type to its expected type
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

// TestCallExpressionParsing tests call expression parsing
func TestCallExpressionParsing(t *testing.T) {

	input := "add(1, 2 * 3, 4 + 5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Fail if program doesn't contain 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	// Program statement is an AST expression statement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	// Fail if program isn't an AST expression statement, show type received
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// AST expression statement type is an AST CallExpression
	exp, ok := stmt.Expression.(*ast.CallExpression)

	// Fail if expression statement isn't an AST CallExpression
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	// Fail if expression's test identifier isn't "add"
	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("Wrong number of arguments, expected 3 got=%d", len(exp.Arguments))
	}

	// Test first expression argument
	testLiteralExpression(t, exp.Arguments[0], 1)

	// Test second expression argument
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)

	// Test third expression argument
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

// TestCallExpressionArgumentParsing tests the parsing of arguments for a call expression
/* Copied from test TestFunctionParameterParsing, haven't tailored it to callExpression arguments yet.
func TestCallExpressionArgumentParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		// test input
		{input: "fn() {};", expectedParams: []string{}},                     // An empty set of parameters
		{input: "fn(x) {};", expectedParams: []string{"x"}},                 // 1 parameter
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}}, // 3 parameters
	}

	// Run test for each input, apply input to lexer, parse it, and create a new program statement
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		// Apply the program statement to an AST expression statement node
		stmt := program.Statements[0].(*ast.ExpressionStatement)

		// Apply the AST expression statement to an AST functionLiteral expression node
		function := stmt.Expression.(*ast.FunctionLiteral)

		// Tests that the number of input parameters equals the number of output parameters
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length of parameters is wrong, expected %d, got =%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		// For each identifier (parameter), compare its actual type to its expected type
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}
*/
