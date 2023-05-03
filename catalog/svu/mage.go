// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package svu

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"

	"github.com/mesosphere/daggers-for-dkp/daggers"
)

// Current runs svu current.
func Current(ctx context.Context) error {
	return SVUWithOptions(ctx, WithCommand(CommandCurrent))
}

// Next runs svu next.
func Next(ctx context.Context) error {
	return SVUWithOptions(ctx, WithCommand(CommandNext))
}

// Major runs svu major.
func Major(ctx context.Context) error {
	return SVUWithOptions(ctx, WithCommand(CommandMajor))
}

// Minor runs svu minor.
func Minor(ctx context.Context) error {
	return SVUWithOptions(ctx, WithCommand(CommandMinor))
}

// Patch runs svu patch.
func Patch(ctx context.Context) error {
	return SVUWithOptions(ctx, WithCommand(CommandPatch))
}

// SVUWithOptions runs svu with specific options.
//
//nolint:revive // Stuttering is fine here to provide a functional options variant of SVU call.
func SVUWithOptions(ctx context.Context, opts ...daggers.Option[config]) error {
	verbose := mg.Verbose() || mg.Debug()

	runtime, err := daggers.NewRuntime(ctx, daggers.WithVerbose(verbose))
	if err != nil {
		return err
	}
	defer runtime.Close()

	output, err := Run(ctx, runtime, opts...)
	if err != nil {
		return err
	}

	fmt.Println(output.Version)

	return nil
}
