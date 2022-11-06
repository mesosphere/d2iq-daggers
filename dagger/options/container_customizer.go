package options

import (
	"context"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/common"
)

// ContainerCustomizer is a function that customizes a container.
type ContainerCustomizer func(*dagger.Container, *dagger.Client) (*dagger.Container, error)

// AppendToPATH appends the given path to the PATH environment variable.
func AppendToPATH(ctx context.Context, path string) ContainerCustomizer {
	return func(c *dagger.Container, _ *dagger.Client) (*dagger.Container, error) {
		existingPATH, err := c.EnvVariable(ctx, "PATH")
		if err != nil {
			return nil, err
		}
		return c.WithEnvVariable("PATH", existingPATH+":"+path), nil
	}
}

// InstallGo installs Go in the container. Currently it's using hardcoded version 1.19.3 for installation.
func InstallGo(ctx context.Context) ContainerCustomizer {
	return func(c *dagger.Container, client *dagger.Client) (*dagger.Container, error) {
		c = c.Exec(dagger.ContainerExecOpts{
			Args: []string{
				"bash", "-ec",
				`curl --location --fail --silent --show-error https://go.dev/dl/go1.19.3.linux-amd64.tar.gz |
				tar -C /usr/local -xz`,
			},
		})

		c, err := AppendToPATH(ctx, "/usr/local/go/bin")(c, nil)
		if err != nil {
			return nil, err
		}

		c = c.WithEnvVariable("GOCACHE", "/go/build-cache").WithEnvVariable("GOMODCACHE", "/go/mod-cache")

		c, err = CacheDirectoryWithKeyFromFileHash(ctx, "/go/build-cache", "go-build-", "go.sum")(c, client)
		if err != nil {
			return nil, err
		}

		return CacheDirectoryWithKeyFromFileHash(ctx, "/go/mod-cache", "go-mod-", "go.sum")(c, client)
	}
}

// DownloadFile downloads the given URL to the given destination file.
func DownloadFile(url, destFile string) ContainerCustomizer {
	return func(c *dagger.Container, _ *dagger.Client) (*dagger.Container, error) {
		return c.Exec(dagger.ContainerExecOpts{
			Args: []string{
				"curl",
				"--location", "--fail", "--silent", "--show-error",
				"--output", destFile,
				url,
			},
		}), nil
	}
}

// DownloadExecutableFile downloads the given URL to the given destination file and makes it executable.
func DownloadExecutableFile(url, destFile string) ContainerCustomizer {
	return func(c *dagger.Container, _ *dagger.Client) (*dagger.Container, error) {
		c, err := DownloadFile(url, destFile)(c, nil)
		if err != nil {
			return nil, err
		}
		return c.Exec(dagger.ContainerExecOpts{
			Args: []string{
				"chmod", "755", destFile,
			},
		}), nil
	}
}

// CacheDirectoryWithKeyFromFileHash creates a cache volume with a key and hash of the given file and mounts
// it to the given directory.
func CacheDirectoryWithKeyFromFileHash(ctx context.Context, cacheDir, cacheKeyPrefix string, filesToHash ...string) ContainerCustomizer {
	return func(c *dagger.Container, client *dagger.Client) (*dagger.Container, error) {
		cacheVol, err := common.NewCacheVolumeWithFileHashKeys(ctx, client, client.Host().Workdir(), cacheKeyPrefix, filesToHash...)
		if err != nil {
			return nil, err
		}

		return c.WithMountedCache(cacheDir, cacheVol), nil
	}
}
