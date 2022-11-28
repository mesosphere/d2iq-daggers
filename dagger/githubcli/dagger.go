package githubcli

import (
	"context"
	"strings"
	"time"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/common"
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

	container, err := GetContainer(ctx, client, &cfg)
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
func GetContainer(ctx context.Context, client *dagger.Client, cfg *config) (*dagger.Container, error) {
	var err error

	container, err := common.GetGolangContainer(ctx, client, cfg.GolangImageConfig)
	if err != nil {
		return nil, err
	}

	container, err = common.InstallGithubCLI(ctx, container, cfg.GithubCLIConfig)
	if err != nil {
		return nil, err
	}

	container, err = common.SetupGitAuth(ctx, client, container)
	if err != nil {
		return nil, err
	}

	return container.WithEntrypoint([]string{"gh"}), nil
}
