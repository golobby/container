// Package container provides an IoC container for Go projects.
// It provides simple, fluent, performance-in-mid, and easy-to-use interface
// to make dependency injection in GoLang easier.
package container

import (
	"errors"
	"reflect"
)

// binding keeps a binding resolver and instance (for singleton bindings).
type binding struct {
	resolver interface{} // resolver function
	instance interface{} // instance stored for singleton bindings
}

// resolve will return the concrete of related abstraction.
func (b binding) resolve(c Container) (interface{}, error) {
	if b.instance != nil {
		return b.instance, nil
	}

	return c.invoke(b.resolver)
}

// Container is a map of reflect.Type to binding
type Container map[reflect.Type]binding

// NewContainer returns a new instance of Container
func NewContainer() Container {
	return make(Container)
}

// bind will map an abstraction to a concrete and set instance if it's a singleton binding.
func (c Container) bind(resolver interface{}, singleton bool) error {
	var err error
	resolverTypeOf := reflect.TypeOf(resolver)
	if resolverTypeOf.Kind() != reflect.Func {
		return errors.New("the resolver must be a function")
	}

	for i := 0; i < resolverTypeOf.NumOut(); i++ {
		var instance interface{}
		if singleton {
			instance, err = c.invoke(resolver)
			if err != nil {
				return err
			}
		}

		c[resolverTypeOf.Out(i)] = binding{
			resolver: resolver,
			instance: instance,
		}
	}

	return nil
}

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func (c Container) invoke(function interface{}) (interface{}, error) {
	args, err := c.arguments(function)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(function).Call(args)[0].Interface(), nil
}

// arguments will return resolved arguments of the given function.
func (c Container) arguments(function interface{}) ([]reflect.Value, error) {
	var err error
	functionTypeOf := reflect.TypeOf(function)
	argumentsCount := functionTypeOf.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := functionTypeOf.In(i)

		var instance interface{}

		if concrete, ok := c[abstraction]; ok {
			instance, err = concrete.resolve(c)
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
func (c Container) Singleton(resolver interface{}) error {
	return c.bind(resolver, true)
}

// Transient will bind an abstraction to a concrete for further transient resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func (c Container) Transient(resolver interface{}) error {
	return c.bind(resolver, false)
}

// Reset will reset the container and remove all the bindings.
func (c Container) Reset() {
	for k := range c {
		delete(c, k)
	}
}

// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func (c Container) Make(receiver interface{}) error {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		return errors.New("cannot detect type of the receiver")
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		if concrete, ok := c[abstraction]; ok {
			instance, err := concrete.resolve(c)
			if err != nil {
				return err
			}
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return nil
		}

		return errors.New("no concrete found for the abstraction " + abstraction.String())
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments, err := c.arguments(receiver)
		if err != nil {
			return err
		}
		reflect.ValueOf(receiver).Call(arguments)
		return nil
	}

	return errors.New("the receiver must be either a reference or a callback")
}
