package build

import (
	"time"
)

type config struct {
	env                 map[string]string
	config              string
	id                  string
	output              string
	parallelism         string
	rmDist              bool
	singleTarget        bool
	skipPostCommitHooks bool
	skipValidate        bool
	snapshot            bool
	timeout             string
}

// Option is a function that configures the goreleaser build options.
type Option func(config config) config

// WithEnv append extra env variables to goreleaser build process.
func WithEnv(env map[string]string) Option {
	return func(config config) config {
		config.env = env
		return config
	}
}

// WithConfig sets  --config flag.
func WithConfig(configPath string) Option {
	return func(config config) config {
		config.config = configPath
		return config
	}
}

// WithID sets --id flag.
func WithID(id string) Option {
	return func(config config) config {
		config.id = id
		return config
	}
}

// WithOutput sets --output.
func WithOutput(output string) Option {
	return func(config config) config {
		config.output = output
		return config
	}
}

// WithParallelism sets --parallelism.
func WithParallelism(parallelism string) Option {
	return func(config config) config {
		config.parallelism = parallelism
		return config
	}
}

// WithRmDist sets --rm-dist.
func WithRmDist(rmDist bool) Option {
	return func(config config) config {
		config.rmDist = rmDist
		return config
	}
}

// WithSingleTarget sets --single-target.
func WithSingleTarget(singleTarget bool) Option {
	return func(config config) config {
		config.singleTarget = singleTarget
		return config
	}
}

// SkipPostCommitHooks sets--skip-post-hooks.
func SkipPostCommitHooks(skipPostCommitHooks bool) Option {
	return func(config config) config {
		config.skipPostCommitHooks = skipPostCommitHooks
		return config
	}
}

// SkipValidate sets --skip-validate.
func SkipValidate(skipValidate bool) Option {
	return func(config config) config {
		config.skipValidate = skipValidate
		return config
	}
}

// WithSnapshot sets --snapshot.
func WithSnapshot(snapshot bool) Option {
	return func(config config) config {
		config.snapshot = snapshot
		return config
	}
}

// WithTimeout sets --timeout duration.
func WithTimeout(timeout time.Duration) Option {
	return func(config config) config {
		config.timeout = timeout.String()
		return config
	}
}

func (c *config) toArgs() []string {
	var args []string

	args = appendNonEmptyStringVal(args, "--config", c.config)
	args = appendNonEmptyStringVal(args, "--id", c.id)
	args = appendNonEmptyStringVal(args, "--output", c.output)
	args = appendNonEmptyStringVal(args, "--parallelism", c.parallelism)
	args = appendBoolVal(args, "--rm-dist", c.rmDist)
	args = appendBoolVal(args, "--single-target", c.singleTarget)
	args = appendBoolVal(args, "--skip-post-hooks", c.skipPostCommitHooks)
	args = appendBoolVal(args, "--skip-validate", c.skipValidate)
	args = appendBoolVal(args, "--snapshot", c.snapshot)
	args = appendNonEmptyStringVal(args, "--timeout", c.timeout)

	return args
}

func appendNonEmptyStringVal(slice []string, flag, val string) []string {
	// if val is empty return slice as it's
	if val == "" {
		return slice
	}

	return append(slice, flag, val)
}

func appendBoolVal(slice []string, flag string, val bool) []string {
	// if val is false, no need to append the flag
	if !val {
		return slice
	}

	return append(slice, flag)
}
