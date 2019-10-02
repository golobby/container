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

type Database interface {
	Connect() bool
}

type MySQL struct{}

func (m MySQL) Connect() bool {
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

func TestSingletonItShouldMakeSameObjectEachMake(t *testing.T) {
	Singleton(func() Shape {
		return &Circle{a: 5}
	})

	area := 6

	Make(func(s1 Shape) {
		s1.SetArea(area)
	})

	Make(func(s2 Shape) {
		if a := s2.GetArea(); a != area {
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

	Make(func(s1 Shape) {
		s1.SetArea(6)
	})

	Make(func(s2 Shape) {
		if a := s2.GetArea(); a != area {
			t.Errorf("Expcted %v got %v", area, a)
		}
	})
}

func TestMakeWithSingleInputAndCallback(t *testing.T) {
	Singleton(func() Shape {
		return &Circle{a: 5}
	})

	Make(func(s Shape) {
		if _, ok := s.(*Circle); !ok {
			t.Errorf("Expcted %v", "Circle")
		}
	})
}

func TestMakeWithMultipleInputsAndCallback(t *testing.T) {
	Singleton(func() Shape {
		return &Circle{a: 5}
	})

	Singleton(func() Database {
		return &MySQL{}
	})

	Make(func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Errorf("Expcted %v", "Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Errorf("Expcted %v", "MySQL")
		}
	})
}


func TestMakeWithSingleInputAndReference(t *testing.T) {
	Singleton(func() Shape {
		return &Circle{a: 5}
	})

	var s Shape

	Make(&s)

	if _, ok := s.(*Circle); !ok {
		t.Errorf("Expcted %v", "Circle")
	}
}

func TestMakeWithMultipleInputsAndReference(t *testing.T) {
	Singleton(func() Shape {
		return &Circle{a: 5}
	})

	Singleton(func() Database {
		return &MySQL{}
	})

	var (
		s Shape
		d Database
	)

	Make(&s)
	Make(&d)

	if _, ok := s.(*Circle); !ok {
		t.Errorf("Expcted %v", "Circle")
	}

	if _, ok := d.(*MySQL); !ok {
		t.Errorf("Expcted %v", "MySQL")
	}
}
