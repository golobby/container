// Package container is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package container

import (
	internal "github.com/golobby/container/v2/pkg/container"
)

// NewContainer creates a new standalone instance of Container
func New() internal.Container {
	return make(internal.Container)
}

// container is the global repository of bindings
var container = internal.NewContainer()

// Singleton will bind an abstraction to a concrete for further singleton resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Singleton(resolver interface{}) error {
	return container.Singleton(resolver)
}

// Transient will bind an abstraction to a concrete for further transient resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Transient(resolver interface{}) error {
	return container.Transient(resolver)
}

// Reset will reset the container and remove all the existing bindings.
func Reset() {
	container.Reset()
}

// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func Make(receiver interface{}) error {
	return container.Make(receiver)
}
