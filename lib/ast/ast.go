/*
Abstract Syntax Tree (AST) for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package ast

import "github.com/tmoore2016/interpreter/lib/token"

// Node in AST implements the node interface, providing a TokenLiteral() that returns the literal value its associated with for debugging and testing.
type Node interface {
	TokenLiteral() string
}

// Statement doesn't return values.
type Statement interface {
	Node
	statementNode()
}

// Expression returns values.
type Expression interface {
	Node
	expressionNode()
}

// Program contains input statements
type Program struct { // struct = type with named fields
	Statements []Statement
}

// TokenLiteral is the root node of the AST
// Doorkey statements are contained in program.statements
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return "Nil statement."
	}
}

// LetStatement prepares a let statement node
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier // call Identifier() for IDENT
	Value Expression  // literal values
}

// statementNode contains LetStatement
func (ls *LetStatement) statementNode() {
}

// TokenLiteral returns the literal values of LetStatement's token
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// ReturnStatement prepares a return statement node
type ReturnStatement struct {
	Token       token.Token // the return token
	ReturnValue Expression
}

// statementNode contains ReturnStatement
func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the literal values of ReturnStatement's token
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// Identifier returns the identity value of token
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// expressionNode contains Identifier
func (i *Identifier) expressionNode() {
}

// TokenLiteral contains Literal value
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
