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

	var (
		image       = fmt.Sprintf("%s:%s", cfg.GoImageRepo, cfg.GoImageTag)
		customizers = cfg.ContainerCustomizers
	)

	if cfg.GoModCacheEnabled {
		customizers = append(customizers, containers.WithMountedGoCache(ctx, cfg.GoModDir))
	}

	container, err := containers.CustomizedContainerFromImage(runtime, image, true, customizers...)
	if err != nil {
		return nil, err
	}

	for k, v := range cfg.Env {
		container = container.WithEnvVariable(k, v)
	}

	return container.WithEntrypoint([]string{"go"}).WithExec(cfg.Args), nil
}
