package precommit

import (
	_ "embed"

	"context"

	"dagger.io/dagger"
)

//go:embed scripts/install_dependencies.sh
var InstallDependenciesScript string

func Run(ctx context.Context, client *dagger.Client, workdir *dagger.Directory) error {
	// Create a directory for scripts we need to use in container
	scriptsDir := client.
		Directory().
		WithNewFile("install_dependencies.sh", dagger.DirectoryWithNewFileOpts{Contents: InstallDependenciesScript})

	scriptsDirID, err := scriptsDir.ID(ctx)
	if err != nil {
		return err
	}

	workdirID, err := workdir.ID(ctx)
	if err != nil {
		return err
	}

	// Create a pre-commit container
	container := client.
		Container().
		From("debian:bullseye-slim").
		WithMountedDirectory("/scripts", scriptsDirID).
		WithMountedDirectory("/work", workdirID).
		WithWorkdir("/work").
		WithEntrypoint([]string{"bash"}).
		Exec(dagger.ContainerExecOpts{Args: []string{"/scripts/install_dependencies.sh"}}).
		Exec(dagger.ContainerExecOpts{Args: []string{"-c", "pre-commit", "run -a --show-diff-on-failure"}})

	// Run container and get Exit code
	_, err = container.ExitCode(ctx)
	if err != nil {
		return err
	}

	return nil
}
