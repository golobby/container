package container

import (
	"reflect"
)

// binding holds a resolver, an instance and the resolving information
type binding struct {
	resolver  interface{}
	singleton bool
	instance  interface{}
}

// resolve will return the concrete of related abstraction
func (b binding) resolve() interface{} {
	if b.singleton {
		return b.instance
	}

	return resolve(b.resolver)
}

// resolve will invoke the given function and return the concrete
func resolve(resolver interface{}) interface{} {
	return reflect.ValueOf(resolver).Call([]reflect.Value{})[0].Interface()
}

// Container is the IoC container that holds the bindings
var container = map[string]binding{}

// bind will bind a concrete to an abstraction
func bind(resolver interface{}, singleton bool, instance interface{}) {
	if reflect.TypeOf(resolver).Kind() != reflect.Func {
		panic("the argument passed to Singleton()/Transient() is not a function")
	}

	if reflect.TypeOf(resolver).NumOut() != 1 {
		panic("The resolver must only return with abstraction type")
	}

	container[reflect.TypeOf(resolver).Out(0).String()] = binding{
		singleton: singleton,
		resolver:  resolver,
		instance:  instance,
	}
}

// Singleton will bind a singleton concrete to an abstraction
func Singleton(function interface{}) {
	bind(function, true, resolve(function))
}

// Transient will bind a transient concrete to an abstraction
func Transient(function interface{}) {
	bind(function, false, nil)
}

// Make will resolve the given abstraction and return related concrete
func Make(function interface{}) {
	if reflect.TypeOf(function).Kind() != reflect.Func {
		panic("the argument passed to Make() is not a function")
	}

	if reflect.TypeOf(function).NumIn() != 1 {
		panic("Make() takes one argument which is the abstraction")
	}

	abstraction := reflect.TypeOf(function).In(0).String()

	if concrete, ok := container[abstraction]; ok {
		arguments := []reflect.Value{reflect.ValueOf(concrete.resolve())}
		reflect.ValueOf(function).Call(arguments)
	} else {
		panic("There is no concrete for " + abstraction)
	}
}
