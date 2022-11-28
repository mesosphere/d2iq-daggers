package common

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GolangImageConfig is the configuration for the golang image.
type GolangImageConfig struct {
	GoImageRepo      string `env:"GO_IMAGE_REPO,notEmpty" envDefault:"docker.io/golang"`
	GoImageTag       string `env:"GO_IMAGE_TAG,notEmpty" envDefault:"1.19"`
	GoModCacheEnable bool   `env:"GO_MOD_CACHE_ENABLE" envDefault:"true"`
	GoModDir         string `env:"GO_MOD_DIR" envDefault:"."`
}

// GetGolangContainer returns a container with the golang image.
func GetGolangContainer(
	ctx context.Context, client *dagger.Client, config GolangImageConfig,
) (*dagger.Container, error) {
	// Create a container with the golang image.
	c := client.Container().From(fmt.Sprintf("%s:%s", config.GoImageRepo, config.GoImageTag))

	// load workdir with go.mod and go.sum only.
	workDir := client.Host().Directory(config.GoModDir, dagger.HostDirectoryOpts{Include: []string{"go.mod", "go.sum"}})

	// List entries in the directory.
	entries, err := workDir.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries from workdir: %w", err)
	}

	// if cache is enabled and go.mod and go.sum are present, mount the directory.
	if config.GoModCacheEnable && len(entries) == 2 {
		// Configure go to use the cache volume for the go build cache.
		buildCache, err := NewCacheVolumeWithFileHashKeys(ctx, client, "go-build-", workDir, "go.mod", "go.sum")
		if err != nil {
			return nil, err
		}

		c = c.WithEnvVariable("GOCACHE", "/go/build-cache").WithMountedCache("/go/build-cache", buildCache)

		// Configure go to use the cache volume for the go build cache.
		modCache, err := NewCacheVolumeWithFileHashKeys(ctx, client, "go-mod-", workDir, "go.mod", "go.sum")
		if err != nil {
			return nil, err
		}

		c = c.WithEnvVariable("GOMODCACHE", "/go/mod-cache").WithMountedCache("/go/mod-cache", modCache)
	} else {
		fmt.Println("go.mod and go.sum not found, skipping go mod cache")
	}

	return c, nil
}
