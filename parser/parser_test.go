/*
Parser tests for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package parser

import (
	"testing"

	"github.com/tmoore2016/interpreter/ast"
	"github.com/tmoore2016/interpreter/lexer"
)

// TestLetStatements tests Let statement
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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}

// testLetStatment must contain test case, AST statement with TokenLiteral "let", and identifier to return true.
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
