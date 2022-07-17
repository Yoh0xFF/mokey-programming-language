package object

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	rslt, ok := e.store[name]

	if !ok && e.outer != nil {
		rslt, ok = e.outer.Get(name)
	}

	return rslt, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
