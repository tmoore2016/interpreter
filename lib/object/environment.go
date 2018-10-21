/*
Environment package for
Doorkey, a Monkey Derivative
by Travis Moore
By following "Writing an Interpreter in Go" by Thorsten Ball, https://interpreterbook.com/
*/

// interpreter\object\environment.go

package object

// NewEnvironment creates a hash table (map) that associates strings with object, like a let statement name with its value.
func NewEnvironment() *Environment {

	s := make(map[string]Object)

	return &Environment{store: s}
}

// Environment structure is a hash table that associates a string (name) with an object
type Environment struct {
	store map[string]Object
}

// Get returns an object if the name is associated with an environment (map)
func (e *Environment) Get(name string) (Object, bool) {

	obj, ok := e.store[name]

	return obj, ok
}

// Set associates a name with an object
func (e *Environment) Set(name string, val Object) Object {

	e.store[name] = val

	return val
}
