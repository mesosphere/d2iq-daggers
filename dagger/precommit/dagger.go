package precommit

import (
	"context"

	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
)

const (
	configFileName      = ".pre-commit-config.yaml"
	cacheDir            = "/pre-commit-cache"
	precommitHomeEnvVar = "PRE_COMMIT_HOME"
)

// Run runs the precommit checks.
func Run(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config],
) (string, error) {
	cfg, err := daggers.InitConfig(opts...) /**/
	if err != nil {
		return "", err
	}

	var (
		url         = "https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz"
		dest        = "/usr/local/bin/pre-commit-2.20.0.pyz"
		customizers = cfg.ContainerCustomizers
	)

	customizers = append(customizers, containers.DownloadFile(url, dest))

	container, err := containers.CustomizedContainerFromImage(runtime, cfg.BaseImage, true, customizers...)
	if err != nil {
		return "", err
	}

	// Configure pre-commit to use the cache volume
	cacheVol, err := containers.NewCacheVolumeWithFileHashKeys(
		ctx, runtime.Client, "pre-commit-", runtime.Workdir, configFileName,
	)
	if err != nil {
		return "", err
	}

	container = container.
		WithEnvVariable(precommitHomeEnvVar, cacheDir).
		WithMountedCache(cacheDir, cacheVol).
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
	runtime, err := daggers.NewRuntime(ctx, daggers.WithVerbose(true))
	if err != nil {
		return err
	}
	defer runtime.Client.Close()

	// Print the command output to stdout when the issue https://github.com/dagger/dagger/issues/3192. is fixed.
	// Currently, we set verbose to true to see the output of the command.
	_, err = Run(ctx, runtime, opts...)
	if err != nil {
		return err
	}

	return nil
}
