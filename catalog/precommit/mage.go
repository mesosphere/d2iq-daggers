// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package precommit

import (
	"context"

	"github.com/mesosphere/daggers-for-dkp/daggers"
)

// Precommit runs all the precommit checks. Run `mage help:precommit` for information on available options.
func Precommit(ctx context.Context) error {
	return PrecommitWithOptions(ctx)
}

// PrecommitWithOptions runs all the precommit checks with Dagger options.
//
//nolint:revive // Stuttering is fine here to provide a functional options variant of Precommit function above.
func PrecommitWithOptions(ctx context.Context, opts ...daggers.Option[config]) error {
	runtime, err := daggers.NewRuntime(ctx, daggers.WithVerbose(true))
	if err != nil {
		return err
	}
	defer runtime.Close()

	// Print the command output to stdout when the issue https://github.com/dagger/dagger/issues/3192. is fixed.
	// Currently, we set verbose to true to see the output of the command.
	_, err = Run(ctx, runtime, opts...)
	if err != nil {
		return err
	}

	return nil
}
