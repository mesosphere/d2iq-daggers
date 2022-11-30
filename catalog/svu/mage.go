package svu

import (
	"context"
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
