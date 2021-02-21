// Package container is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package container

import (
	"errors"
	"reflect"
)

// binding keeps a binding resolver and an instance (for singleton bindings).
type binding struct {
	resolver interface{} // resolver function that creates the appropriate implementation of the related abstraction
	instance interface{} // instance stored for reusing in singleton bindings
}

// resolve will create the appropriate implementation of the related abstraction
func (b binding) resolve(c Container) (interface{}, error) {
	if b.instance != nil {
		return b.instance, nil
	}

	return c.invoke(b.resolver)
}

// Container is the repository of bindings
type Container map[reflect.Type]binding

// NewContainer creates a new instance of Container
func NewContainer() Container {
	return make(Container)
}

// bind will map an abstraction to a concrete and set instance if it's a singleton binding.
func (c Container) bind(resolver interface{}, singleton bool) error {
	reflectedResolver := reflect.TypeOf(resolver)
	if reflectedResolver.Kind() != reflect.Func {
		return errors.New("the resolver must be a function")
	}

	for i := 0; i < reflectedResolver.NumOut(); i++ {
		if singleton {
			instance, err := c.invoke(resolver)
			if err != nil {
				return err
			}

			c[reflectedResolver.Out(i)] = binding{resolver: resolver, instance: instance}
		} else {
			c[reflectedResolver.Out(i)] = binding{resolver: resolver}
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
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)

		if concrete, ok := c[abstraction]; ok {
			instance, err := concrete.resolve(c)
			if err != nil {
				return nil, err
			}

			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, errors.New("no concrete found for the abstraction: " + abstraction.String())
		}
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

// Reset will reset the container and remove all the existing bindings.
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
	receiverType := reflect.TypeOf(receiver)
	if receiverType == nil {
		return errors.New("cannot detect type of the receiver")
	}

	if receiverType.Kind() == reflect.Ptr {
		abstraction := receiverType.Elem()

		if concrete, ok := c[abstraction]; ok {
			instance, err := concrete.resolve(c)
			if err != nil {
				return err
			}

			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))

			return nil
		}

		return errors.New("no concrete found for the abstraction: " + abstraction.String())
	}

	if receiverType.Kind() == reflect.Func {
		arguments, err := c.arguments(receiver)
		if err != nil {
			return err
		}

		reflect.ValueOf(receiver).Call(arguments)

		return nil
	}

	return errors.New("the receiver must be either a reference or a callback")
}
