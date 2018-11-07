/*
Object package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

// interpreter\object\object.go

package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tmoore2016/interpreter/lib/ast"
)

// ObjectType represents the Doorkey data types
type ObjectType string

// Strings for Doorkey data types
const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE" // An object for return values
	FUNCTION_OBJ     = "FUNCTION"
	ERROR_OBJ        = "ERROR"
)

// Object represents each data type with a type and value
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer type object.Integer
type Integer struct {
	Value int64
}

// Inspect AST Integer node and return integer value
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type Integer ObjectType
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// Boolean struct wraps a bool value in object.Boolean
type Boolean struct {
	Value bool
}

// Type Boolean ObjectType
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// Inspect AST Boolean node and return a bool
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null is an empty struct
type Null struct{}

// Type Null ObjectType
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// Inspect AST empty node and return a null
func (n *Null) Inspect() string {
	return "null"
}

// ReturnValue structure for Return value objects
type ReturnValue struct {
	Value Object
}

// Type ReturnValue Object
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

// Inspect ReturnValue object value
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Structure for Function object
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // A pointer to the particular environment
}

// Identify the Function object type
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

// Inspect Function and parameters, recreate function string.
func (f *Function) Inspect() string {
	var out bytes.Buffer

	// get function parameters as array and string.
	params := []string{}

	for _, p := range f.Parameters {

		params = append(params, p.String())
	}

	// Adds the function notation, parameters, and function body to the object
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Error structure for error message objects
type Error struct {
	Message string
}

// Type of object: ERROR_OBJ
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

// Inspect Error returns error message (ERROR_OBJ value)
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
