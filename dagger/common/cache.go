package common

import (
	"context"
	"fmt"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers/utils"
)

// NewCacheVolumeWithFileHashKeys creates a new cache volume with generated keys based on the file hashes and prefix.
//
// The key is a SHA256 hash of the prefix and the contents of the given files.
func NewCacheVolumeWithFileHashKeys(
	ctx context.Context, client *dagger.Client, workDir *dagger.Directory, prefix string, files ...string,
) (*dagger.CacheVolume, error) {
	if workDir == nil {
		return nil, fmt.Errorf("%w: workDir", ErrMissingRequiredArgument)
	}

	keys := []string{prefix}

	for _, file := range files {
		contents, err := workDir.File(file).Contents(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", file, err)
		}

		keys = append(keys, contents)
	}

	key := utils.SHA256Sum(keys...)

	return client.CacheVolume(key), nil
}
