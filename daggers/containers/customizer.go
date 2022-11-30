package containers

import (
	"context"
	"fmt"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/common"
	"github.com/mesosphere/daggers/daggers"
)

// ContainerCustomizerFn is a function that customizes a container.
type ContainerCustomizerFn func(*daggers.Runtime, *dagger.Container) (*dagger.Container, error)

// AppendToPATH appends the given path to the PATH environment variable.
func AppendToPATH(ctx context.Context, path string) ContainerCustomizerFn {
	return func(_ *daggers.Runtime, c *dagger.Container) (*dagger.Container, error) {
		existingPATH, err := c.EnvVariable(ctx, "PATH")
		if err != nil {
			return nil, err
		}

		return c.WithEnvVariable("PATH", existingPATH+":"+path), nil
	}
}

// WithMountedGoCache mounts a cache volume for the container's GOCACHE and GOMODCACHE environment variables using
// the contents of the go.mod and go.sum files in the given path. If the path is empty, the current working directory
// is used.
func WithMountedGoCache(ctx context.Context, path string) ContainerCustomizerFn {
	return func(runtime *daggers.Runtime, c *dagger.Container) (*dagger.Container, error) {
		var (
			client     = runtime.Client
			cacheFiles = []string{"go.mod", "go.sum"}
		)

		cacheDir, err := getGoCacheDir(ctx, runtime, path, cacheFiles)
		if err != nil {
			return nil, err
		}

		// Configure go to use the cache volume for the go build cache.
		buildCache, err := common.NewCacheVolumeWithFileHashKeys(ctx, client, "go-build-", cacheDir, cacheFiles...)
		if err != nil {
			return nil, err
		}

		c = c.WithEnvVariable("GOCACHE", "/go/build-cache").WithMountedCache("/go/build-cache", buildCache)

		// Configure go to use the cache volume for the go mod cache.
		modCache, err := common.NewCacheVolumeWithFileHashKeys(ctx, client, "go-mod-", cacheDir, cacheFiles...)
		if err != nil {
			return nil, err
		}

		c = c.WithEnvVariable("GOMODCACHE", "/go/mod-cache").WithMountedCache("/go/mod-cache", modCache)

		return c, nil
	}
}

func getGoCacheDir(
	ctx context.Context, runtime *daggers.Runtime, path string, cacheFiles []string,
) (*dagger.Directory, error) {
	// Default to the current working directory if no path is given.
	if path == "" {
		path = "."
	}

	cacheDir := runtime.Client.Directory()

	for _, cacheFile := range cacheFiles {
		file := runtime.Workdir.Directory(path).File(cacheFile)

		if _, err := file.ID(ctx); err == nil {
			cacheDir = cacheDir.WithFile(path, file)
		}
	}

	// List the files in the cache directory and determine if they exist.
	entries, err := cacheDir.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache files: %w", err)
	}

	if len(entries) != len(cacheFiles) {
		return nil, fmt.Errorf("missing cache files: %v", cacheFiles)
	}

	return cacheDir, nil
}

// InstallGo installs Go in the container using the given version. If the version is empty, the hardcoded "1.19.3" is
// used.
//
// The container must have the "curl" and "tar" binaries installed in order to install Go.
func InstallGo(ctx context.Context, version string) ContainerCustomizerFn {
	return func(runtime *daggers.Runtime, c *dagger.Container) (*dagger.Container, error) {
		// If no version is given, default to 1.19.3.
		if version == "" {
			version = "1.19.3"
		}

		var (
			url = fmt.Sprintf("https://golang.org/dl/go%s.linux-amd64.tar.gz", version)
			cmd = fmt.Sprintf("curl --location --fail --silent --show-error %s | tar -C /usr/local -xz", url)
		)

		c = c.WithExec([]string{"sh", "-ec", cmd})

		return AppendToPATH(ctx, "/usr/local/go/bin")(runtime, c)
	}
}

// DownloadFile downloads the given URL to the given destination file.
func DownloadFile(url, destFile string) ContainerCustomizerFn {
	return func(runtime *daggers.Runtime, c *dagger.Container) (*dagger.Container, error) {
		cmd := fmt.Sprintf("curl --location --fail --silent --show-error %s --output %s", url, destFile)

		return c.WithExec([]string{"sh", "-ec", cmd}), nil
	}
}

// DownloadExecutableFile downloads the given URL to the given destination file and makes it executable.
func DownloadExecutableFile(url, destFile string) ContainerCustomizerFn {
	return func(runtime *daggers.Runtime, c *dagger.Container) (*dagger.Container, error) {
		c, err := DownloadFile(url, destFile)(runtime, c)
		if err != nil {
			return nil, err
		}

		return c.WithExec([]string{"chmod", "755", destFile}), nil
	}
}
