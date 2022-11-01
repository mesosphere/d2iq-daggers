package precommit

import (
	_ "embed"

	"context"

	"dagger.io/dagger"
)

func Run(ctx context.Context, client *dagger.Client, workdir *dagger.Directory) error {
	srcDirID, err := workdir.ID(ctx)
	if err != nil {
		return err
	}

	// Create a pre-commit container
	container := client.
		Container().From("python:3.12.0a1-bullseye").
		WithMountedDirectory("/src", srcDirID).WithWorkdir("/src").
		Exec(dagger.ContainerExecOpts{
			Args: []string{
				"curl",
				"--location", "--fail", "--silent", "--show-error",
				"--output", "/usr/local/bin/pre-commit-2.20.0.pyz",
				"https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz",
			},
		}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{
				"python", "/usr/local/bin/pre-commit-2.20.0.pyz",
				"install-hooks",
			},
		}).
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
