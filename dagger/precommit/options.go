package precommit

import (
	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
)

type Config struct {
	BaseImage            string `env:"PRECOMMIT_BASE_IMAGE" envDefault:"python:3.12.0a1-bullseye"`
	ContainerCustomizers []containers.ContainerCustomizerFn
}

// BaseImage sets the base image for the precommit container.
func BaseImage(img string) daggers.Option[Config] {
	return func(c Config) Config {
		c.BaseImage = img
		return c
	}
}

// CustomizeContainer adds a customizer function to the precommit container.
func CustomizeContainer(customizers ...containers.ContainerCustomizerFn) daggers.Option[Config] {
	return func(c Config) Config {
		c.ContainerCustomizers = append(c.ContainerCustomizers, customizers...)
		return c
	}
}
