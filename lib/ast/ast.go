/*
Abstract Syntax Tree (AST) for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

package ast

import (
	"bytes"
	"strings"

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

// StringLiteral structure for a String literal expression
type StringLiteral struct {
	Token token.Token
	Value string
}

// StringLiteral assigned to AST expression node
func (sl *StringLiteral) expressionNode() {}

// TokenLiteral contains the literal type of StringLiteral
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String writing function for StringLiteral
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
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

// Boolean structure for Boolean values
type Boolean struct {
	Token token.Token
	Value bool
}

// expressionNode receives Boolean to create an AST node
func (b *Boolean) expressionNode() {}

// TokenLiteral receives Boolean for tokenization
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// Boolean is sent to String function for documentation
func (b *Boolean) String() string {
	return b.Token.Literal
}

// IfExpression structure for If statements
type IfExpression struct {
	Token       token.Token     // The 'if' token
	Condition   Expression      // The condition of the If expression that determines the return value.
	Consequence *BlockStatement // The primary consequence
	Alternative *BlockStatement // The alternative consequence
}

// expressionNode receives the IfExpression to create an AST node
func (ie *IfExpression) expressionNode() {}

// TokenLiteral receives the IfExpression to tokenize
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String receives the IfExpression for documentation and testing
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// BlockStatement is a structure for consequences and alternatives of If statements
type BlockStatement struct {
	Token      token.Token // The { token
	Statements []Statement // An array of If statements
}

// statementNode receives the BlockStatement to create an AST node
func (bs *BlockStatement) statementNode() {}

// TokenLiteral receives the BlockStatement to tokenize
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String receives the BlockStatement for documentation and testing purposes
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// FunctionLiteral structure defines a function
type FunctionLiteral struct {
	Token      token.Token     // The 'fn' token
	Parameters []*Identifier   // Function parameters (a,b,c)
	Body       *BlockStatement // Function statement
}

// expressionNode assign an AST node to FunctionLiteral
func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral returns the FunctionLiteral's token value
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String appends each function literal parameter value (an AST identifier), adds parentheses and separates by comma
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters { // *ast.Identifiers
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression structure for Call Expression AST Node, DoorKey example: 'add(2, 3)' , or 'callsFunction(2, 3, fn(x + y) {x + y;};' , 'out.WriteString(strings.Join(args, ","))' # Golang expression call
type CallExpression struct {
	Token     token.Token  // The '(' token
	Function  Expression   // Identifier or function literal
	Arguments []Expression // Arguments are expressions
}

// ExpressionNode creates an AST expression node for the CallExpression
func (ce *CallExpression) expressionNode() {}

// TokenLiteral returns the Token value for Call expression
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String appends each CallExpression argument to a string, adds parentheses, and separates by comma.
func (ce *CallExpression) String() string {
	var out bytes.Buffer // Assign "out" to byte(s) with buffer?

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}

// ArrayLiteral structure for an array, an ordered list of any type, separated by commas, enclosed by brackets.
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

// ExpressionNode creates an AST expression node for ArrayLiterals
func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral returns the token value for array literal
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// String loops through the ArrayLiteral elements and appends each to a string separated by commas, enclosed by brackets
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
