package precommit

import "dagger.io/dagger"

type config struct {
	baseImage            string
	containerCustomizers []ContainerCustomizer
}

func defaultConfig() config {
	return config{
		baseImage: "python:3.12.0a1-bullseye",
	}
}

type Option func(config) config

type ContainerCustomizer func(*dagger.Container) *dagger.Container

func BaseImage(img string) Option {
	return func(c config) config {
		c.baseImage = img
		return c
	}
}

func CustomizeContainer(customizers ...ContainerCustomizer) Option {
	return func(c config) config {
		c.containerCustomizers = append(c.containerCustomizers, customizers...)
		return c
	}
}
