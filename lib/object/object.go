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

// BuiltinFunction type is a Go function that can be called from Doorkey
type BuiltinFunction func(args ...Object) Object

// Strings for Doorkey data types
const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE" // An object for return values
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
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

// String type object.String
type String struct {
	Value string
}

// Type string ObjectType
func (s *String) Type() ObjectType {
	return STRING_OBJ
}

// Inspect AST string node and return its value
func (s *String) Inspect() string {
	return s.Value
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

// Function object structure
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // A pointer to the particular environment
}

// Type check for Function object
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

// Builtin structure for callable Go functions
type Builtin struct {
	Fn BuiltinFunction
}

// Type check for BUILTIN_OBJ
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// Inspect string, return as a Builtin function
func (b *Builtin) Inspect() string { return "builtin function" }

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
