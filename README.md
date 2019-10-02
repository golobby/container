[![GoDoc](https://godoc.org/github.com/golobby/container?status.svg)](https://godoc.org/github.com/golobby/container)
[![Build Status](https://travis-ci.org/golobby/container.svg?branch=master)](https://travis-ci.org/golobby/container)
[![Coverage State](https://coveralls.io/repos/github/golobby/container/badge.svg?branch=master)](https://coveralls.io/github/golobby/container)


# Container
An IoC Container for Go projects.
It provides simple, fluent and easy-to-use APIs to make dependency injection in GoLang very easier.

## Documentation

### Supported Versions
It requires Go `v1.11` or newer versions.

### Installation
To install this package run following command in the root of your project

```bash
go get github.com/golobby/container
```

### Binding
To bind an abstraction to a concrete for further singleton resolutions:

```go
container.Singleton(func() Abstraction {
  return Implementation
})
```

It invokes the resolver function once and always return the same object each time you call `make()` method.

And to bind an abstraction to a concrete for further transient resolutions:

```go
container.Transient(func() Abstraction {
  return Implementation
})
```

It invokes the resolver function to provide a brand new object each time you call `make()` method.

Take a look at examples below:

Singleton example:

```go
container.Singleton(func() Database {
  return &MySQL{}
})
```

Transient example:

```go
container.Transient(func() Shape {
  return &Rectangle{}
})
```

### Resolving

To make a concrete by its abstraction you can use `Make()` method.

#### Using References

To resolve the dependencies using reference:

```go
var x Abstraction
container.Make(&x)
// x will be Implementation of Abstraction
```

For example:

```go
var s Shape
container.Make(&s)
s.Area()
```

#### Using Closures

To resolve the dependencies using closure:

```go
container.Make(func(a Abstraction) {
  // a will be a concrete of Abstraction
})
```

For example:

```go
container.Make(func(db Database) {
  // db is an instance of MySQL
  db.Query("...")
})
```

You can resolve multiple abstractions:

```go
container.Make(func(db Database, s Shape) {
  // db is an instance of MySQL
  // s is an instance of Rectangle
  db.Query("...")
  s.Area()
})
```

## License

GoLobby Container is initially created by 
[@miladrahimi](https://github.com/miladrahimi) and [@amirrezaask](https://github.com/amirrezaask),
and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
