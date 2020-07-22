// Package container provides an IoC container for Go projects.
// It provides simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.
package container

import (
	"reflect"
)

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func invoke(function interface{}, functionType reflect.Type) interface{} {
	return reflect.ValueOf(function).Call(arguments(function, functionType))[0].Interface()
}

// binding keeps a binding resolver and instance (for singleton bindings).
type binding struct {
	resolver     interface{} // resolver function
	instance     interface{} // instance stored for singleton bindings
	resolverType reflect.Type
}

// resolve will return the concrete of related abstraction.
func (b binding) resolve() interface{} {
	if b.instance != nil {
		return b.instance
	}
	return invoke(b.resolver, nil)
}

// container is the IoC container that will keep all of the bindings.
var container = map[reflect.Type]binding{}
var containerPointer = map[reflect.Type]binding{}

// bind will map an abstraction to a concrete and set instance if it's a singleton binding.
func bind(resolver interface{}, singleton bool) {
	resolverTypeOf := reflect.TypeOf(resolver)
	if resolverTypeOf.Kind() != reflect.Func {
		panic("the resolver must be a function")
	}

	for i := 0; i < resolverTypeOf.NumOut(); i++ {
		var instance interface{}
		if singleton {
			instance = invoke(resolver, resolverTypeOf)
		}

		if resolverTypeOf.Out(i).Kind() == reflect.Ptr {
			containerPointer[resolverTypeOf.Out(i)] = binding{
				resolver:     resolver,
				instance:     instance,
				resolverType: resolverTypeOf.Out(i),
			}
		} else {
			container[resolverTypeOf.Out(i)] = binding{
				resolver:     resolver,
				instance:     instance,
				resolverType: resolverTypeOf.Out(i),
			}
		}
	}
}

// arguments will return resolved arguments of the given function.
func arguments(function interface{}, functionTypeOf reflect.Type) []reflect.Value {
	if functionTypeOf == nil {
		functionTypeOf = reflect.TypeOf(function)
	}
	argumentsCount := functionTypeOf.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := functionTypeOf.In(i)

		var instance reflect.Value

		if abstraction.Kind() == reflect.Ptr {
			if concrete, ok := containerPointer[abstraction]; ok {
				instance = reflect.ValueOf(concrete.resolve())
			} else {
				if concrete, ok := container[abstraction.Elem()]; ok {
					//https://github.com/a8m/reflect-examples#wrap-a-reflectvalue-with-pointer-t--t
					data := concrete.resolve()
					pt := reflect.PtrTo(reflect.TypeOf(data)) // create a *T type.
					pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
					pv.Elem().Set(reflect.ValueOf(data))              // sets pv to point to underlying value of v.
					instance = pv
				} else {
					panic("no concrete found for the abstraction: " + abstraction.String())
				}
			}
		} else {
			if concrete, ok := container[abstraction]; ok {
				instance = reflect.ValueOf(concrete.resolve())
			} else {
				if concrete, ok := containerPointer[reflect.PtrTo(abstraction)]; ok {
					instance = reflect.ValueOf(concrete.resolve()).Elem()
				} else {
					panic("no concrete found for the abstraction: " + abstraction.String())
				}
			}
		}
		arguments[i] = instance
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
	containerPointer = map[reflect.Type]binding{}
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

		if abstraction.Kind() == reflect.Ptr {
			if concrete, ok := containerPointer[abstraction]; ok {
				instance := concrete.resolve()
				reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
				return
			} else {
				if concrete, ok := container[abstraction.Elem()]; ok {
					instance := concrete.resolve()
					field := reflect.New(reflect.TypeOf(instance))
					field.Elem().Set(reflect.ValueOf(instance))
					reflect.ValueOf(receiver).Elem().Set(field)
					return
				} else {
					panic("no concrete found for the abstraction " + abstraction.String())
				}
			}
		} else {
			if concrete, ok := container[abstraction]; ok {
				instance := concrete.resolve()
				reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
				return
			} else {
				if concrete, ok := containerPointer[reflect.PtrTo(abstraction)]; ok {
					instance := concrete.resolve()
					reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance).Elem())
					return
				} else {
					panic("no concrete found for the abstraction " + abstraction.String())
				}
			}
		}
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments := arguments(receiver, receiverTypeOf)
		reflect.ValueOf(receiver).Call(arguments)
		return
	}

	panic("the receiver must be either a reference or a callback")
}
