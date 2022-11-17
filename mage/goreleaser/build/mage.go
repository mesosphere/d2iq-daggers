package build

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/magefile/mage/mg"

	"github.com/mesosphere/daggers/mage/goreleaser/cli"
)

const (
	goreleaserConfigEnvVar              = "GORELEASER_BUILD_CONFIG"
	goreleaserIDEnvVar                  = "GORELEASER_BUILD_ID"
	goreleaserOutputEnvVar              = "GORELEASER_BUILD_OUTPUT"
	goreleaserParallelismEnvVar         = "GORELEASER_BUILD_PARALLELISM"
	goreleaserRmDistEnvVar              = "GORELEASER_BUILD_RM_DIST"
	goreleaserSingleTargetEnvVar        = "GORELEASER_BUILD_SINGLE_TARGET"
	goreleaserSkipPostCommitHooksEnvVar = "GORELEASER_BUILD_SKIP_POST_COMMIT_HOOKS"
	goreleaserSkipValidateEnvVar        = "GORELEASER_BUILD_SKIP_VALIDATE"
	goreleaserSnapshotEnvVar            = "GORELEASER_BUILD_SNAPSHOT"
	goreleaserTimeoutEnvVar             = "GORELEASER_BUILD_TIMEOUT"
)

// Build runs goreleaser build with --rm-dist and --single-target flags.
func Build(_ context.Context) error {
	result, err := BuildWithOptions(
		WithRmDist(true),
		WithSingleTarget(true),
	)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Build is successful for project: %s and version: %s\n",
		result.Metadata.ProjectName,
		result.Metadata.Version,
	)

	return nil
}

// BuildSnapshot runs goreleaser build with --snapshot, --rm-dist and --single-target flags.
//
//nolint:revive // Disable stuttering check.
func BuildSnapshot(_ context.Context) error {
	result, err := BuildWithOptions(
		WithRmDist(true),
		WithSingleTarget(true),
		WithSnapshot(true),
	)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Build snapshot is successful for project: %s and version: %s\n",
		result.Metadata.ProjectName,
		result.Metadata.Version,
	)

	return nil
}

// BuildWithOptions runs goreleaser build with specific options.
//
//nolint:revive // Disable stuttering check.
func BuildWithOptions(opts ...Option) (*cli.Result, error) {
	debug := mg.Debug() || mg.Verbose()

	envOpts, err := loadOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	// Combine options from environment variables and options passed to this function. Environment variables
	// take precedence to allow overriding from the arguments passed to this function.
	var combinedOpts []Option

	combinedOpts = append(combinedOpts, envOpts...)
	combinedOpts = append(combinedOpts, opts...)

	options := config{}
	for _, opt := range combinedOpts {
		options = opt(options)
	}

	return cli.Run(cli.CommandBuild, debug, options.env, options.toArgs())
}

func loadOptionsFromEnv() ([]Option, error) {
	var opts []Option

	opts = append(stringOptFromEnvVar(goreleaserConfigEnvVar, WithConfig), opts...)
	opts = append(stringOptFromEnvVar(goreleaserIDEnvVar, WithID), opts...)
	opts = append(stringOptFromEnvVar(goreleaserOutputEnvVar, WithOutput), opts...)
	opts = append(stringOptFromEnvVar(goreleaserParallelismEnvVar, WithParallelism), opts...)

	boolOpts, err := boolOptFromEnvVar(goreleaserRmDistEnvVar, WithRmDist)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpts, opts...)

	boolOpts, err = boolOptFromEnvVar(goreleaserSingleTargetEnvVar, WithSingleTarget)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpts, opts...)

	boolOpts, err = boolOptFromEnvVar(goreleaserSkipPostCommitHooksEnvVar, SkipPostCommitHooks)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpts, opts...)

	boolOpts, err = boolOptFromEnvVar(goreleaserSkipValidateEnvVar, SkipValidate)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpts, opts...)

	boolOpts, err = boolOptFromEnvVar(goreleaserSnapshotEnvVar, WithSnapshot)
	if err != nil {
		return nil, err
	}
	opts = append(boolOpts, opts...)

	durationOpts, err := durationOptFromEnvVar(goreleaserTimeoutEnvVar, WithTimeout)
	if err != nil {
		return nil, err
	}
	opts = append(durationOpts, opts...)

	return opts, nil
}

// stringOptFromEnvVar returns a slice containing a single option if the specified environment variable is set.
// This function returns a slice to make it easier to chain calls in optsFromEnvVars.
// Returns a nil slice if the env var is not set.
func stringOptFromEnvVar(envVarName string, optFn func(string) Option) []Option {
	if envVarValue, ok := os.LookupEnv(envVarName); ok {
		return []Option{optFn(envVarValue)}
	}

	return nil
}

// boolOptFromEnvVar returns a slice containing a single option if the specified environment variable is set.
// This function returns a slice to make it easier to chain calls in optsFromEnvVars.
// Returns a nil slice if the env var is not set.
func boolOptFromEnvVar(envVarName string, optFn func(bool) Option) ([]Option, error) {
	if envVarValue, ok := os.LookupEnv(envVarName); ok {
		boolVal, err := strconv.ParseBool(envVarValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a boolean: %w", envVarName, err)
		}
		return []Option{optFn(boolVal)}, nil
	}

	return nil, nil
}

// durationOptFromEnvVar returns a slice containing a single option if the specified environment variable is set.
// This function returns a slice to make it easier to chain calls in optsFromEnvVars.
// Returns a nil slice if the env var is not set.
func durationOptFromEnvVar(envVarName string, optFn func(duration time.Duration) Option) ([]Option, error) {
	if envVarValue, ok := os.LookupEnv(envVarName); ok {
		durationVal, err := time.ParseDuration(envVarValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a duration: %w", envVarName, err)
		}
		return []Option{optFn(durationVal)}, nil
	}

	return nil, nil
}
