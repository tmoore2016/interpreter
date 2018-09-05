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

	if len(errors) == 0 { // exit if no errors
		return
	}

	t.Errorf("Parser has %d errors", len(errors)) // return number of errors
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg) // return error message
	}
	t.FailNow()
}

// TestLetStatements tests integrity of input from lexer and parser.
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
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
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

	// Ok to assign the program statement to an ast expression statement node
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

// testIntegerLiteral is a smaller integer literal test
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

// TestParsingPrefixExpressions will test prefix expressions ! and -
func TestParsingPrefixExpressions(t *testing.T) {

	// Declare input types, prevents having to rewrite the same test for new input.
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		// Each set of input
		{"!8;", "!", 8},
		{"-16;", "-", 16},
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

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

// TestParsingInfixExpressions tests parsing of infix expressions
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
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

		// OK if the expression statement (+) is an ast infix expression statement
		exp, ok := stmt.Expression.(*ast.InfixExpression)

		// Fails if the expression statement (+) isn't an ast infix expression statement
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		// Fails if left expression statement isn't an integer literal
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		// Fails if current expression operator != program statement expression operator
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		// Fails if right expression statement isn't an integer value
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// input values for testing operator precedence
		{
			"-a * b",     // expected string
			"((-a) * b)", // actual string
		},
		{
			"!-a",     // expected string
			"(!(-a))", // actual string
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
