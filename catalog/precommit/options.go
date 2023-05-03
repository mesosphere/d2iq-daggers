// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package precommit

import (
	"github.com/mesosphere/d2iq-daggers/daggers"
	"github.com/mesosphere/d2iq-daggers/daggers/containers"
)

type config struct {
	BaseImage string `env:"PRECOMMIT_BASE_IMAGE" envDefault:"python:3.12.0a1-bullseye"`

	Env                  map[string]string
	ContainerCustomizers []containers.ContainerCustomizerFn
}

// BaseImage sets the base image for the precommit container.
func BaseImage(img string) daggers.Option[config] {
	return func(c config) config {
		c.BaseImage = img
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

// CustomizeContainer adds a customizer function to the precommit container.
func CustomizeContainer(customizers ...containers.ContainerCustomizerFn) daggers.Option[config] {
	return func(c config) config {
		c.ContainerCustomizers = append(c.ContainerCustomizers, customizers...)
		return c
	}
}
