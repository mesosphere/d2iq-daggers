package githubcli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
)

// Run runs the ginkgo run command with given options.
func Run(ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config]) (string, error) {
	container, err := GetContainer(ctx, runtime, opts...)
	if err != nil {
		return "", err
	}

	// TODO: this is necessary to get args from the config. We should find a way to do this without any duplication.
	cfg, err := daggers.InitConfig(opts...)
	if err != nil {
		return "", err
	}

	container = containers.MountRuntimeWorkdir(runtime, container).
		WithEnvVariable("CACHE_BUSTER", time.Now().String()). // Workaround for stop caching after this step
		WithExec(cfg.Args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// GetContainer returns a dagger container instance with github cli as entrypoint.
func GetContainer(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config],
) (*dagger.Container, error) {
	var err error

	cfg, err := daggers.InitConfig(opts...)
	if err != nil {
		return nil, err
	}

	container := containers.ContainerFromImage(runtime, fmt.Sprintf("%s:%s", cfg.GoBaseImage, cfg.GoVersion))

	container, err = containers.InstallGithubCli(cfg.GithubCliVersion)(runtime, container)
	if err != nil {
		return nil, err
	}

	container = container.WithEntrypoint([]string{"gh"})

	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}
