package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golobby/container/pkg/container"
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

	err := instance.Singleton(func() Shape {
		return &Circle{a: area}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, area, a)
	})
	assert.NoError(t, err)
}

func TestSingletonItShouldMakeSameObjectEachMake(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	area := 6

	err = instance.Make(func(s1 Shape) {
		s1.SetArea(area)
	})
	assert.NoError(t, err)

	err = instance.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
	assert.NoError(t, err)
}

func TestSingletonWithNonFunctionResolverItShouldPanic(t *testing.T) {
	expectedErr := "the resolver must be a function"
	err := instance.Singleton("STRING!")
	assert.EqualError(t, err, expectedErr)
}

func TestSingletonItShouldResolveResolverArguments(t *testing.T) {
	area := 5
	err := instance.Singleton(func() Shape {
		return &Circle{a: area}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), area)
		return &MySQL{}
	})
	assert.NoError(t, err)
}

func TestTransientItShouldMakeDifferentObjectsOnMake(t *testing.T) {
	area := 5

	err := instance.Transient(func() Shape {
		return &Circle{a: area}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s1 Shape) {
		s1.SetArea(6)
	})
	assert.NoError(t, err)

	err = instance.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
	assert.NoError(t, err)
}

func TestTransientItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	err := instance.Transient(func() Shape {
		return &Circle{a: area}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, a, area)
	})
	assert.NoError(t, err)
}

func TestMakeWithSingleInputAndCallback(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s Shape) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}
	})
	assert.NoError(t, err)
}

func TestMakeWithMultipleInputsAndCallback(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
	assert.NoError(t, err)
}

func TestMakeWithSingleInputAndReference(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	var s Shape

	err = instance.Make(&s)
	assert.NoError(t, err)

	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}
}

func TestMakeWithMultipleInputsAndReference(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	var (
		s Shape
		d Database
	)

	err = instance.Make(&s)
	assert.NoError(t, err)
	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	err = instance.Make(&d)
	assert.NoError(t, err)
	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestMakeWithUnsupportedReceiver(t *testing.T) {
	expectedErr := "the receiver must be either a reference or a callback"
	err := instance.Make("STRING!")
	assert.EqualError(t, err, expectedErr)
}

func TestMakeWithNonReference(t *testing.T) {
	expectedErr := "cannot detect type of the receiver, make sure your are passing reference of the object"
	var s Shape

	err := instance.Make(s)
	assert.EqualError(t, err, expectedErr)
}

func TestMakeWithUnboundedAbstraction(t *testing.T) {
	expectedErr := "no concrete found for the abstraction container_test.Shape"
	var s Shape
	instance.Reset()

	err := instance.Make(&s)
	assert.EqualError(t, err, expectedErr)
}

func TestMakeWithCallbackThatHasAUnboundedAbstraction(t *testing.T) {
	expectedErr := "no concrete found for the abstraction: container_test.Database"
	instance.Reset()

	err := instance.Singleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s Shape, d Database) {})
	assert.EqualError(t, err, expectedErr)
}
