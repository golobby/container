[![GoDoc](https://godoc.org/github.com/golobby/container/v2?status.svg)](https://godoc.org/github.com/golobby/container/v2)
[![Build Status](https://travis-ci.org/golobby/container.svg?branch=master)](https://travis-ci.org/golobby/container)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/container)](https://goreportcard.com/report/github.com/golobby/container)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome) 
[![Coverage Status](https://coveralls.io/repos/github/golobby/container/badge.svg?branch=master)](https://coveralls.io/github/golobby/container?branch=master)

# Container
A lightweight yet powerful IoC dependency injection container for Go projects.
It provides an easy-to-use interface and performance-in-mind dependency injection container to be your ultimate requirement.

## Documentation

### Required Go Versions
It requires Go `v1.11` or newer versions.

### Installation
To install this package, run the following command in the root of your project.

```bash
go get github.com/golobby/container/v2
```

### Introduction
GoLobby Container like any other IoC dependency injection container is used to bind abstractions to their implementations.
Binding is the process of introducing appropriate concretes (implementations) of abstractions to an IoC container.
In this process, you also determine the resolving type, singleton or transient.
In singleton bindings, the container provides an instance once and returns it for all requests.
In transient bindings, the container always returns a brand new instance for each request.
After the binding process, you can ask the IoC container to make the appropriate implementation of the abstraction that your code depends on.
Now your code depends on abstractions, not implementations.

### Quick Start

The following example demonstrates a simple binding and resolving.

```go
// Bind Config (interface) to JsonConfig
err := container.Singleton(func() Config {
    return &JsonConfig{...}
})

var c Config
err := container.Bind(&c)
// `c` will be an instance of JsonConfig
```

### Binding

#### Singleton

Singleton binding using Container:

```go
err := container.Singleton(func() Abstraction {
  return Implementation
})
```

It takes a resolver function whose return type is the abstraction and the function body returns the concrete (implementation).

Example for singleton binding:

```go
err := container.Singleton(func() Database {
  return &MySQL{}
})
```

In the example above, the container makes a MySQL instance once and returns it for all requests.

#### Transient

Transient binding is also similar to singleton binding.

Example for transient binding:

```go
err := container.Transient(func() Shape {
  return &Rectangle{}
})
```

In the example above, the container always returns a brand new Rectangle instance for each request.

### Resolving

Container resolves the dependencies with the method `make()`.

#### Using References

The `Bind()` method takes a reference of the abstraction type and fills it with the appropriate concrete from the container.

```go
var a Abstraction
err := container.Bind(&a)
// `a` will be an implementation of the Abstraction
```

Example of resolving using refrences:

```go
var m Mailer
err := container.Bind(&m)
// `m` will be an implementation of the Mailer interface
m.Send("contact@miladrahimi.com", "Hello Milad!")
```

#### Using Closures

The `Call()` method takes a function (receiver) with arguments of abstractions you need and invokes it with parameters of appropriate concretes from the container.

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

You can also resolve multiple abstractions like tho follwing example:

```go
err := container.Call(func(db Database, s Shape) {
  db.Query("...")
  s.Area()
})
```

#### Using Structs

The `Fill()` method takes a struct (pointer) with fields of abstractions you need and fills the fields.

```go
type Dependencies struct {
    a Abstraction `container:"inject"`
    x string
}

dep := Dependencies{}

err := container.Fill(&dep)
// `dep.a` will be an implementation of the Abstraction
// `dep.x` will be ignored since it has no `container:"inject"` tag
```

Example of resolving using refrences:

```go
type App struct {
    m Mailer   `container:"inject"`
    d Database `container:"inject"`
    x int
}

myApp := App{}

err := container.Fill(&myApp)
// `myApp.m` will be an implementation of the Mailer interface
// `myApp.s` will be an implementation of the Database interface
// `myApp.x` will be ignored since it has no `container:"inject"` tag
```

#### Binding time

You can resolve dependencies at the binding time if you need previous bindings for the new one.

Example of resolving in binding time:

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

Notice: You can only resolve the dependencies in a binding resolver function that exists in the container.

### Standalone instance

In default, the Container keeps your bindings in the global instance.
Sometimes you may want to create a standalone instance for a part of your application.
If so, create a standalone instance like this example:

```go
c := container.New() // returns a container.Container

err := c.Singleton(func() Database {
    return &MySQL{}
})

err := c.Call(func(db Database) {
    db.Query("...")
})
```

The rest stays the same.
The global container is still available.

### Usage Tips

#### Performance
The package Container inevitably uses reflection in binding and resolving processes. 
If performance is a concern, you should use this package more carefully. 
Try to bind and resolve the dependencies out of the processes that are going to run many times 
(for example, on each request), put it where that run only once when you run your applications 
like main and init functions.

## License

GoLobby Container is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
