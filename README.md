[![Go Reference](https://pkg.go.dev/badge/github.com/golobby/container.svg)](https://pkg.go.dev/github.com/golobby/container)
[![CI](https://github.com/golobby/container/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/container/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/container/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/golobby/container/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/container)](https://goreportcard.com/report/github.com/golobby/container)
[![Coverage Status](https://coveralls.io/repos/github/golobby/container/badge.svg)](https://coveralls.io/github/golobby/container?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

# Container
GoLobby Container is a lightweight yet powerful IoC (dependency injection) container for Go projects.
It's built neat, easy-to-use, and performance-in-mind to be your ultimate requirement.

Features:
- Singleton and Transient bindings
- Named dependencies (bindings)
- Resolve by functions, variables, and structs
- Must helpers that convert errors to panics
- Optional lazy loading of bindings
- Global instance for small applications
- 100% Test coverage!

## Documentation
### Required Go Versions
It requires Go `v1.11` or newer versions.

### Installation
To install this package, run the following command in your project directory.

```bash
go get github.com/golobby/container/v3
```

Next, include it in your application:

```go
import "github.com/golobby/container/v3"
```

### Introduction
GoLobby Container is used to bind abstractions to their implementations.
Binding is the process of introducing appropriate concretes (implementations) of abstractions to an IoC container.
In this process, you also determine the resolving type, singleton or transient.
In singleton bindings, the container provides an instance once and returns it for all the requests.
In transient bindings, the container always returns a brand-new instance for each request.
After the binding process, you can ask the IoC container to make the appropriate implementation of the abstraction that your code needs.
Then your code will depend on abstractions, not implementations!

### Quick Start
The following example demonstrates a simple binding and resolving.

```go
// Bind Config interface to JsonConfig struct
err := container.Singleton(func() Config {
    return &JsonConfig{...}
})

var c Config
err := container.Resolve(&c)
// `c` will be the instance of JsonConfig
```

### Typed Binding
#### Singleton
The following snippet expresses singleton binding.

```go
err := container.Singleton(func() Abstraction {
  return Implementation
})

// If you might return an error...

err := container.Singleton(func() (Abstraction, error) {
  return Implementation, nil
})
```

It takes a resolver (function) whose return type is the abstraction and the function body returns the concrete (implementation).

The example below shows a singleton binding.

```go
err := container.Singleton(func() Database {
  return &MySQL{}
})
```

#### Transient
The example below shows a transient binding.

```go
err := container.Transient(func() Shape {
  return &Rectangle{}
})
```

### Named Bindings
You may have different concretes for an abstraction.
In this case, you can use named bindings instead of typed bindings.
Named bindings take the dependency name into account as well.
The rest is similar to typed bindings.
The following examples demonstrate some named bindings.

```go
// Singleton
err := container.NamedSingleton("square", func() Shape {
    return &Rectangle{}
})
err := container.NamedSingleton("rounded", func() Shape {
    return &Circle{}
})

// Transient
err := container.NamedTransient("sql", func() Database {
    return &MySQL{}
})
err := container.NamedTransient("noSql", func() Database {
    return &MongoDB{}
})
```

### Resolver Errors

The process of creating concrete (resolving) might face an error.
In this case, you can return the error as the second return value like the example below.

```go
err := container.Transient(func() (Shape, error) {
  return nil, errors.New("my-app: cannot create a Shape implementation")
})
```

It could be applied to other binding types.

### Resolving
Container resolves the dependencies with the `Resolve()`, `Call()`, and `Fill()` methods.

#### Using References
The `Resolve()` method takes reference of the abstraction type and fills it with the appropriate concrete.

```go
var a Abstraction
err := container.Resolve(&a)
// `a` will be an implementation of the Abstraction
```

Example of resolving using references:

```go
var m Mailer
err := container.Resolve(&m)
// `m` will be an implementation of the Mailer interface
m.Send("contact@miladrahimi.com", "Hello Milad!")
```

Example of named-resolving using references:

```go
var s Shape
err := container.NamedResolve(&s, "rounded")
// `s` will be an implementation of the Shape that named rounded
```

#### Using Closures
The `Call()` method takes a receiver (function) with arguments of abstractions you need.
It calls it with parameters of appropriate concretes.

```go
err := container.Call(func(a Abstraction) {
    // `a` will be an implementation of the Abstraction
})
```

Example of resolving using closures:

```go
err := container.Call(func(db Database) {
  // `db` will be an implementation of the Database interface
  db.Query("...")
})
```

You can also resolve multiple abstractions like the following example:

```go
err := container.Call(func(db Database, s Shape) {
  db.Query("...")
  s.Area()
})
```

You are able to raise an error in your receiver function, as well.

```go
err := container.Call(func(db Database) error {
  return db.Ping()
})
// err could be `db.Ping()` error.
```

Caution: The `Call()` method does not support named bindings.

#### Using Structs
The `Fill()` method takes a struct (pointer) and resolves its fields.

The example below expresses how the `Fill()` method works.

```go
type App struct {
    mailer Mailer   `container:"type"`
    sql    Database `container:"name"`
    noSql  Database `container:"name"`
    other  int
}

myApp := App{}

err := container.Fill(&myApp)

// [Typed Bindings]
// `myApp.mailer` will be an implementation of the Mailer interface

// [Named Bindings]
// `myApp.sql` will be a sql implementation of the Database interface
// `myApp.noSql` will be a noSql implementation of the Database interface

// `myApp.other` will be ignored since it has no `container` tag
```

#### Binding time
You can resolve dependencies at the binding time if you need previous dependencies for the new one.

The following example shows resolving dependencies at binding time.

```go
// Bind Config to JsonConfig
err := container.Singleton(func() Config {
    return &JsonConfig{...}
})

// Bind Database to MySQL
err := container.Singleton(func(c Config) Database {
    // `c` will be an instance of `JsonConfig`
    return &MySQL{
        Username: c.Get("DB_USERNAME"),
        Password: c.Get("DB_PASSWORD"),
    }
})
```

### Standalone Instance
By default, the Container keeps your bindings in the global instance.
Sometimes you may want to create a standalone instance for a part of your application.
If so, create a standalone instance like the example below.

```go
c := container.New()

err := c.Singleton(func() Database {
    return &MySQL{}
})

err := c.Call(func(db Database) {
    db.Query("...")
})
```

The rest stays the same.
The global container is still available.

### Must Helpers

You might believe that the container shouldn't raise any error and/or you prefer panics.
In this case, Must helpers are for you.
Must helpers are global methods that panic instead of returning errors.

```go
c := container.New()
// Global instance:
// c := container.Global

container.MustSingleton(c, func() Shape {
    return &Circle{a: 13}
})

container.MustCall(c, func(s Shape) {
    // ...
})

// Other Must Helpers:
// container.MustSingleton()
// container.MustSingletonLazy()
// container.MustNamedSingleton()
// container.MustNamedSingletonLazy()
// container.MustTransient()
// container.MustTransientLazy()
// container.MustNamedTransient()
// container.MustNamedTransientLazy()
// container.MustCall()
// container.MustResolve()
// container.MustNamedResolve()
// container.MustFill()
```

### Lazy Binding
Both the singleton and transient binding calls have a lazy version.
Lazy versions defer calling the provided resolver function until the first call.
For singleton bindings, It calls the resolver function only once and stores the result.

Lazy binding methods:
* `container.SingletonLazy()`
* `container.NamedSingletonLazy()`
* `container.TransientLazy()`
* `container.NamedTransientLazy()`

### Performance
The package Container inevitably uses reflection for binding and resolving processes. 
If performance is a concern, try to bind and resolve the dependencies where it runs only once, like the main and init functions.

## License

GoLobby Container is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
