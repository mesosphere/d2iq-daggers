package svu

import (
	"context"
	"fmt"
	"strings"

	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
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
	ctx context.Context, runtime *daggers.Runtime, options ...daggers.Option[Config],
) (*Output, error) {
	cfg, err := daggers.InitConfig(options...)
	if err != nil {
		return nil, err
	}

	var (
		image    = fmt.Sprintf("ghcr.io/caarlos0/svu:%s", cfg.Version)
		svuFlags = flagsFromConfig(&cfg)
	)

	container, err := containers.CustomizedContainerFromImage(runtime, image, true)
	if err != nil {
		return nil, err
	}

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

func flagsFromConfig(cfg *Config) []string {
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
