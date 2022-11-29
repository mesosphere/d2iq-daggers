package daggers

import (
	"dagger.io/dagger"
)

type runtimeConfig struct {
	verbose     bool
	workdir     string
	workdirOpts dagger.HostDirectoryOpts
}

// WithVerbose sets the verbose option for the runtime config.
func WithVerbose(verbose bool) Option[runtimeConfig] {
	return func(rc runtimeConfig) runtimeConfig {
		rc.verbose = verbose
		return rc
	}
}

// WithWorkdir sets the workdir option for the runtime config.
func WithWorkdir(workdir string) Option[runtimeConfig] {
	return func(rc runtimeConfig) runtimeConfig {
		rc.workdir = workdir
		return rc
	}
}

// WithWorkdirOpts sets the workdir options for the runtime config.
func WithWorkdirOpts(workdirOpts dagger.HostDirectoryOpts) Option[runtimeConfig] {
	return func(rc runtimeConfig) runtimeConfig {
		rc.workdirOpts = workdirOpts
		return rc
	}
}
