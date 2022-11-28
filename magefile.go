//go:build mage

package main

import (
	"context"
	"os"

	"dagger.io/dagger"

	"github.com/magefile/mage/mg"

	"github.com/mesosphere/daggers/dagger/golang"
)

// Test is a collection of test targets.
//
//goland:noinspection GoUnusedExportedType // used by mage
type Test mg.Namespace

func (Test) Go(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir("."))
	if err != nil {
		return err
	}

	args := []string{"test", "-v", "-race", "-coverprofile", "coverage.txt", "-covermode", "atomic", "./..."}

	_, dir, err := golang.RunCommand(
		ctx,
		client,
		client.Host().Directory("."),
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
