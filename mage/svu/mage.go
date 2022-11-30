package svu

import (
	"context"

	svudagger "github.com/mesosphere/daggers/catalog/svu"
)

// Current runs svu current.
func Current(ctx context.Context) error {
	return svudagger.SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandCurrent))
}

// Next runs svu next.
func Next(ctx context.Context) error {
	return svudagger.SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandNext))
}

// Major runs svu major.
func Major(ctx context.Context) error {
	return svudagger.SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandMajor))
}

// Minor runs svu minor.
func Minor(ctx context.Context) error {
	return svudagger.SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandMinor))
}

// Patch runs svu patch.
func Patch(ctx context.Context) error {
	return svudagger.SVUWithOptions(ctx, svudagger.WithCommand(svudagger.CommandPatch))
}
