//go:build go1.18

package container

func ResolveT[T any](c Container) (T, error) {
	var defaultvalue T

	if err := c.Resolve(&defaultvalue); err != nil {
		return defaultvalue, err
	}

	return defaultvalue, nil
}

func MustResolveT[T any](c Container) T {
	return must(ResolveT[T](c))
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
