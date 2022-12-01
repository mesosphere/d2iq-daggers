package githubcli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/options"
	"github.com/mesosphere/daggers/daggers"
)

const (
	// url template for downloading github cli from github releases.
	ghURLTemplate = "https://github.com/cli/cli/releases/download/v%s/gh_%s_linux_amd64.tar.gz"

	// standard source path.
	srcDir = "/src"
)

// Run runs the ginkgo run command with given options.
func Run(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config],
) (string, error) {
	cfg, err := daggers.InitConfig(opts...)
	if err != nil {
		return "", err
	}

	container, err := GetContainer(ctx, runtime, &cfg)
	if err != nil {
		return "", err
	}

	container = container.
		WithMountedDirectory(srcDir, runtime.Workdir()).
		WithWorkdir(srcDir).
		WithEnvVariable("CACHE_BUSTER", time.Now().String()). // Workaround for stop caching after this step
		WithExec(cfg.Args)

	output, err := container.Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// GetContainer returns a dagger container instance with github cli as entrypoint.
func GetContainer(
	ctx context.Context, runtime *daggers.Runtime, cfg *config,
) (*dagger.Container, error) {
	var err error

	// Source url for downloading the Github CLI
	srcURL := fmt.Sprintf(ghURLTemplate, cfg.GithubCliVersion, cfg.GithubCliVersion)

	// Destination file to download tar file contains Github CLI
	dstFile := "/tmp/gh_linux_amd64.tar.gz"

	// Extract Directory
	extractDir := "/tmp"

	// Cli path after extracting downloaded tar to extract directory
	cliPath := fmt.Sprintf("/tmp/gh_%s_linux_amd64/bin/gh", cfg.GithubCliVersion)

	var customizers []options.ContainerCustomizer

	customizers = append(
		customizers, options.WithMountedGoCache(ctx, runtime.Workdir()), options.DownloadFile(srcURL, dstFile),
	)

	container := runtime.Client().Container().From(fmt.Sprintf("%s:%s", cfg.GoBaseImage, cfg.GoVersion))

	for _, customizer := range customizers {
		container, err = customizer(container, runtime.Client())
		if err != nil {
			return nil, err
		}
	}

	token := runtime.Client().Host().EnvVariable("GITHUB_TOKEN").Secret()

	container = container.
		WithSecretVariable("GITHUB_TOKEN", token).
		WithExec([]string{"tar", "-xf", dstFile, "-C", extractDir}).
		WithExec([]string{"mv", cliPath, "/usr/local/bin/gh"}).
		WithExec([]string{"rm", "-rf", "/tmp/*"}).
		WithEntrypoint([]string{"/usr/local/bin/gh"})

	for _, extension := range cfg.Extensions {
		container = container.WithExec([]string{"extension", "install", extension})
	}

	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}
