// Package container is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package container

import (
	internal "github.com/golobby/container/v3/pkg/container"
)

// New creates a new standalone instance of Container
func New() internal.Container {
	return internal.New()
}

// container is the global repository of bindings
var container = New()

// Singleton will bind an abstraction to a concrete for further singleton resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Singleton(resolver interface{}) error {
	return container.Singleton(resolver)
}

// NamedSingleton binds like the Singleton method but for named bindings.
func NamedSingleton(name string, resolver interface{}) error {
	return container.NamedSingleton(name, resolver)
}

// Transient will bind an abstraction to a concrete for further transient resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Transient(resolver interface{}) error {
	return container.Transient(resolver)
}

// NamedTransient binds like the Transient method but for named bindings.
func NamedTransient(name string, resolver interface{}) error {
	return container.NamedTransient(name, resolver)
}

// Reset will reset the container and remove all the existing bindings.
func Reset() {
	container.Reset()
}

// Call takes a function with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func Call(receiver interface{}) error {
	return container.Call(receiver)
}

// Resolve takes an abstraction (interface reference) and fill it with the related implementation.
func Resolve(abstraction interface{}) error {
	return container.Resolve(abstraction)
}

// NamedResolve resolves like the Resolve method but for named bindings.
func NamedResolve(abstraction interface{}, name string) error {
	return container.NamedResolve(abstraction, name)
}

// Fill takes a struct and fills the fields with the tag `container:"inject"`
func Fill(receiver interface{}) error {
	return container.Fill(receiver)
}
