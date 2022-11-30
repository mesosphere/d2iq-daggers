package svu

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"

	svudagger "github.com/mesosphere/daggers/dagger/svu"
	"github.com/mesosphere/daggers/daggers"
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
// TODO: Refactor this to make it more generic and reusable. Temporary solution to get svu working.
//
//nolint:revive // Stuttering is fine here to provide a functional options variant of SVU call.
func SVUWithOptions(ctx context.Context, opts ...daggers.Option[svudagger.Config]) error {
	verbose := mg.Verbose() || mg.Debug()

	runtime, err := daggers.NewRuntime(ctx, daggers.WithVerbose(verbose))
	if err != nil {
		return err
	}
	defer runtime.Client.Close()

	output, err := svudagger.Run(ctx, runtime, opts...)
	if err != nil {
		return err
	}

	fmt.Println(output.Version)

	return nil
}
