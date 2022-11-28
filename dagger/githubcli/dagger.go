package githubcli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/common"
	"github.com/mesosphere/daggers/dagger/options"
)

// standard source path.
const srcDir = "/src"

// Run runs the ginkgo run command with given options.
func Run(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option) (string, error) {
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return "", err
	}

	for _, o := range opts {
		cfg = o(cfg)
	}

	container, err := GetContainer(ctx, client, workdir, &cfg)
	if err != nil {
		return "", err
	}

	container = container.
		WithMountedDirectory(srcDir, workdir).
		WithWorkdir(srcDir).
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
	ctx context.Context, client *dagger.Client, workdir *dagger.Directory, cfg *config,
) (*dagger.Container, error) {
	var err error

	var customizers []options.ContainerCustomizer

	customizers = append(customizers, options.WithMountedGoCache(ctx, workdir))

	container := client.Container().From(fmt.Sprintf("%s:%s", cfg.GoBaseImage, cfg.GoVersion))

	for _, customizer := range customizers {
		container, err = customizer(container, client)
		if err != nil {
			return nil, err
		}
	}

	container, err = common.InstallGithubCLI(ctx, container, cfg.GithubCLIConfig)
	if err != nil {
		return nil, err
	}

	token := client.Host().EnvVariable("GITHUB_TOKEN").Secret()

	container = container.
		WithSecretVariable("GITHUB_TOKEN", token).
		WithEntrypoint([]string{"gh"})

	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}
