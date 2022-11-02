package precommit

import "github.com/aweris/tools/dagger/options"

type config struct {
	baseImage            string
	containerCustomizers []options.ContainerCustomizer
}

func defaultConfig() config {
	return config{
		baseImage: "python:3.12.0a1-bullseye",
	}
}

type Option func(config) config

func BaseImage(img string) Option {
	return func(c config) config {
		c.baseImage = img
		return c
	}
}

func CustomizeContainer(customizers ...options.ContainerCustomizer) Option {
	return func(c config) config {
		c.containerCustomizers = append(c.containerCustomizers, customizers...)
		return c
	}
}
