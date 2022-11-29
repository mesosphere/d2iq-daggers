package golang

import "github.com/mesosphere/daggers/daggers"

type config struct {
	GoImageRepo string   `env:"GO_IMAGE_REPO,notEmpty" envDefault:"docker.io/golang"`
	GoImageTag  string   `env:"GO_IMAGE_TAG,notEmpty" envDefault:"1.19"`
	Args        []string `env:"ARGS" envDefault:""  envSeparator:" "`
	Env         map[string]string
}

// WithGoImageRepo sets the go base image to use for the container.
func WithGoImageRepo(image string) daggers.Option[config] {
	return func(c config) config {
		c.GoImageRepo = image
		return c
	}
}

// WithGoImageTag sets the go version to use for the container.
func WithGoImageTag(version string) daggers.Option[config] {
	return func(c config) config {
		c.GoImageTag = version
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
