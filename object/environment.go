package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// NewEnvironment returns an Environment ready for use
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// Environment maps identifiers to values
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get returns the value for the passed identifier
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set binds an identifier to a value
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
