// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package daggers

import (
	"dagger.io/dagger"
)

type runtimeConfig struct {
	verbose   bool
	workdirFn func(client *dagger.Client) *dagger.Directory
}

// WithVerbose sets the verbose option for the runtime config.
func WithVerbose(verbose bool) Option[runtimeConfig] {
	return func(rc runtimeConfig) runtimeConfig {
		rc.verbose = verbose
		return rc
	}
}

// WithWorkdirFn sets the workdir function for getting workdir information.
func WithWorkdirFn(workdirFn func(client *dagger.Client) *dagger.Directory) Option[runtimeConfig] {
	return func(rc runtimeConfig) runtimeConfig {
		rc.workdirFn = workdirFn
		return rc
	}
}

// WithWorkdirFromHostPath sets the workdir option from host path for the runtime config.
func WithWorkdirFromHostPath(workdir string, opts ...dagger.HostDirectoryOpts) Option[runtimeConfig] {
	return func(rc runtimeConfig) runtimeConfig {
		rc.workdirFn = func(client *dagger.Client) *dagger.Directory {
			return client.Host().Directory(workdir, opts...)
		}
		return rc
	}
}
