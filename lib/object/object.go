/*
Object package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

// interpreter\object\object.go

package object

import (
	"fmt"
)

// ObjectType represents the Doorkey data types
type ObjectType string

// Strings for Doorkey data types
const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE" // An object for returns
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
