package build

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Env                 map[string]string
	Config              string `env:"CONFIG"`
	ID                  string `env:"ID"`
	Output              string `env:"OUTPUT"`
	Parallelism         string `env:"PARALLELISM"`
	RmDist              bool   `env:"RM_DIST"`
	SingleTarget        bool   `env:"SINGLE_TARGET"`
	SkipPostCommitHooks bool   `env:"SKIP_POST_COMMIT_HOOKS"`
	SkipValidate        bool   `env:"SKIP_VALIDATE"`
	Snapshot            bool   `env:"SNAPSHOT"`
	Timeout             string `env:"TIMEOUT"`
}

// Option is a function that configures the goreleaser build options.
type Option func(config config) config

func loadConfigFromEnv() (config, error) {
	cfg := config{}

	if err := env.Parse(&cfg, env.Options{Prefix: "GORELEASER_BUILD_"}); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// WithEnv append extra env variables to goreleaser build process.
func WithEnv(envMap map[string]string) Option {
	return func(config config) config {
		config.Env = envMap
		return config
	}
}

// WithConfig sets  --config flag.
func WithConfig(configPath string) Option {
	return func(config config) config {
		config.Config = configPath
		return config
	}
}

// WithID sets --id flag.
func WithID(id string) Option {
	return func(config config) config {
		config.ID = id
		return config
	}
}

// WithOutput sets --output.
func WithOutput(output string) Option {
	return func(config config) config {
		config.Output = output
		return config
	}
}

// WithParallelism sets --parallelism.
func WithParallelism(parallelism string) Option {
	return func(config config) config {
		config.Parallelism = parallelism
		return config
	}
}

// WithRmDist sets --rm-dist.
func WithRmDist(rmDist bool) Option {
	return func(config config) config {
		config.RmDist = rmDist
		return config
	}
}

// WithSingleTarget sets --single-target.
func WithSingleTarget(singleTarget bool) Option {
	return func(config config) config {
		config.SingleTarget = singleTarget
		return config
	}
}

// SkipPostCommitHooks sets--skip-post-hooks.
func SkipPostCommitHooks(skipPostCommitHooks bool) Option {
	return func(config config) config {
		config.SkipPostCommitHooks = skipPostCommitHooks
		return config
	}
}

// SkipValidate sets --skip-validate.
func SkipValidate(skipValidate bool) Option {
	return func(config config) config {
		config.SkipValidate = skipValidate
		return config
	}
}

// WithSnapshot sets --snapshot.
func WithSnapshot(snapshot bool) Option {
	return func(config config) config {
		config.Snapshot = snapshot
		return config
	}
}

// WithTimeout sets --timeout duration.
func WithTimeout(timeout time.Duration) Option {
	return func(config config) config {
		config.Timeout = timeout.String()
		return config
	}
}

func (c *config) toArgs() []string {
	var args []string

	args = appendNonEmptyStringVal(args, "--config", c.Config)
	args = appendNonEmptyStringVal(args, "--id", c.ID)
	args = appendNonEmptyStringVal(args, "--output", c.Output)
	args = appendNonEmptyStringVal(args, "--parallelism", c.Parallelism)
	args = appendBoolVal(args, "--rm-dist", c.RmDist)
	args = appendBoolVal(args, "--single-target", c.SingleTarget)
	args = appendBoolVal(args, "--skip-post-hooks", c.SkipPostCommitHooks)
	args = appendBoolVal(args, "--skip-validate", c.SkipValidate)
	args = appendBoolVal(args, "--snapshot", c.Snapshot)
	args = appendNonEmptyStringVal(args, "--timeout", c.Timeout)

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
