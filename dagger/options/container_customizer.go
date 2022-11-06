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

		workDir := client.Host().Workdir()

		// Configure go to use the cache volume for the go build cache.
		buildCache, err := common.NewCacheVolumeWithFileHashKeys(ctx, client, workDir, "go-build", "go.mod", "go.sum")
		if err != nil {
			return nil, err
		}

		c = c.WithEnvVariable("GOCACHE", "/go/build-cache").WithMountedCache("/go/build-cache", buildCache)

		// Configure go to use the cache volume for the go build cache.
		modCache, err := common.NewCacheVolumeWithFileHashKeys(ctx, client, workDir, "go-mod", "go.mod", "go.sum")
		if err != nil {
			return nil, err
		}

		c = c.WithEnvVariable("GOMODCACHE", "/go/mod-cache").WithMountedCache("/go/mod-cache", modCache)

		return c, nil
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
