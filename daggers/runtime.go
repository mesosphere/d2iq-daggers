package daggers

import (
	"context"

	"dagger.io/dagger"

	loggerdagger "github.com/mesosphere/daggers/dagger/logger"
)

// Runtime defines the runtime for a dagger.
type Runtime struct {
	Client  *dagger.Client
	Workdir *dagger.Directory
}

func NewRuntime(ctx context.Context, opts ...RuntimeOption) (*Runtime, error) {
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
func getRuntimeConfig(opts []RuntimeOption) runtimeConfig {
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
	logger, err := loggerdagger.NewLogger(verbose)
	if err != nil {
		return nil, err
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(logger))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ContainerWithAddress returns a new dagger container with given address and load workdir to /src directory
func (r Runtime) ContainerWithAddress(address string) *dagger.Container {
	return r.Client.
		Container().
		From(address).
		WithMountedDirectory("/src", r.Workdir).
		WithWorkdir("/src")
}
