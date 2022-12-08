// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package daggers

import "github.com/caarlos0/env/v6"

// Option represents an option that can be applied to a config.
type Option[T any] func(T) T

// InitConfig initialize new config using env variables and given options.
func InitConfig[T any](modifiers ...Option[T]) (T, error) {
	//nolint:gocritic // replace `*new(T)` with `T(nil)` is not possible
	cfg := *new(T)

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	for _, o := range modifiers {
		cfg = o(cfg)
	}

	return cfg, nil
}
