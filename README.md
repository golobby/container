# Container
An IoC Container written in Go

## Documentation

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

And to bind an abstraction to a concrete for further transient resolutions:

```go
container.Transient(func() Abstraction {
  return Implementation
})
```

For example:

```go
import "github.com/golobby/container"

container.Singleton(func() Mailer {
  return &Gmail{}
})
```

### Resolving

To make (resolve) a concrete by its abstraction:

```go
container.Make(func(a Abstraction) {
  // a will be an concrete of Abstraction
})
```

For example:

```go
container.Make(func(m Mailer) {
  m.Send("info@miladrahimi.com", "Hello!")
})
```

## License

GoLobby IoC is initially created by [@miladrahimi](https://github.com/miladrahimi) and [@amirrezaask](https://github.com/amirrezaask) and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
