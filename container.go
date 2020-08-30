// Package container provides an IoC container for Go projects.
// It provides simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.
package container

import (
	"reflect"
)

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func invoke(function interface{}) interface{} {
	return reflect.ValueOf(function).Call(arguments(function))[0].Interface()
}

// binding keeps a binding resolver and instance (for singleton bindings).
type binding struct {
	resolver interface{} // resolver function
	instance interface{} // instance stored for singleton bindings
}

// resolve will return the concrete of related abstraction.
func (b binding) resolve() interface{} {
	if b.instance != nil {
		return b.instance
	}

	return invoke(b.resolver)
}

// container is the IoC container that will keep all of the bindings.
var container = map[reflect.Type]binding{}

// bind will map an abstraction to a concrete and set instance if it's a singleton binding.
func bind(resolver interface{}, singleton bool) {
	resolverTypeOf := reflect.TypeOf(resolver)
	if resolverTypeOf.Kind() != reflect.Func {
		panic("the resolver must be a function")
	}

	for i := 0; i < resolverTypeOf.NumOut(); i++ {
		var instance interface{}
		if singleton {
			instance = invoke(resolver)
		}

		container[resolverTypeOf.Out(i)] = binding{
			resolver: resolver,
			instance: instance,
		}
	}
}

// arguments will return resolved arguments of the given function.
func arguments(function interface{}) []reflect.Value {
	functionTypeOf := reflect.TypeOf(function)
	argumentsCount := functionTypeOf.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := functionTypeOf.In(i)

		var instance interface{}

		if concrete, ok := container[abstraction]; ok {
			instance = concrete.resolve()
		} else {
			panic("no concrete found for the abstraction: " + abstraction.String())
		}

		arguments[i] = reflect.ValueOf(instance)
	}

	return arguments
}

// Singleton will bind an abstraction to a concrete for further singleton resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Singleton(resolver interface{}) {
	bind(resolver, true)
}

// Transient will bind an abstraction to a concrete for further transient resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Transient(resolver interface{}) {
	bind(resolver, false)
}

// Reset will reset the container and remove all the bindings.
func Reset() {
	container = map[reflect.Type]binding{}
}

// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func Make(receiver interface{}) {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		panic("cannot detect type of the receiver, make sure your are passing reference of the object")
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		if concrete, ok := container[abstraction]; ok {
			instance := concrete.resolve()
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return
		}

		panic("no concrete found for the abstraction " + abstraction.String())
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments := arguments(receiver)
		reflect.ValueOf(receiver).Call(arguments)
		return
	}

	panic("the receiver must be either a reference or a callback")
}
