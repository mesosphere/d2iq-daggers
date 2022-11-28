package golang

import (
	"context"
	"fmt"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/options"
)

// standard source path.
const srcDir = "/src"

// RunCommand runs a go command with given working directory and options and returns command output and
// working directory.
func RunCommand(
	ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option,
) (string, *dagger.Directory, error) {
	container, err := GetContainer(ctx, client, workdir, opts...)
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
	ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option,
) (*dagger.Container, error) {
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		cfg = o(cfg)
	}

	container := client.
		Container().
		From(fmt.Sprintf("%s:%s", cfg.GoBaseImage, cfg.GoVersion)).
		WithMountedDirectory(srcDir, workdir).
		WithWorkdir(srcDir).
		WithEntrypoint([]string{"go"})

	for k, v := range cfg.Env {
		container = container.WithEnvVariable(k, v)
	}

	container, err = options.WithMountedGoCache(ctx, workdir)(container, client)
	if err != nil {
		return nil, err
	}

	return container.WithExec(cfg.Args), nil
}
