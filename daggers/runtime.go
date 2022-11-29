package daggers

import (
	"context"

	"dagger.io/dagger"
)

// Runtime defines the runtime for a dagger.
type Runtime struct {
	Client  *dagger.Client
	Workdir *dagger.Directory
}

// NewRuntime returns a new runtime with given options.
func NewRuntime(ctx context.Context, opts ...Option[runtimeConfig]) (*Runtime, error) {
	rc := getRuntimeConfig(opts)

	client, err := getDaggerClient(ctx, rc.verbose)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Client:  client,
		Workdir: client.Host().Directory(rc.workdir, rc.workdirOpts),
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
