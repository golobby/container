package container

// MustSingleton wraps the `Singleton` method and panics on errors instead of returning the errors.
func MustSingleton(c Container, resolver interface{}) {
	if err := c.Singleton(resolver); err != nil {
		panic(err)
	}
}

// MustNamedSingleton wraps the `NamedSingleton` method and panics on errors instead of returning the errors.
func MustNamedSingleton(c Container, name string, resolver interface{}) {
	if err := c.NamedSingleton(name, resolver); err != nil {
		panic(err)
	}
}

// MustTransient wraps the `Transient` method and panics on errors instead of returning the errors.
func MustTransient(c Container, resolver interface{}) {
	if err := c.Transient(resolver); err != nil {
		panic(err)
	}
}

// MustNamedTransient wraps the `NamedTransient` method and panics on errors instead of returning the errors.
func MustNamedTransient(c Container, name string, resolver interface{}) {
	if err := c.NamedTransient(name, resolver); err != nil {
		panic(err)
	}
}

// MustCall wraps the `Call` method and panics on errors instead of returning the errors.
func MustCall(c Container, receiver interface{}) {
	if err := c.Call(receiver); err != nil {
		panic(err)
	}
}

// MustResolve wraps the `Resolve` method and panics on errors instead of returning the errors.
func MustResolve(c Container, abstraction interface{}) {
	if err := c.Resolve(abstraction); err != nil {
		panic(err)
	}
}

// MustNamedResolve wraps the `NamedResolve` method and panics on errors instead of returning the errors.
func MustNamedResolve(c Container, abstraction interface{}, name string) {
	if err := c.NamedResolve(abstraction, name); err != nil {
		panic(err)
	}
}

// MustFill wraps the `Fill` method and panics on errors instead of returning the errors.
func MustFill(c Container, receiver interface{}) {
	if err := c.Fill(receiver); err != nil {
		panic(err)
	}
}
