// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package containers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"dagger.io/dagger"
)

// NewCacheVolumeWithFileHashKeys creates a new cache volume with generated keys based on the file hashes and prefix.
func NewCacheVolumeWithFileHashKeys(
	ctx context.Context, client *dagger.Client, cacheKeyPrefix string, workDir *dagger.Directory, fileNames ...string,
) (*dagger.CacheVolume, error) {
	if workDir == nil {
		return nil, fmt.Errorf("%w: workDir", ErrMissingRequiredArgument)
	}

	key, err := cacheKeyFromFiles(ctx, cacheKeyPrefix, workDir, fileNames...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache key from files: %w", err)
	}

	return client.CacheVolume(key), nil
}

// cacheKeyFromFiles returns the string constructed from the supplied prefix suffixed with a SHA256 sum
// of the contents of the requested filenames.
func cacheKeyFromFiles(
	ctx context.Context, cacheKeyPrefix string, workDir *dagger.Directory, fileNames ...string,
) (string, error) {
	h := sha256.New()

	for _, file := range fileNames {
		contents, err := workDir.File(file).Contents(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", file, err)
		}
		_, err = h.Write([]byte(contents))
		if err != nil {
			return "", err
		}
	}

	return cacheKeyPrefix + hex.EncodeToString(h.Sum(nil)), nil
}
