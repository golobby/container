// Package container provides an IoC container for Go projects.
// It provides simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.
package container

import (
	"errors"
	"reflect"
)

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func invoke(function interface{}) (interface{}, error) {
	args, err := arguments(function)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(function).Call(args)[0].Interface(), nil
}

// binding keeps a binding resolver and instance (for singleton bindings).
type binding struct {
	resolver interface{} // resolver function
	instance interface{} // instance stored for singleton bindings
}

// resolve will return the concrete of related abstraction.
func (b binding) resolve() (interface{}, error) {
	if b.instance != nil {
		return b.instance, nil
	}

	return invoke(b.resolver)
}

// container is the IoC container that will keep all of the bindings.
var container = map[reflect.Type]binding{}

// bind will map an abstraction to a concrete and set instance if it's a singleton binding.
func bind(resolver interface{}, singleton bool) error {
	resolverTypeOf := reflect.TypeOf(resolver)
	if resolverTypeOf.Kind() != reflect.Func {
		return errors.New("the resolver must be a function")
	}

	for i := 0; i < resolverTypeOf.NumOut(); i++ {
		var instance interface{}
		var err error
		if singleton {
			instance, err = invoke(resolver)
			if err != nil {
				return err
			}
		}

		container[resolverTypeOf.Out(i)] = binding{
			resolver: resolver,
			instance: instance,
		}
	}

	return nil
}

// arguments will return resolved arguments of the given function.
func arguments(function interface{}) ([]reflect.Value, error) {
	functionTypeOf := reflect.TypeOf(function)
	argumentsCount := functionTypeOf.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := functionTypeOf.In(i)

		var instance interface{}
		var err error

		if concrete, ok := container[abstraction]; ok {
			instance, err = concrete.resolve()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("no concrete found for the abstraction: " + abstraction.String())
		}

		arguments[i] = reflect.ValueOf(instance)
	}

	return arguments, nil
}

// Singleton will bind an abstraction to a concrete for further singleton resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Singleton(resolver interface{}) error {
	return bind(resolver, true)
}

// Transient will bind an abstraction to a concrete for further transient resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Transient(resolver interface{}) error {
	return bind(resolver, false)
}

// Reset will reset the container and remove all the bindings.
func Reset() {
	container = map[reflect.Type]binding{}
}

// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func Make(receiver interface{}) error {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		return errors.New("cannot detect type of the receiver, make sure your are passing reference of the object")
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		if concrete, ok := container[abstraction]; ok {
			instance, err := concrete.resolve()
			if err != nil {
				return err
			}

			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))

			return nil
		}

		return errors.New("no concrete found for the abstraction " + abstraction.String())
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments, err := arguments(receiver)
		if err != nil {
			return err
		}

		reflect.ValueOf(receiver).Call(arguments)

		return nil
	}

	return errors.New("the receiver must be either a reference or a callback")
}
