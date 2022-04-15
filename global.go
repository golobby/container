package container

// Global is the global instance of the Container.
var Global = New()

// Singleton calls the same method of the global instance.
func Singleton(resolver interface{}) error {
	return Global.Singleton(resolver)
}

// NamedSingleton calls the same method of the global instance.
func NamedSingleton(name string, resolver interface{}) error {
	return Global.NamedSingleton(name, resolver)
}

// Transient calls the same method of the global instance.
func Transient(resolver interface{}) error {
	return Global.Transient(resolver)
}

// NamedTransient calls the same method of the global instance.
func NamedTransient(name string, resolver interface{}) error {
	return Global.NamedTransient(name, resolver)
}

// Reset calls the same method of the global instance.
func Reset() {
	Global.Reset()
}

// Call calls the same method of the global instance.
func Call(receiver interface{}) error {
	return Global.Call(receiver)
}

// Resolve calls the same method of the global instance.
func Resolve(abstraction interface{}) error {
	return Global.Resolve(abstraction)
}

// NamedResolve calls the same method of the global instance.
func NamedResolve(abstraction interface{}, name string) error {
	return Global.NamedResolve(abstraction, name)
}

// Fill calls the same method of the global instance.
func Fill(receiver interface{}) error {
	return Global.Fill(receiver)
}
