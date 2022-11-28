package golang

import "github.com/caarlos0/env/v6"

type config struct {
	GoBaseImage string   `env:"BASE_IMAGE,notEmpty" envDefault:"docker.io/golang"`
	GoVersion   string   `env:"VERSION,notEmpty" envDefault:"1.19"`
	Args        []string `env:"ARGS" envDefault:""  envSeparator:" "`
	Env         map[string]string
}

func loadConfigFromEnv() (config, error) {
	cfg := config{}

	if err := env.Parse(&cfg, env.Options{Prefix: "GO_"}); err != nil {
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
