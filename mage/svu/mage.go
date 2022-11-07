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
//
//nolint:revive // Stuttering is fine here to provide a functional options variant of SVU call.
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

	optsFromEnv, err := optsFromEnvVars()
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

// optsFromEnvVars loads environment variables into options.
func optsFromEnvVars() ([]svudagger.Option, error) {
	var opts []svudagger.Option

	opts = append(stringOptFromEnvVar(svuVersionEnvVar, svudagger.SVUVersion), opts...)
	opts = append(stringOptFromEnvVar(svuPatternEnvVar, svudagger.WithPattern), opts...)
	opts = append(stringOptFromEnvVar(svuPrefixEnvVar, svudagger.WithPrefix), opts...)
	opts = append(stringOptFromEnvVar(svuSuffixEnvVar, svudagger.WithSuffix), opts...)
	opts = append(stringOptFromEnvVar(svuTagModeEnvVar, withTagModeString), opts...)

	boolOpt, err := boolOptFromEnvVar(svuMetadataEnvVar, svudagger.WithMetadata)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpt, opts...)

	boolOpt, err = boolOptFromEnvVar(svuPreReleaseEnvVar, svudagger.WithPreRelease)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpt, opts...)

	boolOpt, err = boolOptFromEnvVar(svuBuildEnvVar, svudagger.WithBuild)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpt, opts...)

	return opts, nil
}

// stringOptFromEnvVar returns a slice containing a single option if the specified environment variable is set.
// This function returns a slice to make it easier to chain calls in optsFromEnvVars.
// Returns a nil slice if the env var is not set.
func stringOptFromEnvVar(envVarName string, optFn func(string) svudagger.Option) []svudagger.Option {
	if envVarValue, ok := os.LookupEnv(envVarName); ok {
		return []svudagger.Option{optFn(envVarValue)}
	}

	return nil
}

// boolOptFromEnvVar returns a slice containing a single option if the specified environment variable is set.
// This function returns a slice to make it easier to chain calls in optsFromEnvVars.
// Returns a nil slice if the env var is not set.
func boolOptFromEnvVar(envVarName string, optFn func(bool) svudagger.Option) ([]svudagger.Option, error) {
	if envVarValue, ok := os.LookupEnv(envVarName); ok {
		boolVal, err := strconv.ParseBool(envVarValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a boolean: %w", envVarName, err)
		}
		return []svudagger.Option{optFn(boolVal)}, nil
	}

	return nil, nil
}

// withTagModeString sets the tag mode to use when searching for tags set via a string.
func withTagModeString(tagMode string) svudagger.Option {
	return svudagger.WithTagMode(svudagger.TagMode(tagMode))
}
