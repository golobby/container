package container_test

import (
	"github.com/golobby/container/pkg/container"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Shape interface {
	SetArea(int)
	GetArea() int
}

type Circle struct {
	a int
}

func (c *Circle) SetArea(a int) {
	c.a = a
}

func (c Circle) GetArea() int {
	return c.a
}

type Database interface {
	Connect() bool
}

type MySQL struct{}

func (m MySQL) Connect() bool {
	return true
}

var instance = container.NewContainer()

func TestSingletonItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	instance.Singleton(func() Shape {
		return &Circle{a: area}
	})

	instance.Make(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, area, a)
	})
}

func TestSingletonItShouldMakeSameObjectEachMake(t *testing.T) {
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	area := 6

	instance.Make(func(s1 Shape) {
		s1.SetArea(area)
	})

	instance.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestSingletonWithNonFunctionResolverItShouldPanic(t *testing.T) {
	value := "the resolver must be a function"
	assert.PanicsWithValue(t, value, func() {
		instance.Singleton("STRING!")
	}, "Expected panic")
}

func TestSingletonItShouldResolveResolverArguments(t *testing.T) {
	area := 5
	instance.Singleton(func() Shape {
		return &Circle{a: area}
	})

	instance.Singleton(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), area)
		return &MySQL{}
	})
}

func TestTransientItShouldMakeDifferentObjectsOnMake(t *testing.T) {
	area := 5

	instance.Transient(func() Shape {
		return &Circle{a: area}
	})

	instance.Make(func(s1 Shape) {
		s1.SetArea(6)
	})

	instance.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestTransientItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	instance.Transient(func() Shape {
		return &Circle{a: area}
	})

	instance.Make(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestMakeWithSingleInputAndCallback(t *testing.T) {
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	instance.Make(func(s Shape) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}
	})
}

func TestMakeWithMultipleInputsAndCallback(t *testing.T) {
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	instance.Singleton(func() Database {
		return &MySQL{}
	})

	instance.Make(func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
}

func TestMakeWithSingleInputAndReference(t *testing.T) {
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	var s Shape

	instance.Make(&s)

	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}
}

func TestMakeWithMultipleInputsAndReference(t *testing.T) {
	instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	instance.Singleton(func() Database {
		return &MySQL{}
	})

	var (
		s Shape
		d Database
	)

	instance.Make(&s)
	instance.Make(&d)

	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestMakeWithUnsupportedReceiver(t *testing.T) {
	value := "the receiver must be either a reference or a callback"
	assert.PanicsWithValue(t, value, func() {
		instance.Make("STRING!")
	}, "Expected panic")
}

func TestMakeWithNonReference(t *testing.T) {
	value := "cannot detect type of the receiver, make sure your are passing reference of the object"
	assert.PanicsWithValue(t, value, func() {
		var s Shape
		instance.Make(s)
	}, "Expected panic")
}

func TestMakeWithUnboundedAbstraction(t *testing.T) {
	value := "no concrete found for the abstraction container_test.Shape"
	assert.PanicsWithValue(t, value, func() {
		var s Shape
		instance.Reset()
		instance.Make(&s)
	}, "Expected panic")
}

func TestMakeWithCallbackThatHasAUnboundedAbstraction(t *testing.T) {
	value := "no concrete found for the abstraction: container_test.Database"
	assert.PanicsWithValue(t, value, func() {
		instance.Reset()
		instance.Singleton(func() Shape {
			return &Circle{}
		})
		instance.Make(func(s Shape, d Database) {})
	}, "Expected panic")
}
