package precommit

import (
	"context"
)

// Precommit runs all the precommit checks. Run `mage help:precommit` for information on available options.
func Precommit(ctx context.Context) error {
	return PrecommitWithOptions(ctx)
}
