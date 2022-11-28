package precommit

import (
	"context"

	"dagger.io/dagger"

	loggerdagger "github.com/mesosphere/daggers/dagger/logger"
	precommitdagger "github.com/mesosphere/daggers/dagger/precommit"
)

// Precommit runs all the precommit checks. Run `mage help:precommit` for information on available options.
func Precommit(ctx context.Context) error {
	return PrecommitWithOptions(ctx)
}

// PrecommitWithOptions runs all the precommit checks with Dagger options.
//
//nolint:revive // Stuttering is fine here to provide a functional options variant of Precommit function above.
func PrecommitWithOptions(ctx context.Context, opts ...precommitdagger.Option) error {
	// There is a known issue in dagger, if exec command is failed, dagger will not return stdout or stderr.
	// So we need to set verbose to true to see the output of the command until the issue is fixed.
	// issue: https://github.com/dagger/dagger/issues/3192.
	logger, err := loggerdagger.NewLogger(true)
	if err != nil {
		return err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logger))
	if err != nil {
		return err
	}
	defer client.Close()

	// Print the command output to stdout when the issue https://github.com/dagger/dagger/issues/3192. is fixed.
	// Currently, we set verbose to true to see the output of the command.
	_, err = precommitdagger.Run(ctx, client, client.Host().Directory("."), opts...)
	if err != nil {
		return err
	}

	return nil
}
