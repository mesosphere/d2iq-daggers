package precommit

import (
	_ "embed"

	"context"

	"dagger.io/dagger"
	"github.com/mesosphere/daggers/dagger/options"
)

const (
	configFileName      = ".pre-commit-config.yaml"
	cacheDir            = "/pre-commit-cache"
	precommitHomeEnvVar = "PRE_COMMIT_HOME"
)

func Run(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option) error {
	cfg := defaultConfig()
	for _, o := range opts {
		cfg = o(cfg)
	}

	srcDirID, err := workdir.ID(ctx)
	if err != nil {
		return err
	}

	// Create a pre-commit container
	container := client.
		Container().From(cfg.baseImage)

	for _, c := range cfg.containerCustomizers {
		container, err = c(container, client)
		if err != nil {
			return err
		}
	}

	container, err = options.DownloadFile(
		ctx,
		"https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz",
		"/usr/local/bin/pre-commit-2.20.0.pyz",
	)(container, client)
	if err != nil {
		return err
	}

	container, err = options.CacheDirectoryWithKeyFromFileHash(
		ctx, cacheDir, "precommit-hooks-", configFileName,
	)(container, client)
	if err != nil {
		return err
	}

	container = container.WithEnvVariable(precommitHomeEnvVar, cacheDir).
		WithMountedDirectory("/src", srcDirID).WithWorkdir("/src").
		Exec(dagger.ContainerExecOpts{
			Args: []string{
				"python", "/usr/local/bin/pre-commit-2.20.0.pyz",
				"run", "--all-files", "--show-diff-on-failure",
			},
		})

	// Run container and get Exit code
	_, err = container.ExitCode(ctx)
	if err != nil {
		return err
	}

	return nil
}
