package container

import "testing"

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

type Mailer interface {
	Send() bool
}

type GMail struct{}

func (g GMail) Send() bool {
	return true
}

func TestSingletonItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	Singleton(func() Shape {
		return &Circle{a: area}
	})

	Make(func(s Shape) {
		if a := s.GetArea(); a != area {
			t.Errorf("Expcted %v got %v", area, a)
		}
	})
}

func TestSingletonItShouldMakeSameObjectOnMake(t *testing.T) {
	Singleton(func() Shape {
		return &Circle{a: 5}
	})

	area := 6

	Make(func(s Shape) {
		s.SetArea(area)
	})

	Make(func(s Shape) {
		if a := s.GetArea(); a != area {
			t.Errorf("Expcted %v got %v", area, a)
		}
	})
}

func TestTransientItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	Transient(func() Shape {
		return &Circle{a: area}
	})

	Make(func(s Shape) {
		if a := s.GetArea(); a != area {
			t.Errorf("Expcted %v got %v", area, a)
		}
	})
}

func TestSingletonItShouldMakeDifferentObjectsOnMake(t *testing.T) {
	area := 5

	Transient(func() Shape {
		return &Circle{a: area}
	})

	Make(func(s Shape) {
		s.SetArea(6)
	})

	Make(func(s Shape) {
		if a := s.GetArea(); a != area {
			t.Errorf("Expcted %v got %v", area, a)
		}
	})
}

func TestMakeWithMultipleInputs(t *testing.T) {
	Singleton(func() Shape {
		return &Circle{a: 5}
	})

	Singleton(func() Mailer {
		return &GMail{}
	})

	Make(func(s Shape, m Mailer) {
		if _, ok := s.(*Circle); !ok {
			t.Errorf("Expcted %v", "Circle")
		}

		if _, ok := m.(*GMail); !ok {
			t.Errorf("Expcted %v", "GMail")
		}
	})
}
