[![GoDoc](https://godoc.org/github.com/golobby/container?status.svg)](https://godoc.org/github.com/golobby/container)
[![Build Status](https://travis-ci.org/golobby/container.svg?branch=master)](https://travis-ci.org/golobby/container)
[![Coverage Status](https://coveralls.io/repos/github/golobby/container/badge.png?branch=master)](https://coveralls.io/github/golobby/container?branch=master)

# Container
An IoC Container for Go projects. It provides simple, fluent and easy-to-use interface to make dependency injection in 
GoLang very easier.

## Documentation

### Supported Versions
It requires Go `v1.11` or newer versions.

### Installation
To install this package run the following command in the root of your project

```bash
go get github.com/golobby/container
```

### Binding
Binding is a process that you introduce the container that which concrete (implementation) is appropriate for each 
abstraction. In the binding process, you also determine how it must be resolved, singleton or transient. 
In singleton binding, the container provides an instance once and it'd return the instance for each request. 
In transient binding, the container always returns a brand new instance for each request.

Singleton binding using Container:

```go
container.Singleton(func() Abstraction {
  return Implementation
})
```

It takes a resolver function that its return type is the abstraction and the function body configures the related 
concrete (implementation) and returns it.

Transient binding is also similar to singleton binding, see the snippet below.

```go
container.Transient(func() Abstraction {
  return Implementation
})
```

Example for a singleton binding:

```go
container.Singleton(func() Database {
  return &MySQL{}
})
```

And an example for transient binding:

```go
container.Transient(func() Shape {
  return &Rectangle{}
})
```

### Resolving

After bindings, you normally need to resolve the dependencies and receive appropriate implementations of the 
abstractions your code needs.

Container resolves the dependencies with the method `make()`.

#### Using References

One way to get the appropriate implementation of the abstraction you need is to declare an instance of the type of 
abstraction and pass its reference to Container this way:

```go
var x Abstraction
container.Make(&x)
// x will be the implementation of Abstraction
```

Example:

```go
var m Mailer
container.Make(&m)
m.Send("info@miladrahimi.com", "Hello Milad!")
```

#### Using Closures

Another way to resolve the dependencies is by passing a function (receiver) that its arguments are the abstractions you 
need their implementations. Container will invoke the function and pass the related implementation for each abstraction.

```go
container.Make(func(a Abstraction) {
  // a will be the implementation of Abstraction
})
```

Example:

```go
container.Make(func(db Database) {
  // db is an instance of MySQL
  db.Query("...")
})
```

You can also resolve multiple abstractions:

```go
container.Make(func(db Database, s Shape) {
  // db is an instance of MySQL
  db.Query("...")
  // s is an instance of Rectangle
  s.Area()
})
```

## License

GoLobby Container is initially created by 
[@miladrahimi](https://github.com/miladrahimi) and [@amirrezaask](https://github.com/amirrezaask),
and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
