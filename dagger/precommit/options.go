package precommit

import (
	"github.com/mesosphere/daggers/dagger/options"
	"github.com/mesosphere/daggers/daggers"
)

type config struct {
	BaseImage            string `env:"PRECOMMIT_BASE_IMAGE" envDefault:"python:3.12.0a1-bullseye"`
	ContainerCustomizers []options.ContainerCustomizer
}

// BaseImage sets the base image for the precommit container.
func BaseImage(img string) daggers.Option[config] {
	return func(c config) config {
		c.BaseImage = img
		return c
	}
}

// CustomizeContainer adds a customizer function to the precommit container.
func CustomizeContainer(customizers ...options.ContainerCustomizer) daggers.Option[config] {
	return func(c config) config {
		c.ContainerCustomizers = append(c.ContainerCustomizers, customizers...)
		return c
	}
}
