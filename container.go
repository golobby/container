// Container package provides an IoC Container for Go projects.
// It provides simple, fluent and easy-to-use interface to make dependency injection in GoLang very easier.
package container

import (
	"reflect"
)

// invoke will call the given resolver function and return its return value
func invoke(resolver interface{}) interface{} {
	return reflect.ValueOf(resolver).Call([]reflect.Value{})[0].Interface()
}

// binding is a struct to holds a resolver and its resolving information
type binding struct {
	resolver interface{}
	instance interface{}
}

// resolve will return the concrete of related abstraction
func (b binding) resolve() interface{} {
	if b.instance != nil {
		return b.instance
	}

	return invoke(b.resolver)
}

// container is the IoC container which keeps all of the bindings
var container = map[string]binding{}

// bind will bind an abstraction to a concrete, it also set instances for singleton bindings
func bind(resolver interface{}, singleton bool) {
	if reflect.TypeOf(resolver).Kind() != reflect.Func {
		panic("the resolver passed to Singleton() or Transient() methods must be a function")
	}

	if reflect.TypeOf(resolver).NumIn() != 0 {
		panic("the resolver function cannot take any argument")
	}

	for i := 0; i < reflect.TypeOf(resolver).NumOut(); i++ {
		var instance interface{}
		if singleton {
			instance = invoke(resolver)
		}

		container[reflect.TypeOf(resolver).Out(i).String()] = binding{
			resolver: resolver,
			instance: instance,
		}
	}
}

// Singleton will bind an abstraction to a concrete for further singleton resolutions.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
func Singleton(resolver interface{}) {
	bind(resolver, true)
}

// Transient will bind an abstraction to a concrete for further transient resolutions.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
func Transient(resolver interface{}) {
	bind(resolver, false)
}

// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
// It takes an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, the Container invokes the receiver function and pass the related implementations.
func Make(receiver interface{}) {
	if reflect.TypeOf(receiver) == nil {
		return
	}

	if reflect.TypeOf(receiver).Kind() == reflect.Ptr {
		key := reflect.TypeOf(receiver).Elem().String()

		if concrete, ok := container[key]; ok {
			instance := concrete.resolve()
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return
		}
	}

	if reflect.TypeOf(receiver).Kind() == reflect.Func {
		argumentsCount := reflect.TypeOf(receiver).NumIn()
		arguments := make([]reflect.Value, argumentsCount)

		for i := 0; i < argumentsCount; i++ {
			abstraction := reflect.TypeOf(receiver).In(i).String()

			var instance interface{}

			if concrete, ok := container[abstraction]; ok {
				instance = concrete.resolve()
			} else {
				instance = nil
			}

			arguments[i] = reflect.ValueOf(instance)
		}

		reflect.ValueOf(receiver).Call(arguments)
	}
}
