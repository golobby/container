# IoC Container
A IoC Container written in Go

## Documentation

### Installation
To install this package run following command in the root of your project

```bash
go get github.com/golobby/ioc
```

### Binding
To bind an abstraction to a concrete for further singletion resolution:

```go
i := ioc.Container{}
i.Singleton(func() Repository {
  return &UserRepository{}
})
```
And to bind an abstraction to a concrete for further transient resolution:

```go
i := ioc.Container{}
i.Transient(func() Repository {
  return &UserRepository{}
})
```

### Resolving

To make (resolve) an abstraction:

```go
i.Make(func(r Repository) {
  // r will be an instance of UserRepository
})
```

## License

GoLobby IoC is initially created by [@miladrahimi](https://github.com/miladrahimi) and [@amirrezaask](https://github.com/amirrezaask) and released under the [MIT License](http://opensource.org/licenses/mit-license.php).
