package svu

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"dagger.io/dagger"
	"github.com/magefile/mage/mg"

	loggerdagger "github.com/mesosphere/daggers/dagger/logger"
	svudagger "github.com/mesosphere/daggers/dagger/svu"
)

const (
	svuVersionEnvVar    = "SVU_VERSION"
	svuMetadataEnvVar   = "SVU_METADATA"
	svuPatternEnvVar    = "SVU_PATTERN"
	svuPreReleaseEnvVar = "SVU_PRERELEASE"
	svuBuildEnvVar      = "SVU_BUILD"
	svuPrefixEnvVar     = "SVU_PREFIX"
	svuSuffixEnvVar     = "SVU_SUFFIX"
	svuTagModeEnvVar    = "SVU_TAG_MODE"
)

// Current runs svu current.
func Current(ctx context.Context) error {
	return SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandCurrent))
}

// Next runs svu next.
func Next(ctx context.Context) error {
	return SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandNext))
}

// Major runs svu major.
func Major(ctx context.Context) error {
	return SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandMajor))
}

// Minor runs svu minor.
func Minor(ctx context.Context) error {
	return SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandMinor))
}

// Patch runs svu patch.
func Patch(ctx context.Context) error {
	return SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandPatch))
}

// SVUWithOptions runs svu with specific options.
func SVUWithOptions(ctx context.Context, opts ...svudagger.Option) error {
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

	optsFromEnv, err := loadEnvVars()
	if err != nil {
		return err
	}

	// Combine options from environment variables and options passed to this function. Environment variables
	// take precedence to allow overriding from the arguments passed to this function.
	var combinedOpts []svudagger.Option

	combinedOpts = append(combinedOpts, optsFromEnv...)
	combinedOpts = append(combinedOpts, opts...)

	output, err := svudagger.Run(ctx, client, client.Host().Workdir(), combinedOpts...)
	if err != nil {
		return err
	}

	fmt.Println(output.Version)

	return nil
}

// loads environment variables into options.
func loadEnvVars() ([]svudagger.Option, error) {
	var opts []svudagger.Option

	if svuVersion, ok := os.LookupEnv(svuVersionEnvVar); ok {
		opts = append([]svudagger.Option{svudagger.SVUVersion(svuVersion)}, opts...)
	}

	if svuMetadata, ok := os.LookupEnv(svuMetadataEnvVar); ok {
		metadata, err := strconv.ParseBool(svuMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a boolean: %w", svuMetadataEnvVar, err)
		}
		opts = append([]svudagger.Option{svudagger.WithMetadata(metadata)}, opts...)
	}

	if svuPreRelease, ok := os.LookupEnv(svuPreReleaseEnvVar); ok {
		preRelease, err := strconv.ParseBool(svuPreRelease)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a boolean: %w", svuPreReleaseEnvVar, err)
		}
		opts = append([]svudagger.Option{svudagger.WithPreRelease(preRelease)}, opts...)
	}

	if svuBuild, ok := os.LookupEnv(svuBuildEnvVar); ok {
		build, err := strconv.ParseBool(svuBuild)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a boolean: %w", svuBuildEnvVar, err)
		}
		opts = append([]svudagger.Option{svudagger.WithBuild(build)}, opts...)
	}

	if svuPattern, ok := os.LookupEnv(svuPatternEnvVar); ok {
		opts = append([]svudagger.Option{svudagger.WithPattern(svuPattern)}, opts...)
	}

	if svuPrefix, ok := os.LookupEnv(svuPrefixEnvVar); ok {
		opts = append([]svudagger.Option{svudagger.WithPrefix(svuPrefix)}, opts...)
	}

	if svuSuffix, ok := os.LookupEnv(svuSuffixEnvVar); ok {
		opts = append([]svudagger.Option{svudagger.WithSuffix(svuSuffix)}, opts...)
	}

	if svuTagMode, ok := os.LookupEnv(svuTagModeEnvVar); ok {
		opts = append([]svudagger.Option{svudagger.WithTagMode(svudagger.TagMode(svuTagMode))}, opts...)
	}

	return opts, nil
}
