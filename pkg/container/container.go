// Package container is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package container

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// binding holds a binding resolver and an instance (for singleton bindings).
type binding struct {
	resolver interface{} // resolver function that creates the appropriate implementation of the related abstraction
	instance interface{} // instance stored for reusing in singleton bindings
}

// resolve creates an appropriate implementation of the related abstraction
func (b binding) resolve(c Container) (interface{}, error) {
	if b.instance != nil {
		return b.instance, nil
	}

	return c.invoke(b.resolver)
}

// Container holds all of the declared bindings
type Container map[reflect.Type]map[string]binding

// New creates a new instance of the Container
func New() Container {
	return make(Container)
}

// bind maps an abstraction to a concrete and sets an instance if it's a singleton binding.
func (c Container) bind(resolver interface{}, name string, singleton bool) error {
	reflectedResolver := reflect.TypeOf(resolver)
	if reflectedResolver.Kind() != reflect.Func {
		return errors.New("container: the resolver must be a function")
	}

	for i := 0; i < reflectedResolver.NumOut(); i++ {
		if _, exist := c[reflectedResolver.Out(i)]; !exist {
			c[reflectedResolver.Out(i)] = make(map[string]binding)
		}

		if singleton {
			instance, err := c.invoke(resolver)
			if err != nil {
				return err
			}

			c[reflectedResolver.Out(i)][name] = binding{resolver: resolver, instance: instance}
		} else {
			c[reflectedResolver.Out(i)][name] = binding{resolver: resolver}
		}
	}

	return nil
}

// invoke calls a function and returns the yielded value.
// It only works for functions that return a single value.
func (c Container) invoke(function interface{}) (interface{}, error) {
	args, err := c.arguments(function)
	if err != nil {
		return nil, err
	}

	if reflect.TypeOf(function).NumOut() == 1 {
		return reflect.ValueOf(function).Call(args)[0].Interface(), nil
	} else if reflect.TypeOf(function).NumOut() == 2 {
		values := reflect.ValueOf(function).Call(args)
		return values[0].Interface(), values[1].Interface().(error)
	}

	return nil, errors.New("container: resolver function signature is invalid")
}

// arguments returns container-resolved arguments of a function.
func (c Container) arguments(function interface{}) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)

		if concrete, exist := c[abstraction][""]; exist {
			instance, _ := concrete.resolve(c)

			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, errors.New("container: no concrete found for: " + abstraction.String())
		}
	}

	return arguments, nil
}

// Singleton binds an abstraction to concrete for further singleton resolves.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c Container) Singleton(resolver interface{}) error {
	return c.bind(resolver, "", true)
}

// NamedSingleton binds like the Singleton method but for named bindings.
func (c Container) NamedSingleton(name string, resolver interface{}) error {
	return c.bind(resolver, name, true)
}

// Transient binds an abstraction to concrete for further transient resolves.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c Container) Transient(resolver interface{}) error {
	return c.bind(resolver, "", false)
}

// NamedTransient binds like the Transient method but for named bindings.
func (c Container) NamedTransient(name string, resolver interface{}) error {
	return c.bind(resolver, name, false)
}

// Reset deletes all the existing bindings and empties the container instance.
func (c Container) Reset() {
	for k := range c {
		delete(c, k)
	}
}

// Call takes a function (receiver) with one or more arguments of the abstractions (interfaces).
// It invokes the function (receiver) and passes the related implementations.
func (c Container) Call(function interface{}) error {
	receiverType := reflect.TypeOf(function)
	if receiverType == nil || receiverType.Kind() != reflect.Func {
		return errors.New("container: invalid function")
	}

	arguments, err := c.arguments(function)
	if err != nil {
		return err
	}

	reflect.ValueOf(function).Call(arguments)

	return nil
}

// Resolve takes an abstraction (interface reference) and fills it with the related implementation.
func (c Container) Resolve(abstraction interface{}) error {
	return c.NamedResolve(abstraction, "")
}

// NamedResolve resolves like the Resolve method but for named bindings.
func (c Container) NamedResolve(abstraction interface{}, name string) error {
	receiverType := reflect.TypeOf(abstraction)
	if receiverType == nil {
		return errors.New("container: invalid abstraction")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()

		if concrete, exist := c[elem][name]; exist {
			if instance, err := concrete.resolve(c); err == nil {
				reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
				return nil
			} else {
				return err
			}
		}

		return errors.New("container: no concrete found for: " + elem.String())
	}

	return errors.New("container: invalid abstraction")
}

// Fill takes a struct and resolves the fields with the tag `container:"inject"`
func (c Container) Fill(structure interface{}) error {
	receiverType := reflect.TypeOf(structure)
	if receiverType == nil {
		return errors.New("container: invalid structure")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()
		if elem.Kind() == reflect.Struct {
			s := reflect.ValueOf(structure).Elem()

			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)

				if t, exist := s.Type().Field(i).Tag.Lookup("container"); exist {
					var name string

					if t == "type" {
						name = ""
					} else if t == "name" {
						name = s.Type().Field(i).Name
					} else {
						return errors.New(
							fmt.Sprintf("container: %v has an invalid struct tag", s.Type().Field(i).Name),
						)
					}

					if concrete, exist := c[f.Type()][name]; exist {
						instance, _ := concrete.resolve(c)

						ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
						ptr.Set(reflect.ValueOf(instance))

						continue
					}

					return errors.New(fmt.Sprintf("container: cannot resolve %v field", s.Type().Field(i).Name))
				}
			}

			return nil
		}
	}

	return errors.New("container: invalid structure")
}
