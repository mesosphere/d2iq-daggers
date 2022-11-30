package build

import (
	"github.com/mesosphere/daggers/daggers"
)

type config struct {
	Env map[string]string

	Args []string `env:"GORELEASER_BUILD_ARGS" envDefault:""  envSeparator:" "`
}

// WithEnv append extra env variables to goreleaser build process.
func WithEnv(envMap map[string]string) daggers.Option[config] {
	return func(config config) config {
		config.Env = envMap
		return config
	}
}

// WithArgs sets args to goreleaser build process.
func WithArgs(args ...string) daggers.Option[config] {
	return func(config config) config {
		config.Args = args
		return config
	}
}
