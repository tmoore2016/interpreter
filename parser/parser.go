/*
Parser for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package parser

import (
	"github.com/tmoore2016/interpreter/ast"
	"github.com/tmoore2016/interpreter/lexer"
	"github.com/tmoore2016/interpreter/token"
)

// Parser structure to contain current Token and next Token from Lexer
type Parser struct {
	l         *lexer.Lexer // l is the pointer
	curToken  token.Token  // current token
	peekToken token.Token  // next token
}

// New Parser for lexer tokens
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken() // set curToken
	p.nextToken() // set peekToken

	return p
}

// nextToken increments to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses the tokens to create the AST
func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
