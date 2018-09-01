/*
Parser for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package parser

import (
	"fmt"
	"strconv"

	"github.com/tmoore2016/interpreter/lib/ast"
	"github.com/tmoore2016/interpreter/lib/lexer"
	"github.com/tmoore2016/interpreter/lib/token"
)

// Parser precedence, lowest to highest
const (
	_           int = iota // iota assigns values in ascending order
	LOWEST                 // lowest precedence
	EQUALS                 // ==
	LESSGREATER            // > or <
	SUM                    // +
	PRODUCT                // *
	PREFIX                 // -X or !X
	CALL                   // myFunction(X)
)

// Parser structure, pulls data from lexer
type Parser struct {
	l              *lexer.Lexer                      // l is the pointer
	errors         []string                          // error handling
	curToken       token.Token                       // current token's type
	peekToken      token.Token                       // next token's type
	prefixParseFns map[token.TokenType]prefixParseFn // hash table to compare prefix and infix expressions
	infixParseFns  map[token.TokenType]infixParseFn
}

// New Parser for lexer tokens
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{}, // error handling
	}

	p.nextToken() // set curToken
	p.nextToken() // set peekToken

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) // Initialize prefixParseFns map
	p.registerPrefix(token.IDENT, p.parseIdentifier)           // Register an Identifier parsing function
	p.registerPrefix(token.INT, p.parseIntegerLiteral)         // Register an Integer Literal parsing function
	p.registerPrefix(token.NOT, p.parsePrefixExpression)       // Register a ! prefix expression
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)     // Register a - prefix expression
	return p
}

// parseIdentifier returns the AST identifier and its value, it doesn't advance the token or call nextToken.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError appends errors to message if unexpected token is encountered
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken increments to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// sets the current token
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// sets the next token
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek confirmst that peekToken equals nextToken, or throws peekError
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Prefix and Infix parsing functions set prefix and infix expression nodes
type (
	prefixParseFn func() ast.Expression               // create a prefix expression
	infixParseFn  func(ast.Expression) ast.Expression // puts the prefix expression on the left of the infix expression
)

// registerPrefix adds entries to prefixParseFns map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// ParseProgram parses the tokens to create the root node for the AST
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// If statements aren't empty, append program statements until End of File token.
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatment checks token type to determine statement type
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	// Let statement
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement() // if the statement isn't a let or a return, treat it as an expression (named var).
	}
}

// parseLetStatement creates a let statement node
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// let statement expects an identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Uses the identifier to create an AST identifier node
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// let statment expects a assignment (=)
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Skipping expressions until we encounter a semicolon

	// Stop progressing when a semicolon is encountered
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement creates a return statement node
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: Skipping expressions until we encounter a semicolon

	// Stop progressing when a semicolon is encountered
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement creates an expression node
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) { // The expression statement continues until the next token is a ";".
		p.nextToken()
	}

	return stmt
}

// parseExpression checks if there is a parsing function associated with the current token.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type) // calls noPrefixParseFnError if prefix type is nil
		return nil
	}

	leftExp := prefix() // assigns prefix expression to left expression.

	return leftExp
}

// parseIntegerLiteral parses integer literal expressions, returns the AST identifier and its value, it doesn't advance the token or call nextToken.
func (p *Parser) parseIntegerLiteral() ast.Expression {

	lit := &ast.IntegerLiteral{Token: p.curToken}

	// Convert string value to Int64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// noPrefixParseFnError appends invalid type information for prefix expressions to parser errors
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for %s found", t) // If there isn't a valid prefix expression type, throw an error and return the actual type.
	p.errors = append(p.errors, msg)                               // Append error message to parser errors
}

// parsePrefixExpression parses ! and - prefixes, and their associated expressions
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// Advance parser to next token after prefix
	p.nextToken()

	// Applies the nextToken's ast node to the right side of the prefix expression
	expression.Right = p.parseExpression(PREFIX)

	return expression
}
