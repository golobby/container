// Package container is a lightweight yet powerful IoC container for Go projects.
// It provides an easy-to-use interface and performance-in-mind container to be your ultimate requirement.
package container

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// binding holds a resolver and a concrete (if already resolved).
// It is the break for the Container wall!
type binding struct {
	resolver    interface{} // resolver is the function that is responsible for making the concrete.
	concrete    interface{} // concrete is the stored instance for singleton bindings.
	isSingleton bool        // isSingleton is true if the binding is a singleton.
}

// make resolves the binding if needed and returns the resolved concrete.
func (b *binding) make(c Container) (interface{}, error) {
	if b.concrete != nil {
		return b.concrete, nil
	}

	retVal, err := c.invoke(b.resolver)
	if b.isSingleton {
		b.concrete = retVal
	}

	return retVal, err
}

// Container holds the bindings and provides methods to interact with them.
// It is the entry point in the package.
type Container map[reflect.Type]map[string]*binding

// New creates a new concrete of the Container.
func New() Container {
	return make(Container)
}

// bind maps an abstraction to concrete and instantiates if it is a singleton binding.
func (c Container) bind(resolver interface{}, name string, isSingleton bool, isLazy bool) error {
	reflectedResolver := reflect.TypeOf(resolver)
	if reflectedResolver.Kind() != reflect.Func {
		return errors.New("container: the resolver must be a function")
	}

	if reflectedResolver.NumOut() > 0 {
		if _, exist := c[reflectedResolver.Out(0)]; !exist {
			c[reflectedResolver.Out(0)] = make(map[string]*binding)
		}
	}

	if err := c.validateResolverFunction(reflectedResolver); err != nil {
		return err
	}

	var concrete interface{}
	if !isLazy {
		var err error
		concrete, err = c.invoke(resolver)
		if err != nil {
			return err
		}
	}

	if isSingleton {
		c[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, concrete: concrete, isSingleton: isSingleton}
	} else {
		c[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, isSingleton: isSingleton}
	}

	return nil
}

func (c Container) validateResolverFunction(funcType reflect.Type) error {
	retCount := funcType.NumOut()

	if retCount == 0 || retCount > 2 {
		return errors.New("container: resolver function signature is invalid - it must return abstract, or abstract and error")
	}

	resolveType := funcType.Out(0)
	for i := 0; i < funcType.NumIn(); i++ {
		if funcType.In(i) == resolveType {
			return fmt.Errorf("container: resolver function signature is invalid - depends on abstract it returns")
		}
	}

	return nil
}

// invoke calls a function and its returned values.
// It only accepts one value and an optional error.
func (c Container) invoke(function interface{}) (interface{}, error) {
	arguments, err := c.arguments(function)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(function).Call(arguments)
	if len(values) == 2 && values[1].CanInterface() {
		if err, ok := values[1].Interface().(error); ok {
			return values[0].Interface(), err
		}
	}
	return values[0].Interface(), nil
}

// arguments returns the list of resolved arguments for a function.
func (c Container) arguments(function interface{}) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)
		if concrete, exist := c[abstraction][""]; exist {
			instance, err := concrete.make(c)
			if err != nil {
				return nil, err
			}
			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, errors.New("container: no concrete found for: " + abstraction.String())
		}
	}

	return arguments, nil
}

// Reset deletes all the existing bindings and empties the container.
func (c Container) Reset() {
	for k := range c {
		delete(c, k)
	}
}

// Singleton binds an abstraction to concrete in singleton mode.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c Container) Singleton(resolver interface{}) error {
	return c.bind(resolver, "", true, false)
}

// SingletonLazy binds an abstraction to concrete lazily in singleton mode.
// The concrete is resolved only when the abstraction is resolved for the first time.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c Container) SingletonLazy(resolver interface{}) error {
	return c.bind(resolver, "", true, true)
}

// NamedSingleton binds a named abstraction to concrete in singleton mode.
func (c Container) NamedSingleton(name string, resolver interface{}) error {
	return c.bind(resolver, name, true, false)
}

// NamedSingleton binds a named abstraction to concrete lazily in singleton mode.
// The concrete is resolved only when the abstraction is resolved for the first time.
func (c Container) NamedSingletonLazy(name string, resolver interface{}) error {
	return c.bind(resolver, name, true, true)
}

// Transient binds an abstraction to concrete in transient mode.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c Container) Transient(resolver interface{}) error {
	return c.bind(resolver, "", false, false)
}

// TransientLazy binds an abstraction to concrete lazily in transient mode.
// Normally the resolver will be called during registration, but that is skipped in lazy mode.
// It takes a resolver function that returns the concrete, and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have been declared in the Container already.
func (c Container) TransientLazy(resolver interface{}) error {
	return c.bind(resolver, "", false, true)
}

// NamedTransient binds a named abstraction to concrete lazily in transient mode.
func (c Container) NamedTransient(name string, resolver interface{}) error {
	return c.bind(resolver, name, false, false)
}

// NamedTransient binds a named abstraction to concrete in transient mode.
// Normally the resolver will be called during registration, but that is skipped in lazy mode.
func (c Container) NamedTransientLazy(name string, resolver interface{}) error {
	return c.bind(resolver, name, false, true)
}

// Call takes a receiver function with one or more arguments of the abstractions (interfaces).
// It invokes the receiver function and passes the related concretes.
func (c Container) Call(function interface{}) error {
	receiverType := reflect.TypeOf(function)
	if receiverType == nil || receiverType.Kind() != reflect.Func {
		return errors.New("container: invalid function")
	}

	arguments, err := c.arguments(function)
	if err != nil {
		return err
	}

	result := reflect.ValueOf(function).Call(arguments)

	if len(result) == 0 {
		return nil
	} else if len(result) == 1 && result[0].CanInterface() {
		if result[0].IsNil() {
			return nil
		}
		if err, ok := result[0].Interface().(error); ok {
			return err
		}
	}

	return errors.New("container: receiver function signature is invalid")
}

// Resolve takes an abstraction (reference of an interface type) and fills it with the related concrete.
func (c Container) Resolve(abstraction interface{}) error {
	return c.NamedResolve(abstraction, "")
}

// NamedResolve takes abstraction and its name and fills it with the related concrete.
func (c Container) NamedResolve(abstraction interface{}, name string) error {
	receiverType := reflect.TypeOf(abstraction)
	if receiverType == nil {
		return errors.New("container: invalid abstraction")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()

		if concrete, exist := c[elem][name]; exist {
			if instance, err := concrete.make(c); err == nil {
				reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
				return nil
			} else {
				return fmt.Errorf("container: encountered error while making concrete for: %s. Error encountered: %w", elem.String(), err)
			}
		}

		return errors.New("container: no concrete found for: " + elem.String())
	}

	return errors.New("container: invalid abstraction")
}

// Fill takes a struct and resolves the fields with the tag `container:"inject"`
func (c Container) Fill(structure interface{}) error {
	receiverType := reflect.TypeOf(structure)
	if receiverType == nil {
		return errors.New("container: invalid structure")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()
		if elem.Kind() == reflect.Struct {
			s := reflect.ValueOf(structure).Elem()

			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)

				if t, exist := s.Type().Field(i).Tag.Lookup("container"); exist {
					var name string

					if t == "type" {
						name = ""
					} else if t == "name" {
						name = s.Type().Field(i).Name
					} else {
						return fmt.Errorf("container: %v has an invalid struct tag", s.Type().Field(i).Name)
					}

					if concrete, exist := c[f.Type()][name]; exist {
						instance, err := concrete.make(c)
						if err != nil {
							return err
						}

						ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
						ptr.Set(reflect.ValueOf(instance))

						continue
					}

					return fmt.Errorf("container: cannot make %v field", s.Type().Field(i).Name)
				}
			}

			return nil
		}
	}

	return errors.New("container: invalid structure")
}
