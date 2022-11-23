package githubcli

import (
	"github.com/caarlos0/env/v6"
)

type config struct {
	GoBaseImage      string   `env:"GO_BASE_IMAGE,notEmpty" envDefault:"docker.io/golang"`
	GoVersion        string   `env:"GO_VERSION,notEmpty" envDefault:"1.19"`
	GithubCliVersion string   `env:"VERSION,notEmpty" envDefault:"2.20.2"`
	Extensions       []string `env:"EXTENSIONS" envDefault:""`
	Args             []string `env:"ARGS" envDefault:""  envSeparator:" "`
}

func loadConfigFromEnv() (config, error) {
	cfg := config{}

	if err := env.Parse(&cfg, env.Options{Prefix: "GH_"}); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Option is a function that configures the precommit checks.
type Option func(config) config

// WithGoBaseImage sets the go base image to use for the container.
func WithGoBaseImage(image string) Option {
	return func(c config) config {
		c.GoBaseImage = image
		return c
	}
}

// WithGoVersion sets the go version to use for the container.
func WithGoVersion(version string) Option {
	return func(c config) config {
		c.GoVersion = version
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
		c.Extensions = extensions
		return c
	}
}

// WithArgs sets the arguments to pass to github cli.
func WithArgs(args ...string) Option {
	return func(c config) config {
		c.Args = args
		return c
	}
}
