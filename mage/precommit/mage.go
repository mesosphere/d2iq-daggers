package precommit

import (
	"context"
	"os"

	"dagger.io/dagger"

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
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	var opts []precommitdagger.Option

	if baseImage, ok := os.LookupEnv(baseImageEnvVar); ok {
		opts = append(opts, precommitdagger.BaseImage(baseImage))
	}

	return precommitdagger.Run(ctx, client, client.Host().Workdir().Read(), opts...)
}

// PrecommitWithOptions runs all the precommit checks with Dagger options.
func PrecommitWithOptions(ctx context.Context, opts ...precommitdagger.Option) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	if baseImage, ok := os.LookupEnv(baseImageEnvVar); ok {
		opts = append([]precommitdagger.Option{precommitdagger.BaseImage(baseImage)}, opts...)
	}

	return precommitdagger.Run(ctx, client, client.Host().Workdir().Read(), opts...)
}
