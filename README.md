# Daggers For DKP

`daggers-for-dkp` is a component library for [dagger](https://github.com/dagger/dagger) that provides a collection of tasks and
utilities to make it easier to use.

**WARNING:** The library is still in development and may introduce breaking changes in the future without notice.

## Installation

To install daggers, use the go command:

```bash
$ go get github.com/mesosphere/daggers-for-dkp
```

## Usage

To use daggers, import the package into your project:

```go
import "github.com/mesosphere/daggers-for-dkp/daggers"
```

Then, use the `daggers` package to create a new Runtime:

```go
runtime, err := daggers.NewRuntime(ctx, daggers.WithVerbose(true))
if err != nil {
  panic(err)
}
```

With the runtime, you can then use the catalog to create tasks:

```go
import "github.com/mesosphere/daggers/catalog/golang"
```

```go
_, dir, err := golang.RunCommand(
    ctx,
    runtime,
    golang.WithArgs([]string{"test", "-v", "-race", "-coverprofile", "coverage.txt", "-covermode", "atomic", "./..."}),
)
if err != nil {
    panic(err)
}

_, err = dir.File("coverage.txt").Export(ctx, ".output/coverage.txt")
if err != nil {
    panic(err)
}
```

## License

Apache License 2.0, see [LICENSE](LICENSE).
