package precommit

import (
	"github.com/caarlos0/env/v6"

	"github.com/mesosphere/daggers/dagger/options"
)

type config struct {
	BaseImage            string `env:"BASE_IMAGE" envDefault:"python:3.12.0a1-bullseye"`
	ContainerCustomizers []options.ContainerCustomizer
}

func loadConfigFromEnv() (config, error) {
	cfg := config{}

	if err := env.Parse(&cfg, env.Options{Prefix: "PRECOMMIT_"}); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Option is a function that configures the precommit checks.
type Option func(config) config

// BaseImage sets the base image for the precommit container.
func BaseImage(img string) Option {
	return func(c config) config {
		c.BaseImage = img
		return c
	}
}

// CustomizeContainer adds a customizer function to the precommit container.
func CustomizeContainer(customizers ...options.ContainerCustomizer) Option {
	return func(c config) config {
		c.ContainerCustomizers = append(c.ContainerCustomizers, customizers...)
		return c
	}
}
