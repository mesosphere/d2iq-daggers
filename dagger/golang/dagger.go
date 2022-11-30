package golang

import (
	"context"
	"fmt"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
)

// standard source path.
const srcDir = "/src"

// RunCommand runs a go command with given working directory and options and returns command output and
// working directory.
func RunCommand(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config],
) (string, *dagger.Directory, error) {
	container, err := GetContainer(ctx, runtime, opts...)
	if err != nil {
		return "", nil, err
	}

	out, err := container.Stdout(ctx)
	if err != nil {
		return "", nil, err
	}

	return out, container.Directory(srcDir), nil
}

// GetContainer returns a dagger container with given working directory and options.
func GetContainer(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config],
) (*dagger.Container, error) {
	cfg, err := daggers.InitConfig(opts...)
	if err != nil {
		return nil, err
	}

	container := runtime.Client.
		Container().
		From(fmt.Sprintf("%s:%s", cfg.GoBaseImage, cfg.GoVersion)).
		WithMountedDirectory(srcDir, runtime.Workdir).
		WithWorkdir(srcDir).
		WithEntrypoint([]string{"go"})

	for k, v := range cfg.Env {
		container = container.WithEnvVariable(k, v)
	}

	container, err = containers.ApplyCustomizations(runtime, container, containers.WithMountedGoCache(ctx, "."))
	if err != nil {
		return nil, err
	}

	return container.WithExec(cfg.Args), nil
}
