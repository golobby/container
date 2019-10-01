// Container provides an IoC Container for Go projects.
// It provides simple, fluent and easy-to-use APIs to make dependency injection in GoLang very easier.
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

// container is the IoC container which keeps all of the bindings
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

// Singleton will bind an abstraction to a concrete for further singleton resolutions.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
func Singleton(resolverFunction interface{}) {
	bind(resolverFunction, true, invoke(resolverFunction))
}

// Transient will bind an abstraction to a concrete for further transient resolutions.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
func Transient(resolverFunction interface{}) {
	bind(resolverFunction, false, nil)
}

// Make will resolve the dependency and return the concrete of given abstraction.
// It takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be resolved,
// the Container invokes the receiver function and pass the related concretes.
func Make(receiverFunction interface{}) {
	if reflect.TypeOf(receiverFunction).Kind() == reflect.Ptr {
		key := reflect.TypeOf(receiverFunction).Elem().String()
		if concrete, ok := container[key]; ok {
			reflect.ValueOf(receiverFunction).Elem().Set(reflect.ValueOf(concrete.resolve()))
			return
		} else {
			panic("There is no concrete bound for " + reflect.TypeOf(receiverFunction).String())
		}
	}

	if reflect.TypeOf(receiverFunction).Kind() != reflect.Func {
		panic("the argument (receiver) passed to Make() is not a function")
	}

	argumentsCount := reflect.TypeOf(receiverFunction).NumIn();
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflect.TypeOf(receiverFunction).In(i).String()

		if concrete, ok := container[abstraction]; ok {
			arguments[i] = reflect.ValueOf(concrete.resolve())
		} else {
			panic("There is no concrete bound for " + abstraction)
		}
	}

	reflect.ValueOf(receiverFunction).Call(arguments)
}
