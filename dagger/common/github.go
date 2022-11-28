package common

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GithubCLIConfig is the configuration for installing Github CLI.
type GithubCLIConfig struct {
	GithubCliVersion    string   `env:"GH_VERSION,notEmpty" envDefault:"2.20.2"`
	GithubCliExtensions []string `env:"GH_EXTENSIONS" envDefault:""`
}

// InstallGithubCLI is a dagger step that installs Github CLI.
func InstallGithubCLI(
	ctx context.Context, container *dagger.Container, config GithubCLIConfig,
) (*dagger.Container, error) {
	var err error

	version := config.GithubCliVersion

	// Source url for downloading the Github CLI
	srcURL := fmt.Sprintf(
		"https://github.com/cli/cli/releases/download/v%s/gh_%s_linux_amd64.tar.gz", version, version,
	)

	// Destination file to download tar file contains Github CLI
	dstFile := "/tmp/gh_linux_amd64.tar.gz"

	// Extract Directory
	extractDir := "/tmp"

	// Cli path after extracting downloaded tar to extract directory
	cliPath := fmt.Sprintf("/tmp/gh_%s_linux_amd64/bin/gh", version)

	container = container.
		WithExec(
			[]string{"curl", "--location", "--fail", "--silent", "--show-error", "--output", dstFile, srcURL},
		).
		WithExec([]string{"tar", "-xf", dstFile, "-C", extractDir}).
		WithExec([]string{"mv", cliPath, "/usr/local/bin/gh"}).
		WithExec([]string{"rm", "-rf", "/tmp/*"})

	for _, extension := range config.GithubCliExtensions {
		container = container.WithExec([]string{"extension", "install", extension})
	}

	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// SetupGitAuth is a dagger step that sets up git authentication using the given GITHUB_TOKEN.
func SetupGitAuth(ctx context.Context, client *dagger.Client, container *dagger.Container) (*dagger.Container, error) {
	var err error

	token := client.Host().EnvVariable("GITHUB_TOKEN").Secret()

	container = container.
		WithSecretVariable("GITHUB_TOKEN", token).
		WithExec([]string{"gh", "auth", "setup-git"}).
		WithExec([]string{"gh", "auth", "status"})

	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	return container, nil
}
