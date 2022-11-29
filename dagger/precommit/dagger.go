package precommit

import (
	"context"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/common"
	"github.com/mesosphere/daggers/dagger/options"
	"github.com/mesosphere/daggers/daggers"
)

const (
	configFileName      = ".pre-commit-config.yaml"
	cacheDir            = "/pre-commit-cache"
	precommitHomeEnvVar = "PRE_COMMIT_HOME"
)

// Run runs the precommit checks.
func Run(
	ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...daggers.Option[config],
) (string, error) {
	cfg, err := daggers.InitConfig(opts...) /**/
	if err != nil {
		return "", err
	}

	// Create a pre-commit container
	container := client.Container().From(cfg.BaseImage)

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
		WithExec(
			[]string{"python", "/usr/local/bin/pre-commit-2.20.0.pyz", "run", "--all-files", "--show-diff-on-failure"},
		)

	// Run container and get Exit code
	return container.Stdout(ctx)
}

// PrecommitWithOptions runs all the precommit checks with Dagger options.
//
// TODO: Refactor this to make it more generic and reusable. Temporary solution to get precommit working.
//
//nolint:revive // Stuttering is fine here to provide a functional options variant of Precommit function above.
func PrecommitWithOptions(ctx context.Context, opts ...daggers.Option[config]) error {
	// There is a known issue in dagger, if exec command is failed, dagger will not return stdout or stderr.
	// So we need to set verbose to true to see the output of the command until the issue is fixed.
	// issue: https://github.com/dagger/dagger/issues/3192.
	logger, err := daggers.NewLogger(true)
	if err != nil {
		return err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logger))
	if err != nil {
		return err
	}
	defer client.Close()

	// Print the command output to stdout when the issue https://github.com/dagger/dagger/issues/3192. is fixed.
	// Currently, we set verbose to true to see the output of the command.
	_, err = Run(ctx, client, client.Host().Directory("."), opts...)
	if err != nil {
		return err
	}

	return nil
}
