package precommit

import (
	"context"
	"os"

	"dagger.io/dagger"
	"github.com/magefile/mage/mg"

	loggerdagger "github.com/mesosphere/daggers/dagger/logger"
	precommitdagger "github.com/mesosphere/daggers/dagger/precommit"
)

const (
	baseImageEnvVar = "PRECOMMIT_BASE_IMAGE"
)

// Precommit runs all the precommit checks.
// Configurable via the following environment variables:
//
//	PRECOMMIT_BASE_IMAGE - The base image to run pre-commit in.
func Precommit(ctx context.Context) error {
	return PrecommitWithOptions(ctx)
}

// PrecommitWithOptions runs all the precommit checks with Dagger options.
func PrecommitWithOptions(ctx context.Context, opts ...precommitdagger.Option) error {
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
	if err != nil {
		return err
	}
	defer client.Close()

	if baseImage, ok := os.LookupEnv(baseImageEnvVar); ok {
		opts = append([]precommitdagger.Option{precommitdagger.BaseImage(baseImage)}, opts...)
	}

	return precommitdagger.Run(ctx, client, client.Host().Workdir().Read(), opts...)
}
