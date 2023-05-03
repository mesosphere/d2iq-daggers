// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package golang

import (
	"context"
	"fmt"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers-for-dkp/daggers"
	"github.com/mesosphere/daggers-for-dkp/daggers/containers"
)

// standard source path.
const srcDir = "/src"

// RunCommand runs a go command with given working directory and options and returns command output and
// working directory.
func RunCommand(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config],
) (string, *dagger.Directory, error) {
	container, err := GetContainer(ctx, runtime, opts...)
	if err != nil {
		return "", nil, err
	}

	out, err := container.Stdout(ctx)
	if err != nil {
		return "", nil, err
	}

	return out, container.Directory(srcDir), nil
}

// GetContainer returns a dagger container with given working directory and options.
func GetContainer(
	ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config],
) (*dagger.Container, error) {
	cfg, err := daggers.InitConfig(opts...)
	if err != nil {
		return nil, err
	}

	var (
		image       = fmt.Sprintf("%s:%s", cfg.GoImageRepo, cfg.GoImageTag)
		envFn       = containers.WithEnvVariables(cfg.Env)
		customizers = []containers.ContainerCustomizerFn{envFn}
	)

	if cfg.GoModCacheEnabled {
		customizers = append(customizers, containers.WithMountedGoCache(ctx, cfg.GoModDir))
	}

	customizers = append(customizers, cfg.ContainerCustomizers...)

	container, err := containers.CustomizedContainerFromImage(ctx, runtime, image, true, customizers...)
	if err != nil {
		return nil, err
	}

	container = container.WithEntrypoint([]string{"go"})

	if len(cfg.Args) > 0 {
		container = container.WithExec(cfg.Args)
	}

	return container, nil
}
