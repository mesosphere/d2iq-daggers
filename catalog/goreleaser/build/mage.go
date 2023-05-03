// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"

	"github.com/mesosphere/d2iq-daggers/catalog/goreleaser"
	"github.com/mesosphere/d2iq-daggers/daggers"
)

// Build runs goreleaser build with --rm-dist and --single-target flags.
func Build(_ context.Context) error {
	result, err := BuildWithOptions(WithArgs("--rm-dist", "--single-target"))
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
	result, err := BuildWithOptions(WithArgs("--snapshot", "--rm-dist", "--single-target"))
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
func BuildWithOptions(opts ...daggers.Option[config]) (*goreleaser.Result, error) {
	debug := mg.Debug() || mg.Verbose()

	options, err := daggers.InitConfig(opts...)
	if err != nil {
		return nil, err
	}

	return goreleaser.Run(goreleaser.CommandBuild, debug, options.Env, options.Args)
}
