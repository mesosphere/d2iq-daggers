package options

import (
	"context"

	"dagger.io/dagger"
)

type ContainerCustomizer func(*dagger.Container) (*dagger.Container, error)

func AppendToPATH(ctx context.Context, dir string) ContainerCustomizer {
	return func(c *dagger.Container) (*dagger.Container, error) {
		existingPATH, err := c.EnvVariable(ctx, "PATH")
		if err != nil {
			return nil, err
		}
		return c.WithEnvVariable("PATH", existingPATH+":/usr/local/go/bin"), nil
	}
}

func InstallGo(ctx context.Context) ContainerCustomizer {
	return func(c *dagger.Container) (*dagger.Container, error) {
		c = c.Exec(dagger.ContainerExecOpts{
			Args: []string{
				"bash", "-ec",
				`curl --location --fail --silent --show-error https://go.dev/dl/go1.19.3.linux-amd64.tar.gz |
				tar -C /usr/local -xz`,
			},
		})
		return AppendToPATH(ctx, "/usr/local/go/bin")(c)
	}
}

func DownloadFile(ctx context.Context, url, destFile string) ContainerCustomizer {
	return func(c *dagger.Container) (*dagger.Container, error) {
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
	return func(c *dagger.Container) (*dagger.Container, error) {
		c, err := DownloadFile(ctx, url, destFile)(c)
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
