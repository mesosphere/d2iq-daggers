package precommit

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/magefile/mage/mg"

	loggerdagger "github.com/mesosphere/daggers/dagger/logger"
	precommitdagger "github.com/mesosphere/daggers/dagger/precommit"
)

const (
	baseImageEnvVar = "PRECOMMIT_BASE_IMAGE"
)

// Precommit runs all the precommit checks. Run `mage help:precommit` for information on available options.
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

	cmdOut, err := precommitdagger.Run(ctx, client, client.Host().Workdir().Read(), opts...)

	// When verbose flag is false, the output is not printed to the console, only redirected to the log file.
	// To work around this, we print the output to the console if the verbose flag is not set.
	if !verbose {
		fmt.Println(cmdOut)
	}

	if err != nil {
		return err
	}

	return nil
}
