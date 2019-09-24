package container

import (
	"reflect"
)

// invoke will call the given function (resolver) and return its return value
func invoke(resolver interface{}) interface{} {
	return reflect.ValueOf(resolver).Call([]reflect.Value{})[0].Interface()
}

// binding holds a resolver and its resolving information
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

	return invoke(b.resolver)
}

// container is the IoC container which holds all of the bindings
var container = map[string]binding{}

// bind will bind an abstraction to a concrete
func bind(resolver interface{}, singleton bool, instance interface{}) {
	if reflect.TypeOf(resolver).Kind() != reflect.Func {
		panic("the argument passed to Singleton() or Transient() methods is not a function")
	}

	if reflect.TypeOf(resolver).NumIn() != 0 {
		panic("the resolver function cannot take any argument")
	}

	if reflect.TypeOf(resolver).NumOut() != 1 {
		panic("the resolver function must only return the abstraction type")
	}

	container[reflect.TypeOf(resolver).Out(0).String()] = binding{
		singleton: singleton,
		resolver:  resolver,
		instance:  instance,
	}
}

// Singleton will bind an abstraction to a singleton concrete
// It takes a resolver function which returns the concrete and its return type matches the abstraction
func Singleton(function interface{}) {
	bind(function, true, invoke(function))
}

// Transient will bind an abstraction to a transient concrete
// It takes a resolver function which returns the concrete and its return type matches the abstraction
func Transient(function interface{}) {
	bind(function, false, nil)
}

// Make will resolve the given abstraction and return related concrete
// It takes a function with one argument of the abstraction type,
// the Container invokes the function an pass the related create
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
