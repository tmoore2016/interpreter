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
	INDEX                  // array[index]
)

// Assigns parser precedence to tokens
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.DIVIDE:   PRODUCT,
	token.MULTIPLY: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

// Parser structure, pulls data from lexer
type Parser struct {
	l              *lexer.Lexer                      // l is the pointer
	errors         []string                          // error handling
	curToken       token.Token                       // current token
	peekToken      token.Token                       // next token
	prefixParseFns map[token.TokenType]prefixParseFn // hash table to compare prefix and infix expressions
	infixParseFns  map[token.TokenType]infixParseFn
}

// peekPrecedence returns the precedence operator for peek token, defaults to lowest
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// curPrecedence returns the precedence operator for current token, defaults to lowest
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	// Precedence defaults to lowest.
	return LOWEST
}

// New Parser for lexer's tokens
func New(l *lexer.Lexer) *Parser {
	p := &Parser{ // current parser
		l:      l,          // current lexer
		errors: []string{}, // error handling
	}

	p.nextToken() // set curToken
	p.nextToken() // set peekToken

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) // Initialize prefixParseFns map
	p.registerPrefix(token.IDENT, p.parseIdentifier)           // Register an Identifier parsing function
	p.registerPrefix(token.INT, p.parseIntegerLiteral)         // Register an Integer Literal parsing function
	p.registerPrefix(token.STRING, p.parseStringLiteral)       // Register a String Literal expression
	p.registerPrefix(token.NOT, p.parsePrefixExpression)       // Register a ! prefix expression
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)     // Register a - prefix expression
	p.registerPrefix(token.TRUE, p.parseBoolean)               // Register a TRUE prefix expression
	p.registerPrefix(token.FALSE, p.parseBoolean)              // Register a False prefix expression
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)   // Register a ( prefix expression
	p.registerPrefix(token.IF, p.parseIfExpression)            // Register an IF prefix expression
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)   // Register a Function prefix expression
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)      // Register a [ prefix expression for arrays
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)         // Register a { prefix for hash literal expressions

	p.infixParseFns = make(map[token.TokenType]infixParseFn) // Create a hash table of infix expression tokens
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression) // Register a ( infix expression for call expressions

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

// expectPeek confirms that peekToken equals nextToken, or throws peekError
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

// registerInfix adds entries to infixParseFns map
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

// parseStatement checks token type to determine statement type
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

	// let statement expects a assignment (=)
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement creates a return statement node
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	// Stop progressing when a semicolon is encountered
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement creates expression nodes
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	defer untrace(trace("parseExpressionStatement")) // Call parser_tracing to follow this expression

	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST) // First precedence expression statement

	if p.peekTokenIs(token.SEMICOLON) { // The expression statement continues until the next token is a ";"
		p.nextToken()
	}

	return stmt
}

// parseExpression checks if there is a parsing function associated with the current token and assigns it to left expression
func (p *Parser) parseExpression(precedence int) ast.Expression { // Precedence defaults to LOWEST unless a higher precedence is passed from parseInfixExpression

	defer untrace(trace("parseExpression")) // Call parser_tracing to follow this expression, defer doesn't execute until the surrounding function returns.

	prefix := p.prefixParseFns[p.curToken.Type] // Checks if there is a prefixParseFn associated with the token type, (i.e. "1 + 2 + 3;", the 1 is an integer literal expression, so it calls parseIntegerLiteral)

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type) // calls noPrefixParseFnError if prefix type is nil
		return nil
	}

	leftExp := prefix() // assigns prefix expression to left expression

	// Check if the next token is higher precedence than the current left expression, if it is assign the new left expression, continue until the next expression is not higher precedence or a ';'
	// (i.e. "1 + 2 + 3;", the first round "1 +" loops and ast.InfixExpression is "+", ast.IntegerLiteral left is "1", the second round "2 +" it doesn't loop because the first +'s precedence is still applied and astIntegerLiteral right is "2", making the expression "1+2". The third time it loops because the second + is higher precedence than 2 and ast.Infix left is now "1 + 2", parseInfixExpression is called again, it advances the token and "1+2+3;" becomes the left expression.
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { // Token is not a ';' and current left expression precedence is lower than peek precedence
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp // if not an infix expression, left expression
		}

		p.nextToken() // advance to next token

		leftExp = infix(leftExp) // apply the new infix expression to left expression and call parseInfixExpression
	}

	return leftExp
}

// parseIntegerLiteral parses integer literal expressions from parseExpression, returns the AST identifier and its value, converting the string into an integer, it doesn't advance the token or call nextToken
func (p *Parser) parseIntegerLiteral() ast.Expression {

	defer untrace(trace("parseIntegerLiteral")) // Call parser_tracing to follow this expression

	lit := &ast.IntegerLiteral{Token: p.curToken}

	// Convert string value to Int64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64) // call the parser's current token's literal value and convert to integer

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// parseStringLiteral parses String Literal expressions, returns the AST identifier and its value as a single string token.
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

/*
// This is ParseBoolean function from the book. I rewrote this following the parseIntegerLiteral function that converts the string to another type. Good idea?
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}
*/

// parseBoolean parses boolean expressions from parseExpression, returns the AST identifier and its value, converts the string into an integer, and returns the new token type. It doesn't advance the token or call nextToken.
func (p *Parser) parseBoolean() ast.Expression {

	defer untrace(trace("parseBoolean")) // Call parser_tracing to follow this expression

	bo := &ast.Boolean{Token: p.curToken}

	// Convert string value to Boolean
	value, err := strconv.ParseBool(p.curToken.Literal) // call the parser's current token string value and convert to Boolean

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as Boolean", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	bo.Value = value

	return bo
}

// parseArrayLiteral parses elements following an '[' prefix expression through parseExpressionList until the end token ']' is encountered, and returns the list of elements within an ArrayLiteral token.
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

// parseIndexExpression parses index expressions for arrays
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// parseHashLiteral parses hash literal expressions by looping over key-value pairs and calling parseExpression two times for each pair and filling hash.Pairs. If peekToken is }, it returns nil.
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

// parseExpressionList parses through each element in the list, ignoring commas, until it reaches the specified end token (']' for arrays)
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	// Return list if the next token is the end token, for an empty array
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	// Advance to the next token and append the element to the list
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	// If the next token is a comma, advance twice
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	// If there is no end token, return nil
	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parsePrefixExpression parses ! and - prefixes, and their associated expressions
func (p *Parser) parsePrefixExpression() ast.Expression {

	defer untrace(trace("parsePrefixExpression")) // Call parser_tracing to follow this expression

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

// noPrefixParseFnError appends invalid type information for prefix expressions to parser errors
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("Invalid prefix operator, type: %s", t) // If there isn't a valid prefix expression type, throw an error and return the actual type.
	p.errors = append(p.errors, msg)                           // Append error message to parser errors
}

// parseInfixExpression creates an infix expression node
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	defer untrace(trace("parseInfixExpression")) // Call parser_tracing to follow this expression

	expression := &ast.InfixExpression{ // & points the product to ast.InfixExpresssion

		Token:    p.curToken,         // Set token to current token
		Operator: p.curToken.Literal, // set operator to literal
		Left:     left,               // set local left to ast expression left from parsePrefixExpression (i.e. "1 + 2 + 3;" first the 1, then 2, then 1 + 2)
	}

	precedence := p.curPrecedence()                  // saves precedence of the current token, i.e. ("1 + 2 + 3;" the first +)
	p.nextToken()                                    // move to next token
	expression.Right = p.parseExpression(precedence) // add right field to infix expression from parseExpression, (i.e. "1 + 2 + 3;", the 2)

	return expression
}

// parseGroupedExpression parses grouped expressions, "12 / (2+2)" == "(12 / (2+2))"
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression parses IF expressions
func (p *Parser) parseIfExpression() ast.Expression { // Create an AST expression node

	expression := &ast.IfExpression{Token: p.curToken} // Add the current token to an AST If expression node

	if !p.expectPeek(token.LPAREN) { // End if token after "If" isn't a "("
		return nil // expectPeek Returns a parser error if token is the wrong type
	}

	p.nextToken() // Call next token

	expression.Condition = p.parseExpression(LOWEST) // Check for (, ), {, then apply LOWEST precedence to them and continue

	if !p.expectPeek(token.RPAREN) { // End if "If" expression doesn't end with ")"
		return nil // expectPeek Returns a parser error if token is the wrong type
	}

	if !p.expectPeek(token.LBRACE) { // { marks beginning of block statement
		return nil // expectPeek Returns a parser error if token is the wrong type
	}

	expression.Consequence = p.parseBlockStatement() // Apply the expression's consequence from the block statement

	if p.peekTokenIs(token.ELSE) { // If "If" expresion contains an "else", call next token
		p.nextToken()

		if !p.expectPeek(token.LBRACE) { // Next token after "else" should be "{", expectPeek will advance token again if it is
			return nil // expectPeek Returns a parser error if token is the wrong type
		}

		expression.Alternative = p.parseBlockStatement() // Apply the "else alternative" to the block statement
	}

	return expression // Results of If expression
}

// parseBlockStatement parses IF block statements, similar to parseStatement function
func (p *Parser) parseBlockStatement() *ast.BlockStatement { // Create an AST node for block statements
	block := &ast.BlockStatement{Token: p.curToken} // Insert current token into AST node
	block.Statements = []ast.Statement{}            // Insert token into an array of block statement tokens

	p.nextToken() // Call next token

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) { // Continue looping until } or EOF is encountered
		stmt := p.parseStatement() // Parse block statement

		if stmt != nil {
			block.Statements = append(block.Statements, stmt) // Append token to array
		}

		p.nextToken() // Call next token
	}

	return block // Results of block statement
}

// parseFunctionLiterals parses function literals "fn add(a,b){a+b;}"
func (p *Parser) parseFunctionLiteral() ast.Expression {

	lit := &ast.FunctionLiteral{Token: p.curToken} // Create an AST function literal node with current token

	if !p.expectPeek(token.LPAREN) { // End if token after fn name isn't a "("
		return nil // ExpectPeek returns a parser error if token isn't the expected type
	}

	lit.Parameters = p.parseFunctionParameters() // Parse the function's parameters with parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) { // End if token after fn name isn't a "{"
		return nil // ExpectPeek returns a parser error if token isn't the expected type
	}

	lit.Body = p.parseBlockStatement() // Parse the function's body with parseBlockStatement()

	return lit // Final functionLiteral expression
}

// parseFunctionParameters parses function literal expression's parameters
func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{} // Assign to an array of AST identifiers (function parameters)

	// If next token is ")", advance parser to next token and return empty array of parameters
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken() // Advance to token following "("

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // Define a single parameter, assign it to an AST identifier with Token type and value

	identifiers = append(identifiers, ident) // Append individual identifier to the array of identifiers (parameters)

	// When next token is "," -- advance parser twice to move to the token after comma, assign the identifier (parameter)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	// An ")" is expected to follow the parameter list, if this is false, return peek error
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers // Final parameter list
}

// parseCallExpressions parses call expressions
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {

	// Assign to an AST CallExpression node
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	// Parse call expression arguments
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	//exp.Arguments = p.parseCallArguments() // old version

	return exp
}

// parseCallArguments parses call expression arguments
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{} // Put each arg into an array of AST expression nodes

	// If peek token is RPAREN, advance to next token and return arguments.
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	// Advance to token after LPAREN
	p.nextToken()

	// Append current token to arguments array, parse the expression, assign lowest precedence
	args = append(args, p.parseExpression(LOWEST))

	// If peekToken is a comma, advance two tokens and append current token to arguments, parse current token, assign it lowest precedence
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	// If there isn't a RPAREN in call expression, return nil
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
