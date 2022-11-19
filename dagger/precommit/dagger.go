package precommit

import (
	"context"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/common"
	"github.com/mesosphere/daggers/dagger/options"
)

const (
	configFileName      = ".pre-commit-config.yaml"
	cacheDir            = "/pre-commit-cache"
	precommitHomeEnvVar = "PRE_COMMIT_HOME"
)

// Run runs the precommit checks.
func Run(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option) (string, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		cfg = o(cfg)
	}

	// Create a pre-commit container
	container := client.Container().From(cfg.BaseImage)

	var err error

	for _, c := range cfg.ContainerCustomizers {
		container, err = c(container, client)
		if err != nil {
			return "", err
		}
	}

	container, err = options.DownloadFile(
		"https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz",
		"/usr/local/bin/pre-commit-2.20.0.pyz",
	)(container, client)
	if err != nil {
		return "", err
	}

	// Configure pre-commit to use the cache volume
	cacheVol, err := common.NewCacheVolumeWithFileHashKeys(ctx, client, "pre-commit-", workdir, configFileName)
	if err != nil {
		return "", err
	}

	container = container.WithEnvVariable(precommitHomeEnvVar, cacheDir).WithMountedCache(precommitHomeEnvVar, cacheVol)

	container = container.WithMountedDirectory("/src", workdir).WithWorkdir("/src").
		Exec(dagger.ContainerExecOpts{
			Args: []string{
				"python", "/usr/local/bin/pre-commit-2.20.0.pyz",
				"run", "--all-files", "--show-diff-on-failure",
			},
		})

	// Run container and get Exit code
	return container.Stdout().Contents(ctx)
}
