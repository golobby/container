package container

// Global is the global concrete of the Container.
var Global = New()

// Singleton calls the same method of the global concrete.
func Singleton(resolver interface{}) error {
	return Global.Singleton(resolver)
}

// NamedSingleton calls the same method of the global concrete.
func NamedSingleton(name string, resolver interface{}) error {
	return Global.NamedSingleton(name, resolver)
}

// Transient calls the same method of the global concrete.
func Transient(resolver interface{}) error {
	return Global.Transient(resolver)
}

// NamedTransient calls the same method of the global concrete.
func NamedTransient(name string, resolver interface{}) error {
	return Global.NamedTransient(name, resolver)
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
