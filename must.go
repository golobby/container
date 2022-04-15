package container

func MustSingleton(c Container, resolver interface{}) {
	if err := c.Singleton(resolver); err != nil {
		panic(err)
	}
}

func MustNamedSingleton(c Container, name string, resolver interface{}) {
	if err := c.NamedSingleton(name, resolver); err != nil {
		panic(err)
	}
}

func MustTransient(c Container, resolver interface{}) {
	if err := c.Transient(resolver); err != nil {
		panic(err)
	}
}

func MustNamedTransient(c Container, name string, resolver interface{}) {
	if err := c.NamedTransient(name, resolver); err != nil {
		panic(err)
	}
}

func MustCall(c Container, receiver interface{}) {
	if err := c.Call(receiver); err != nil {
		panic(err)
	}
}
