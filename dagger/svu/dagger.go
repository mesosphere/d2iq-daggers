package svu

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"github.com/magefile/mage/mg"

	loggerdagger "github.com/mesosphere/daggers/dagger/logger"
	"github.com/mesosphere/daggers/daggers"
)

// Output is svu command output.
type Output struct {
	// Version
	Version string
	// Version without the prefix
	VersionWithoutPrefix string
}

// Run runs the svu command with the given options.
func Run(
	ctx context.Context, client *dagger.Client, workdir *dagger.Directory, options ...daggers.Option[config],
) (*Output, error) {
	cfg, err := daggers.InitConfig(options...)
	if err != nil {
		return nil, err
	}

	svuFlags := flagsFromConfig(&cfg)

	container := client.Container().
		From(fmt.Sprintf("ghcr.io/caarlos0/svu:%s", cfg.Version)).
		WithMountedDirectory("/src", workdir).
		WithWorkdir("/src")

	container = container.WithExec(append([]string{cfg.Command}, svuFlags...))
	// Run container and get Exit code
	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	version, err := container.Stdout(ctx)
	if err != nil {
		return nil, err
	}

	svuFlags = append(svuFlags, "--strip-prefix")
	container = container.WithExec(append([]string{cfg.Command}, svuFlags...))
	// Run container and get Exit code
	_, err = container.ExitCode(ctx)
	if err != nil {
		return nil, err
	}

	versionWithoutPrefix, err := container.Stdout(ctx)
	if err != nil {
		return nil, err
	}

	return &Output{
		Version:              strings.TrimSpace(version),
		VersionWithoutPrefix: strings.TrimSpace(versionWithoutPrefix),
	}, nil
}

// SVUWithOptions runs svu with specific options.
//
// TODO: Refactor this to make it more generic and reusable. Temporary solution to get svu working.
//
//nolint:revive // Stuttering is fine here to provide a functional options variant of SVU call.
func SVUWithOptions(ctx context.Context, opts ...daggers.Option[config]) error {
	verbose := mg.Verbose() || mg.Debug()

	logger, err := loggerdagger.NewLogger(verbose)
	if err != nil {
		return err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logger))
	if err != nil {
		return err
	}
	defer client.Close()

	output, err := Run(ctx, client, client.Host().Directory("."), opts...)
	if err != nil {
		return err
	}

	fmt.Println(output.Version)

	return nil
}

func flagsFromConfig(cfg *config) []string {
	var flags []string

	if cfg.Pattern != "" {
		flags = append(flags, "--pattern", cfg.Pattern)
	}
	if cfg.Prefix != "" {
		flags = append(flags, "--prefix", cfg.Prefix)
	}
	if cfg.Suffix != "" {
		flags = append(flags, "--suffix", cfg.Suffix)
	}
	if cfg.TagMode != "" {
		flags = append(flags, "--tag-mode", cfg.TagMode)
	}
	if cfg.Metadata {
		flags = append(flags, "--metadata")
	} else {
		flags = append(flags, "--no-metadata")
	}
	if cfg.Prerelease {
		flags = append(flags, "--pre-release")
	} else {
		flags = append(flags, "--no-pre-release")
	}
	if cfg.Build {
		flags = append(flags, "--build")
	} else {
		flags = append(flags, "--no-build")
	}

	return flags
}
