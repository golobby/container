package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golobby/container/v2/pkg/container"
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

var instance = container.New()

func TestContainer_Singleton(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() {})
	assert.NoError(t, err)

	err = instance.Make(func(s1 Shape) {
		s1.SetArea(666)
	})
	assert.NoError(t, err)

	err = instance.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestContainer_Singleton_With_NonFunction_Resolver_It_Should_Fail(t *testing.T) {
	err := instance.Singleton("STRING!")
	assert.EqualError(t, err, "the resolver must be a function")
}

func TestContainer_Singleton_With_Resolvable_Arguments(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), 666)
		return &MySQL{}
	})
	assert.NoError(t, err)
}

func TestContainer_Transient(t *testing.T) {
	err := instance.Transient(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s1 Shape) {
		s1.SetArea(13)
	})
	assert.NoError(t, err)

	err = instance.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestContainer_Make_With_Multiple_Resolving(t *testing.T) {
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

func TestContainer_Make_With_Reference_As_Resolver(t *testing.T) {
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

func TestContainer_Make_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := instance.Make("STRING!")
	assert.EqualError(t, err, "the receiver must be either a reference or a callback")
}

func TestContainer_Make_With_NonReference_Receiver_It_Should_Fail(t *testing.T) {
	var s Shape
	err := instance.Make(s)
	assert.EqualError(t, err, "cannot detect type of the receiver")
}

func TestContainer_Make_With_UnBounded_Reference_It_Should_Fail(t *testing.T) {
	instance.Reset()

	var s Shape
	err := instance.Make(&s)
	assert.EqualError(t, err, "no concrete found for the abstraction: container_test.Shape")
}

func TestContainer_Make_With_Second_UnBounded_Argument(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Make(func(s Shape, d Database) {})
	assert.EqualError(t, err, "no concrete found for the abstraction: container_test.Database")
}

func TestContainer_Fill(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		S Shape    `container:"inject"`
		D Database `container:"inject"`
		X string
	}{}

	err = instance.Fill(&myApp)
	assert.NoError(t, err)

	assert.IsType(t, &Circle{}, myApp.S)
	assert.IsType(t, &MySQL{}, myApp.D)
}

func TestContainer_Fill_With_Invalid_Field(t *testing.T) {
	type App struct {
		S string `container:"inject"`
	}

	myApp := App{}

	err := instance.Fill(&myApp)
	assert.EqualError(t, err, "cannot resolve S field")
}

func TestContainer_Fill_With_Invalid_Struct(t *testing.T) {
	invalidStruct := 0
	err := instance.Fill(&invalidStruct)
	assert.EqualError(t, err, "invalid structure")
}
