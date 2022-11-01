package precommit

import (
	"context"
	"os"

	"dagger.io/dagger"

	precommitDagger "github.com/aweris/tools/dagger/precommit"
)

// Precommit runs all the precommit checks.
func Precommit(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	return precommitDagger.Run(ctx, client, client.Host().Workdir().Read())
}
