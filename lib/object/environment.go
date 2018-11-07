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

	return &Environment{store: s, outer: nil}
}

// Environment structure is a hash table that associates a string (name) with an object. The outer environment allows one environment to wrap another.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get returns an object if the name is associated with an environment (map)
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set associates a name with an object
func (e *Environment) Set(name string, val Object) Object {

	e.store[name] = val

	return val
}

// NewEnclosedEnvironment allows one environment to wrap another.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}
