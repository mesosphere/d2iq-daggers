package svu

import (
	"context"
	"fmt"
	"strings"

	_ "embed"

	"dagger.io/dagger"
)

//go:embed scripts/entrypoint.sh
var entrypointScript string

const (
	baseImage                    = "ghcr.io/caarlos0/svu:v1.9.0"
	versionInfoFile              = ".svu_version"
	versionWithoutPrefixInfoFile = ".svu_version.without_prefix"
)

// Output is svu command output
type Output struct {
	// Version
	Version string
	// Version without the prefix
	VersionWithoutPrefix string
}

func Run(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, options ...Option) (*Output, error) {
	cfg := defaultConfig()
	for _, o := range options {
		cfg = o(cfg)
	}

	// Create a directory for scripts we need to use in container
	scriptsDir := client.
		Directory().
		WithNewFile("entrypoint.sh", dagger.DirectoryWithNewFileOpts{Contents: entrypointScript})

	scriptsDirID, err := scriptsDir.ID(ctx)
	if err != nil {
		return nil, err
	}

	srcDirID, err := workdir.ID(ctx)
	if err != nil {
		return nil, err
	}

	container := client.Container().
		From(baseImage).
		WithMountedDirectory("/scrips", scriptsDirID).
		WithMountedDirectory("/src", srcDirID).
		WithWorkdir("/src").
		WithEnvVariable("SVU_COMMAND", string(cfg.command)).
		WithEnvVariable("SVU_METADATA", fmt.Sprintf("%v", cfg.metadata)).
		WithEnvVariable("SVU_PATTERN", cfg.pattern).
		WithEnvVariable("SVU_PRE_RELEASE", fmt.Sprintf("%v", cfg.preRelease)).
		WithEnvVariable("SVU_BUILD", fmt.Sprintf("%v", cfg.build)).
		WithEnvVariable("SVU_PREFIX", cfg.prefix).
		WithEnvVariable("SVU_SUFFIX", cfg.suffix).
		WithEnvVariable("SVU_TAG_MODE", string(cfg.tagMode)).
		WithEnvVariable("SVU_VERSION_INFO_FILE", versionInfoFile).
		WithEnvVariable("SVU_VERSION_WITHOUT_PREFIX_INFO_FILE", versionWithoutPrefixInfoFile).
		WithEntrypoint([]string{"ash"}).
		Exec(dagger.ContainerExecOpts{Args: []string{"/scrips/entrypoint.sh"}})

	// Run container and get Exit code
	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	version, err := container.Directory("/src").File(versionInfoFile).Contents(ctx)
	if err != nil {
		return nil, err
	}

	versionWithoutPrefix, err := container.Directory("/src").File(versionWithoutPrefixInfoFile).Contents(ctx)
	if err != nil {
		return nil, err
	}

	return &Output{
		Version:              strings.TrimSpace(version),
		VersionWithoutPrefix: strings.TrimSpace(versionWithoutPrefix),
	}, nil
}
