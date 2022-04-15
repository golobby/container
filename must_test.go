package container_test

import (
	"errors"
	"github.com/golobby/container/v3"
	"testing"
)

func TestMustSingleton_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustSingleton(c, func() (Shape, error) {
		return nil, errors.New("error")
	})
	t.Errorf("panic expcted.")
}

func TestMustNamedSingleton_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustNamedSingleton(c, "name", func() (Shape, error) {
		return nil, errors.New("error")
	})
	t.Errorf("panic expcted.")
}

func TestMustTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustTransient(c, func() (Shape, error) {
		return nil, errors.New("error")
	})
	t.Errorf("panic expcted.")
}

func TestMustNamedTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustNamedTransient(c, "name", func() (Shape, error) {
		return nil, errors.New("error")
	})
	t.Errorf("panic expcted.")
}

func TestMustCall_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	defer func() { recover() }()
	container.MustCall(c, func(s Shape) {
		s.GetArea()
	})
	t.Errorf("panic expcted.")
}

func TestMustResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	var s Shape

	defer func() { recover() }()
	container.MustResolve(c, &s)
	t.Errorf("panic expcted.")
}

func TestMustNamedResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	var s Shape

	defer func() { recover() }()
	container.MustNamedResolve(c, &s, "name")
	t.Errorf("panic expcted.")
}

func TestMustFill_It_Should_Panic_On_Error(t *testing.T) {
	c := container.New()

	myApp := struct {
		S Shape `Global:"type"`
	}{}

	defer func() { recover() }()
	container.MustFill(c, &myApp)
	t.Errorf("panic expcted.")
}
