package golang

import "github.com/mesosphere/daggers/daggers"

type config struct {
	GoBaseImage string   `env:"GO_BASE_IMAGE,notEmpty" envDefault:"docker.io/golang"`
	GoVersion   string   `env:"GO_VERSION,notEmpty" envDefault:"1.19"`
	Args        []string `env:"GO_ARGS" envDefault:""  envSeparator:" "`
	Env         map[string]string
}

// WithGoBaseImage sets the go base image to use for the container.
func WithGoBaseImage(image string) daggers.Option[config] {
	return func(c config) config {
		c.GoBaseImage = image
		return c
	}
}

// WithGoVersion sets the go version to use for the container.
func WithGoVersion(version string) daggers.Option[config] {
	return func(c config) config {
		c.GoVersion = version
		return c
	}
}

// WithArgs sets the arguments to pass to go.
func WithArgs(args ...string) daggers.Option[config] {
	return func(c config) config {
		c.Args = args
		return c
	}
}

// WithEnv sets the environment variables to pass to go.
func WithEnv(envMap map[string]string) daggers.Option[config] {
	return func(c config) config {
		c.Env = envMap
		return c
	}
}
