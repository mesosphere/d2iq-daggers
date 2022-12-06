//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	"github.com/mesosphere/daggers/catalog/golang"
	"github.com/mesosphere/daggers/daggers"
)

// Test is a collection of test targets.
//
//goland:noinspection GoUnusedExportedType // used by mage
type Test mg.Namespace

func (Test) Go(ctx context.Context) error {
	runtime, err := daggers.NewRuntime(ctx, daggers.WithVerbose(true))

	args := []string{"test", "-v", "-race", "-coverprofile", "coverage.txt", "-covermode", "atomic", "./..."}

	_, dir, err := golang.RunCommand(
		ctx,
		runtime,
		golang.WithArgs(args...),
		golang.WithEnv(map[string]string{"GOPRIVATE": "github.com/mesosphere"}),
	)
	if err != nil {
		return err
	}

	_, err = dir.File("coverage.txt").Export(ctx, ".output/coverage.txt")
	if err != nil {
		return err
	}

	return nil
}
