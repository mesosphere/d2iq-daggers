package githubcli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/dagger/options"
)

const (
	// url template for downloading github cli from github releases.
	ghURLTemplate = "https://github.com/cli/cli/releases/download/v%s/gh_%s_linux_amd64.tar.gz"

	// standard source path.
	srcDir = "/src"
)

// Run runs the ginkgo run command with given options.
func Run(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option) (string, error) {
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return "", err
	}

	for _, o := range opts {
		cfg = o(cfg)
	}

	container, err := GetContainer(ctx, client, workdir, &cfg)
	if err != nil {
		return "", err
	}

	container = container.
		WithMountedDirectory(srcDir, workdir).
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
	ctx context.Context, client *dagger.Client, workdir *dagger.Directory, cfg *config,
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

	customizers = append(customizers, options.WithMountedGoCache(ctx, workdir), options.DownloadFile(srcURL, dstFile))

	container := client.Container().From(fmt.Sprintf("%s:%s", cfg.GoBaseImage, cfg.GoVersion))

	for _, customizer := range customizers {
		container, err = customizer(container, client)
		if err != nil {
			return nil, err
		}
	}

	token := client.Host().EnvVariable("GITHUB_TOKEN").Secret()

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
