package release

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
	goreleaserConfigEnvVar             = "GORELEASER_RELEASE_CONFIG"
	goreleaserAutoSnapshotEnvVar       = "GORELEASER_RELEASE_AUTO_SNAPSHOT"
	goreleaserParallelismEnvVar        = "GORELEASER_RELEASE_PARALLELISM"
	goreleaserRmDistEnvVar             = "GORELEASER_RELEASE_RM_DIST"
	goreleaserReleaserFooterEnvVar     = "GORELEASER_RELEASE_FOOTER"
	goreleaserReleaserFooterTmplEnvVar = "GORELEASER_RELEASE_FOOTER_TMPL"
	goreleaserReleaserHeaderEnvVar     = "GORELEASER_RELEASE_HEADER"
	goreleaserReleaserHeaderTmplEnvVar = "GORELEASER_RELEASE_HEADER_TMPL"
	goreleaserReleaserNotesEnvVar      = "GORELEASER_RELEASE_NOTES"
	goreleaserReleaseNotesTmplEnvVar   = "GORELEASER_RELEASE_NOTES_TMPL"
	goreleaserSkipAnnounceEnvVar       = "GORELEASER_RELEASE_SKIP_ANNOUNCE"
	goreleaserSkipPublishEnvVar        = "GORELEASER_RELEASE_SKIP_PUBLISH"
	goreleaserSkipSbomEnvVar           = "GORELEASER_RELEASE_SKIP_SBOM"
	goreleaserSkipSignEnvVar           = "GORELEASER_RELEASE_SKIP_SIGN"
	goreleaserSkipValidateEnvVar       = "GORELEASER_RELEASE_SKIP_VALIDATE"
	goreleaserSnapshotEnvVar           = "GORELEASER_RELEASE_SNAPSHOT"
	goreleaserTimeoutEnvVar            = "GORELEASER_RELEASE_TIMEOUT"
)

// Release runs goreleaser release with --rm-dist flags.
func Release(_ context.Context) error {
	result, err := ReleaseWithOptions(WithRmDist(true))
	if err != nil {
		return err
	}

	fmt.Printf(
		"Release is successful for project: %s and version: %s\n",
		result.Metadata.ProjectName,
		result.Metadata.Version,
	)

	return nil
}

// ReleaseSnapshot runs goreleaser release with --snapshot, --rm-dist and --skip-publish flags.
//
//nolint:revive // Disable stuttering check.
func ReleaseSnapshot(_ context.Context) error {
	result, err := ReleaseWithOptions(
		WithRmDist(true),
		SkipPublish(true),
		WithSnapshot(true),
	)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Release snapshot is successful for project: %s and version: %s\n",
		result.Metadata.ProjectName,
		result.Metadata.Version,
	)

	return nil
}

// ReleaseWithOptions runs goreleaser release with specific options.
//
//nolint:revive // Disable stuttering check.
func ReleaseWithOptions(opts ...Option) (*cli.Result, error) {
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

	return cli.Run(cli.CommandRelease, debug, options.Env, options.toArgs())
}

// TODO: make this more readable or come-up with more generic solution
//
//nolint:revive // Disable cognitive-complexity check.
func loadOptionsFromEnv() ([]Option, error) {
	var opts []Option

	opts = append(stringOptFromEnvVar(goreleaserConfigEnvVar, WithConfig), opts...)
	opts = append(stringOptFromEnvVar(goreleaserParallelismEnvVar, WithParallelism), opts...)
	opts = append(stringOptFromEnvVar(goreleaserReleaserFooterEnvVar, WithReleaseFooter), opts...)
	opts = append(stringOptFromEnvVar(goreleaserReleaserFooterTmplEnvVar, WithReleaseFooterTmpl), opts...)
	opts = append(stringOptFromEnvVar(goreleaserReleaserHeaderEnvVar, WithReleaseHeader), opts...)
	opts = append(stringOptFromEnvVar(goreleaserReleaserHeaderTmplEnvVar, WithReleaseHeaderTmpl), opts...)
	opts = append(stringOptFromEnvVar(goreleaserReleaserNotesEnvVar, WithReleaseNotes), opts...)
	opts = append(stringOptFromEnvVar(goreleaserReleaseNotesTmplEnvVar, WithReleaseNotesTmpl), opts...)

	autoSnapshotOpts, err := boolOptFromEnvVar(goreleaserAutoSnapshotEnvVar, WithAutoSnapshot)
	if err != nil {
		return nil, err
	}
	opts = append(autoSnapshotOpts, opts...)

	rmDistOpts, err := boolOptFromEnvVar(goreleaserRmDistEnvVar, WithRmDist)
	if err != nil {
		return nil, err
	}
	opts = append(rmDistOpts, opts...)

	skipAnnounceOpts, err := boolOptFromEnvVar(goreleaserSkipAnnounceEnvVar, SkipAnnounce)
	if err != nil {
		return nil, err
	}
	opts = append(skipAnnounceOpts, opts...)

	skipPublishOpts, err := boolOptFromEnvVar(goreleaserSkipPublishEnvVar, SkipPublish)
	if err != nil {
		return nil, err
	}
	opts = append(skipPublishOpts, opts...)

	skipSbomOpts, err := boolOptFromEnvVar(goreleaserSkipSbomEnvVar, SkipSbom)
	if err != nil {
		return nil, err
	}
	opts = append(skipSbomOpts, opts...)

	skipSignOpts, err := boolOptFromEnvVar(goreleaserSkipSignEnvVar, SkipSign)
	if err != nil {
		return nil, err
	}
	opts = append(skipSignOpts, opts...)

	boolOpts, err := boolOptFromEnvVar(goreleaserSkipValidateEnvVar, SkipValidate)
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
