package golang

import (
	"github.com/caarlos0/env/v6"

	"github.com/mesosphere/daggers/dagger/common"
)

type config struct {
	common.GolangImageConfig
	common.GithubCLIConfig

	Args []string `env:"GO_ARGS" envDefault:""  envSeparator:" "`
	Env  map[string]string
}

func loadConfigFromEnv() (config, error) {
	cfg := config{}

	if err := env.Parse(&cfg, env.Options{}); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Option is a function that configures the precommit checks.
type Option func(config) config

// WithGoImageRepo sets the go image repo to use for the container.
func WithGoImageRepo(image string) Option {
	return func(c config) config {
		c.GoImageRepo = image
		return c
	}
}

// WithGoImageTag sets the go image tag to use for the container.
func WithGoImageTag(version string) Option {
	return func(c config) config {
		c.GoImageTag = version
		return c
	}
}

// WithGithubCliVersion sets the github cli version to use for the container.
func WithGithubCliVersion(version string) Option {
	return func(c config) config {
		c.GithubCliVersion = version
		return c
	}
}

// WithExtensions sets the extensions to install for github cli.
func WithExtensions(extensions ...string) Option {
	return func(c config) config {
		c.GithubCliExtensions = extensions
		return c
	}
}

// WithArgs sets the arguments to pass to go.
func WithArgs(args ...string) Option {
	return func(c config) config {
		c.Args = args
		return c
	}
}

// WithEnv sets the environment variables to pass to go.
func WithEnv(envMap map[string]string) Option {
	return func(c config) config {
		c.Env = envMap
		return c
	}
}
