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

	// CACHE_BUSTER is workaround for stop caching after this step
	container = container.WithEnvVariable("CACHE_BUSTER", time.Now().String()).WithExec(cfg.Args)

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

	var (
		image       = fmt.Sprintf("%s:%s", cfg.GoImageRepo, cfg.GoImageTag)
		installFn   = containers.InstallGithubCli(cfg.GithubCliVersion, cfg.Extensions...)
		customizers = []containers.ContainerCustomizerFn{installFn}
	)

	customizers = append(customizers, cfg.ContainerCustomizers...)

	container, err := containers.CustomizedContainerFromImage(runtime, image, cfg.MountWorkDir, customizers...)
	if err != nil {
		return nil, err
	}

	for k, v := range cfg.Env {
		container = container.WithEnvVariable(k, v)
	}

	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	return container.WithEntrypoint([]string{"gh"}), nil
}
