package release

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"

	"github.com/mesosphere/daggers/catalog/goreleaser"
	"github.com/mesosphere/daggers/daggers"
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
func ReleaseWithOptions(opts ...daggers.Option[config]) (*goreleaser.Result, error) {
	debug := mg.Debug() || mg.Verbose()

	options, err := daggers.InitConfig(opts...)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		options = opt(options)
	}

	return goreleaser.Run(goreleaser.CommandRelease, debug, options.Env, options.toArgs())
}
