package precommit

import (
	"context"

	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
)

const (
	configFileName      = ".pre-commit-Config.yaml"
	cacheDir            = "/pre-commit-cache"
	precommitHomeEnvVar = "PRE_COMMIT_HOME"
)

// Run runs the precommit checks.
func Run(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[Config],
) (string, error) {
	cfg, err := daggers.InitConfig(opts...)
	if err != nil {
		return "", err
	}

	var (
		url         = "https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz"
		dest        = "/usr/local/bin/pre-commit-2.20.0.pyz"
		customizers = cfg.ContainerCustomizers
	)

	// Configure pre-commit to use the cache volume
	cacheVol, err := containers.NewCacheVolumeWithFileHashKeys(
		ctx, runtime.Client, "pre-commit-", runtime.Workdir, configFileName,
	)
	if err != nil {
		return "", err
	}

	customizers = append(
		customizers,
		containers.DownloadFile(url, dest),
		containers.WithMountedCache(cacheVol, cacheDir, precommitHomeEnvVar),
	)

	container, err := containers.CustomizedContainerFromImage(runtime, cfg.BaseImage, true, customizers...)
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
