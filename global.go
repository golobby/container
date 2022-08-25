package container

// Global is the global concrete of the Container.
var Global = New()

// Singleton calls the same method of the global concrete.
func Singleton(resolver interface{}) error {
	return Global.Singleton(resolver)
}

// SingletonLazy calls the same method of the global concrete.
func SingletonLazy(resolver interface{}) error {
	return Global.SingletonLazy(resolver)
}

// NamedSingleton calls the same method of the global concrete.
func NamedSingleton(name string, resolver interface{}) error {
	return Global.NamedSingleton(name, resolver)
}

// NamedSingletonLazy calls the same method of the global concrete.
func NamedSingletonLazy(name string, resolver interface{}) error {
	return Global.NamedSingletonLazy(name, resolver)
}

// Transient calls the same method of the global concrete.
func Transient(resolver interface{}) error {
	return Global.Transient(resolver)
}

// TransientLazy calls the same method of the global concrete.
func TransientLazy(resolver interface{}) error {
	return Global.TransientLazy(resolver)
}

// NamedTransient calls the same method of the global concrete.
func NamedTransient(name string, resolver interface{}) error {
	return Global.NamedTransient(name, resolver)
}

// NamedTransientLazy calls the same method of the global concrete.
func NamedTransientLazy(name string, resolver interface{}) error {
	return Global.NamedTransientLazy(name, resolver)
}

// Reset calls the same method of the global concrete.
func Reset() {
	Global.Reset()
}

// Call calls the same method of the global concrete.
func Call(receiver interface{}) error {
	return Global.Call(receiver)
}

// Resolve calls the same method of the global concrete.
func Resolve(abstraction interface{}) error {
	return Global.Resolve(abstraction)
}

// NamedResolve calls the same method of the global concrete.
func NamedResolve(abstraction interface{}, name string) error {
	return Global.NamedResolve(abstraction, name)
}

// Fill calls the same method of the global concrete.
func Fill(receiver interface{}) error {
	return Global.Fill(receiver)
}
