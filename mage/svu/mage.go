package precommit

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"

	"github.com/magefile/mage/mg"

	loggerdagger "github.com/mesosphere/daggers/dagger/logger"
	svudagger "github.com/mesosphere/daggers/dagger/svu"
)

const (
	svuVersionEnvVar = "SVU_VERSION"
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
func SVUWithOptions(ctx context.Context, opts ...svudagger.Option) error {
	verbose := mg.Verbose() || mg.Debug()

	logger, err := loggerdagger.NewLogger(verbose)
	if err != nil {
		return err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logger))
	if err != nil {
		return err
	}
	defer client.Close()

	if svuVersion, ok := os.LookupEnv(svuVersionEnvVar); ok {
		opts = append([]svudagger.Option{svudagger.SVUVersion(svuVersion)}, opts...)
	}

	output, err := svudagger.Run(ctx, client, client.Host().Workdir().Read(), opts...)
	if err != nil {
		return err
	}

	fmt.Println(output.Version)

	return nil
}
