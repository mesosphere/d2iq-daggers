// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package daggers

import (
	"context"
	"io"
	"os"

	"dagger.io/dagger"
)

var _ io.Closer = new(Runtime)

// Runtime defines the runtime for a dagger.
type Runtime struct {
	client  *dagger.Client
	workdir *dagger.Directory
}

// NewRuntime returns a new runtime with given options.
func NewRuntime(ctx context.Context, opts ...Option[runtimeConfig]) (*Runtime, error) {
	rc := getRuntimeConfig(opts)

	client, err := getDaggerClient(ctx, rc.verbose)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		client:  client,
		workdir: rc.workdirFn(client),
	}, nil
}

// getRuntimeConfig initializes a runtime config with default values and applies given options before returning it.
func getRuntimeConfig(opts []Option[runtimeConfig]) runtimeConfig {
	rc := runtimeConfig{
		verbose:   false,
		workdirFn: func(client *dagger.Client) *dagger.Directory { return client.Host().Directory(".") },
	}

	for _, o := range opts {
		rc = o(rc)
	}

	return rc
}

// getDaggerClient returns a dagger client with given verbose option.
func getDaggerClient(ctx context.Context, verbose bool) (*dagger.Client, error) {
	logger, err := NewLogger(verbose)
	if err != nil {
		return nil, err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logger))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Client returns the dagger client.
func (r *Runtime) Client() *dagger.Client {
	return r.client
}

// Workdir returns the workdir directory.
func (r *Runtime) Workdir() *dagger.Directory {
	return r.workdir
}

// Close closes the dagger client.
func (r *Runtime) Close() error {
	return r.client.Close()
}

// IsCI returns true if the runtime is running in a CI environment. This check rely on CI environment variable set
// by GitHub Actions. If the CI environment variable is not set, it returns false.
//
// https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
func (r *Runtime) IsCI() bool {
	// Using os.Getenv instead client.Host().Env() because this check should be simple and with dagger client
	// we would need context and error handling which is not needed here.
	return os.Getenv("CI") == "true"
}
