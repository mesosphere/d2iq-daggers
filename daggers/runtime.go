package daggers

import (
	"context"
	"io"

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
		workdir: client.Host().Directory(rc.workdir, rc.workdirOpts),
	}, nil
}

// getRuntimeConfig initializes a runtime config with default values and applies given options before returning it.
func getRuntimeConfig(opts []Option[runtimeConfig]) runtimeConfig {
	rc := runtimeConfig{
		verbose:     false,
		workdir:     ".",
		workdirOpts: dagger.HostDirectoryOpts{},
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
