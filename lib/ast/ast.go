/*
Abstract Syntax Tree (AST) for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package ast

import (
	"bytes"

	"github.com/tmoore2016/interpreter/lib/token"
)

// Node in AST implements the node interface, providing a TokenLiteral() that returns its associated literal value for debugging and testing.
type Node interface {
	TokenLiteral() string
	String() string // Each node will write itself as a string for debugging
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

// Identifier returns the identity value of token
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// expressionNode contains Identifier
func (i *Identifier) expressionNode() {
}

// TokenLiteral contains Literal type
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String function for Identifier, return value
func (i *Identifier) String() string {
	return i.Value
}

// Create a buffer and write the value of each new program statement into a string for debugging
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// TokenLiteral is the root node of the AST
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral() // Doorkey statements are contained in program.statements
	} else {
		return "Nil statement."
	}
}

// LetStatement prepares a Let statement node
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier // call Identifier() for IDENT
	Value Expression  // literal type
}

// statementNode contains LetStatement
func (ls *LetStatement) statementNode() {
}

// TokenLiteral returns the literal type of LetStatement's token
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String writing function for let statement
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ReturnStatement prepares a Return statement node
type ReturnStatement struct {
	Token       token.Token // the return token
	ReturnValue Expression
}

// statementNode contains ReturnStatement
func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the literal type of ReturnStatement's token
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String writing function for return statement
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement prepares an Expression statement node type
type ExpressionStatement struct {
	Token      token.Token // This field contains the first token of the expression
	Expression Expression  // This field contains the expression
}

// statementNode contains ExpressionStatement
func (es *ExpressionStatement) statementNode() {}

// TokenLiteral contains the literal type of ExpressionStatement
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String writing function for expression statement
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral structure for an integer literal expression
type IntegerLiteral struct {
	Token token.Token
	Value int64 // value isn't a string
}

// IntegerLiteral is assigned to an AST expression node
func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral contains the literal type of integer literal
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String writing function for IntegerLiteral
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// PrefixExpression structure for a prefix expression
type PrefixExpression struct {
	Token    token.Token // The prefix token
	Operator string      // ! or -
	Right    Expression  // The expression to the right of the operator
}

// PrefixExpression is assigned to an AST expression node
func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral contains the literal type of prefix expression
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String adds parentheses between the operator and the operand of a prefix expression
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression structure for an infix expression
type InfixExpression struct {
	Token    token.Token // The infix operator token, '+'
	Left     Expression
	Operator string
	Right    Expression
}

// InfixExpression is assigned to an ast expression node
func (oe *InfixExpression) expressionNode() {} // operator expression

// TokenLiteral contains the literal type of the infix expression
func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.Literal
}

// Write the entire infix expression to a string
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())        // Left operand expression
	out.WriteString(" " + oe.Operator + " ") // Infix operator
	out.WriteString(oe.Right.String())       // Right operand expression
	out.WriteString(")")

	return out.String()
}
