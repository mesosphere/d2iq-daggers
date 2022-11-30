package precommit

import (
	"context"

	precommitdagger "github.com/mesosphere/daggers/catalog/precommit"
)

// Precommit runs all the precommit checks. Run `mage help:precommit` for information on available options.
func Precommit(ctx context.Context) error {
	return precommitdagger.PrecommitWithOptions(ctx)
}
