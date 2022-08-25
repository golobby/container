package container_test

import (
	"testing"

	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	container.Reset()

	err := container.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestSingletonLazy(t *testing.T) {
	container.Reset()

	err := container.SingletonLazy(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedSingleton(t *testing.T) {
	container.Reset()

	err := container.NamedSingleton("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedSingletonLazy(t *testing.T) {
	container.Reset()

	err := container.NamedSingletonLazy("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestTransient(t *testing.T) {
	container.Reset()

	err := container.Transient(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestTransientLazy(t *testing.T) {
	container.Reset()

	err := container.TransientLazy(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedTransient(t *testing.T) {
	container.Reset()

	err := container.NamedTransient("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedTransientLazy(t *testing.T) {
	container.Reset()

	err := container.NamedTransientLazy("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestCall(t *testing.T) {
	container.Reset()

	err := container.Call(func() {})
	assert.NoError(t, err)
}

func TestResolve(t *testing.T) {
	container.Reset()

	var s Shape

	err := container.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = container.Resolve(&s)
	assert.NoError(t, err)
}

func TestNamedResolve(t *testing.T) {
	container.Reset()

	var s Shape

	err := container.NamedSingleton("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = container.NamedResolve(&s, "rounded")
	assert.NoError(t, err)
}

func TestFill(t *testing.T) {
	container.Reset()

	err := container.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	myApp := struct {
		s Shape `Global:"type"`
	}{}

	err = container.Fill(&myApp)
	assert.NoError(t, err)
}
