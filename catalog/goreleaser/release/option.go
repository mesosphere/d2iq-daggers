// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package release

import "github.com/mesosphere/d2iq-daggers/daggers"

type config struct {
	Env map[string]string

	Args []string `env:"GORELEASER_RELEASE_ARGS" envDefault:""  envSeparator:" "`
}

// WithEnv append extra env variables to goreleaser build process.
func WithEnv(envMap map[string]string) daggers.Option[config] {
	return func(config config) config {
		config.Env = envMap
		return config
	}
}

// WithArgs sets args to goreleaser release process.
func WithArgs(args ...string) daggers.Option[config] {
	return func(config config) config {
		config.Args = args
		return config
	}
}
