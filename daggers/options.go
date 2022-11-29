package daggers

import (
	"dagger.io/dagger"
)

// Option is a function that configures the given generic type.
type Option[T any] func(T) T

// RuntimeOption is a function that configures the runtime config.
type RuntimeOption Option[runtimeConfig]

type runtimeConfig struct {
	verbose     bool
	workdir     string
	workdirOpts dagger.HostDirectoryOpts
}

// WithVerbose sets the verbose option for the runtime config.
func WithVerbose(verbose bool) RuntimeOption {
	return func(rc runtimeConfig) runtimeConfig {
		rc.verbose = verbose
		return rc
	}
}

// WithWorkdir sets the workdir option for the runtime config.
func WithWorkdir(workdir string) RuntimeOption {
	return func(rc runtimeConfig) runtimeConfig {
		rc.workdir = workdir
		return rc
	}
}

// WithWorkdirOpts sets the workdir options for the runtime config.
func WithWorkdirOpts(workdirOpts dagger.HostDirectoryOpts) RuntimeOption {
	return func(rc runtimeConfig) runtimeConfig {
		rc.workdirOpts = workdirOpts
		return rc
	}
}
