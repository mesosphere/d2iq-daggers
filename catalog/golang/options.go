// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package golang

import (
	"github.com/mesosphere/daggers-for-dkp/daggers"
	"github.com/mesosphere/daggers-for-dkp/daggers/containers"
)

type config struct {
	GoImageRepo       string   `env:"GO_IMAGE_REPO,notEmpty" envDefault:"docker.io/golang"`
	GoImageTag        string   `env:"GO_IMAGE_TAG,notEmpty" envDefault:"1.20"`
	GoModCacheEnabled bool     `env:"GO_MOD_CACHE_ENABLE" envDefault:"true"`
	GoModDir          string   `env:"GO_MOD_DIR" envDefault:"."`
	Args              []string `env:"GO_ARGS" envDefault:""  envSeparator:" "`

	Env                  map[string]string
	ContainerCustomizers []containers.ContainerCustomizerFn
}

// WithGoImageRepo sets whether to enable go module caching. Optional, defaults to docker.io/golang.
func WithGoImageRepo(repo string) daggers.Option[config] {
	return func(c config) config {
		c.GoImageRepo = repo
		return c
	}
}

// WithGoImageTag sets the go image tag to use for the container. Optional, defaults to 1.19.
func WithGoImageTag(tag string) daggers.Option[config] {
	return func(c config) config {
		c.GoImageTag = tag
		return c
	}
}

// WithGoModCacheEnabled sets whether to enable go module caching. Optional, defaults to true.
func WithGoModCacheEnabled(enable bool) daggers.Option[config] {
	return func(c config) config {
		c.GoModCacheEnabled = enable
		return c
	}
}

// WithGoModDir sets the go module directory to use for the container. Optional, defaults to the current directory.
func WithGoModDir(dir string) daggers.Option[config] {
	return func(c config) config {
		c.GoModDir = dir
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

// WithContainerCustomizers adds the container customizers to use for the container.
func WithContainerCustomizers(customizers ...containers.ContainerCustomizerFn) daggers.Option[config] {
	return func(c config) config {
		c.ContainerCustomizers = append(c.ContainerCustomizers, customizers...)
		return c
	}
}
