package options

import (
	"context"

	"dagger.io/dagger"
	"github.com/aweris/tools/utils"
)

type ContainerCustomizer func(*dagger.Container, *dagger.Client) (*dagger.Container, error)

func AppendToPATH(ctx context.Context, dir string) ContainerCustomizer {
	return func(c *dagger.Container, _ *dagger.Client) (*dagger.Container, error) {
		existingPATH, err := c.EnvVariable(ctx, "PATH")
		if err != nil {
			return nil, err
		}
		return c.WithEnvVariable("PATH", existingPATH+":/usr/local/go/bin"), nil
	}
}

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

func DownloadFile(ctx context.Context, url, destFile string) ContainerCustomizer {
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

func DownloadExecutableFile(ctx context.Context, url, destFile string) ContainerCustomizer {
	return func(c *dagger.Container, _ *dagger.Client) (*dagger.Container, error) {
		c, err := DownloadFile(ctx, url, destFile)(c, nil)
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

func CacheDirectory(ctx context.Context, cacheDir, cacheKey string) ContainerCustomizer {
	return func(c *dagger.Container, client *dagger.Client) (*dagger.Container, error) {
		cacheID, err := client.CacheVolume(cacheKey).ID(ctx)
		if err != nil {
			return nil, err
		}
		return c.WithMountedCache(cacheID, cacheDir), nil
	}
}

func CacheDirectoryWithKeyFromFileHash(ctx context.Context, cacheDir, cacheKeyPrefix, fileToHash string) ContainerCustomizer {
	return func(c *dagger.Container, client *dagger.Client) (*dagger.Container, error) {
		fileHash, err := utils.SHA256SumFile(fileToHash)
		if err != nil {
			return nil, err
		}

		return CacheDirectory(ctx, cacheDir, cacheKeyPrefix+fileHash)(c, client)
	}
}
