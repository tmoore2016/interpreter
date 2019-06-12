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
	"hash/fnv"
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
	ARRAY_OBJ        = "ARRAY"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE" // An object for return values
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
	HASH_OBJ         = "HASH"
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

// Array structure for an array object
type Array struct {
	Elements []Object
}

// Type assigns an Array.Type to an array object
func (ao *Array) Type() ObjectType {
	return ARRAY_OBJ
}

// Inspect loops through the elements of an array object and appends their index ID to each
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
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

// HashKey structure for hash keys. Type is any object type, value is an integer.
// Caching HashKey methods would be a performance optimization
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// HashKey function for boolean comparisons
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// HashKey function for comparing integer values
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey function for comparing string values
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// HashPair structure contains the objects that generated the HashKey, their type and values.
type HashPair struct {
	Key   Object
	Value Object
}

// Hash structure points to the HashKey and the HashPair
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns HASH_OBJ type
func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

// Inspect iterates over hash pairs and returns their key and value as a string.
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Hashable determines whether the type given is suitable for hashing.
type Hashable interface {
	HashKey() HashKey
}
