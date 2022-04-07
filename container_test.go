package container_test

import (
	"errors"
	"testing"

	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/assert"
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

func TestContainer_With_The_Global_Instance(t *testing.T) {
	err := container.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = container.Call(func(s Shape) {})
	assert.NoError(t, err)

	var sh Shape
	err = container.Resolve(&sh)
	assert.NoError(t, err)

	err = container.Transient(func() Shape {
		return &Circle{a: 33}
	})
	assert.NoError(t, err)

	err = container.Resolve(&sh)
	assert.NoError(t, err)

	err = container.NamedSingleton("theCircle", func() Shape {
		return &Circle{a: 33}
	})
	assert.NoError(t, err)

	err = container.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)

	err = container.NamedTransient("theCircle", func() Shape {
		return &Circle{a: 66}
	})
	assert.NoError(t, err)

	err = container.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)

	myApp := struct {
		S Shape `container:"type"`
	}{}

	err = container.Fill(&myApp)
	assert.NoError(t, err)

	container.Reset()
}

func TestContainer_Singleton(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() {})
	assert.NoError(t, err)

	err = instance.Call(func(s1 Shape) {
		s1.SetArea(666)
	})
	assert.NoError(t, err)

	err = instance.Call(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestContainer_Singleton_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.Singleton(func() (Shape, error) {
		return nil, errors.New("app: error")
	})
	assert.Error(t, err, "app: error")
}

func TestContainer_Singleton_With_NonFunction_Resolver_It_Should_Fail(t *testing.T) {
	err := instance.Singleton("STRING!")
	assert.EqualError(t, err, "container: the resolver must be a function")
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

func TestContainer_Singleton_With_Non_Resolvable_Arguments(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func(s Shape) Shape {
		return &Circle{a: s.GetArea()}
	})
	assert.EqualError(t, err, "container: no concrete found for: container_test.Shape")
}

func TestContainer_NamedSingleton(t *testing.T) {
	err := instance.NamedSingleton("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestContainer_Transient(t *testing.T) {
	err := instance.Transient(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s1 Shape) {
		s1.SetArea(13)
	})
	assert.NoError(t, err)

	err = instance.Call(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestContainer_Transient_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.Transient(func() (Shape, error) {
		return nil, errors.New("app: error")
	})
	assert.Error(t, err, "app: error")

	firstCall := true
	err = instance.Transient(func() (Database, error) {
		if firstCall {
			firstCall = false
			return &MySQL{}, nil
		}
		return nil, errors.New("app: second call error")
	})
	assert.NoError(t, err)

	var db Database
	err = instance.Resolve(&db)
	assert.Error(t, err, "app: second call error")
}

func TestContainer_Transient_With_Resolve_With_Invalid_Signature_It_Should_Fail(t *testing.T) {
	err := instance.Transient(func() (Shape, Database, error) {
		return nil, nil, nil
	})
	assert.Error(t, err, "container: resolver function signature is invalid")
}

func TestContainer_NamedTransient(t *testing.T) {
	err := instance.NamedTransient("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestContainer_Call_With_Multiple_Resolving(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
	assert.NoError(t, err)
}

func TestContainer_Call_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := instance.Call("STRING!")
	assert.EqualError(t, err, "container: invalid function")
}

func TestContainer_Call_With_Second_UnBounded_Argument(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s Shape, d Database) {})
	assert.EqualError(t, err, "container: no concrete found for: container_test.Database")
}

func TestContainer_Resolve_With_Reference_As_Resolver(t *testing.T) {
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

	err = instance.Resolve(&s)
	assert.NoError(t, err)
	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	err = instance.Resolve(&d)
	assert.NoError(t, err)
	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestContainer_Resolve_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := instance.Resolve("STRING!")
	assert.EqualError(t, err, "container: invalid abstraction")
}

func TestContainer_Resolve_With_NonReference_Receiver_It_Should_Fail(t *testing.T) {
	var s Shape
	err := instance.Resolve(s)
	assert.EqualError(t, err, "container: invalid abstraction")
}

func TestContainer_Resolve_With_UnBounded_Reference_It_Should_Fail(t *testing.T) {
	instance.Reset()

	var s Shape
	err := instance.Resolve(&s)
	assert.EqualError(t, err, "container: no concrete found for: container_test.Shape")
}

func TestContainer_Fill_With_Struct_Pointer(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.NamedSingleton("C", func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		S Shape    `container:"type"`
		D Database `container:"type"`
		C Shape    `container:"name"`
		X string
	}{}

	err = instance.Fill(&myApp)
	assert.NoError(t, err)

	assert.IsType(t, &Circle{}, myApp.S)
	assert.IsType(t, &MySQL{}, myApp.D)
}

func TestContainer_Fill_Unexported_With_Struct_Pointer(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		s Shape    `container:"type"`
		d Database `container:"type"`
		y int
	}{}

	err = instance.Fill(&myApp)
	assert.NoError(t, err)

	assert.IsType(t, &Circle{}, myApp.s)
	assert.IsType(t, &MySQL{}, myApp.d)
}

func TestContainer_Fill_With_Invalid_Field_It_Should_Fail(t *testing.T) {
	err := instance.NamedSingleton("C", func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	type App struct {
		S string `container:"name"`
	}

	myApp := App{}

	err = instance.Fill(&myApp)
	assert.EqualError(t, err, "container: cannot resolve S field")
}

func TestContainer_Fill_With_Invalid_Tag_It_Should_Fail(t *testing.T) {
	type App struct {
		S string `container:"invalid"`
	}

	myApp := App{}

	err := instance.Fill(&myApp)
	assert.EqualError(t, err, "container: S has an invalid struct tag")
}

func TestContainer_Fill_With_Invalid_Field_Name_It_Should_Fail(t *testing.T) {
	type App struct {
		S string `container:"name"`
	}

	myApp := App{}

	err := instance.Fill(&myApp)
	assert.EqualError(t, err, "container: cannot resolve S field")
}

func TestContainer_Fill_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	invalidStruct := 0
	err := instance.Fill(&invalidStruct)
	assert.EqualError(t, err, "container: invalid structure")
}

func TestContainer_Fill_With_Invalid_Pointer_It_Should_Fail(t *testing.T) {
	var s Shape
	err := instance.Fill(s)
	assert.EqualError(t, err, "container: invalid structure")
}
