package object

// NewEnvironment returns an Environment ready for use
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Environment maps identifiers to values
type Environment struct {
	store map[string]Object
}

// Get returns the value for the passed identifier
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set binds an identifier to a value
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
