// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package svu

import (
	"context"
	"fmt"
	"strings"

	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
)

// Output is svu command output.
type Output struct {
	// Version
	Version string
	// Version without the prefix
	VersionWithoutPrefix string
}

// Run runs the svu command with the given options.
func Run(ctx context.Context, runtime *daggers.Runtime, options ...daggers.Option[config]) (*Output, error) {
	cfg, err := daggers.InitConfig(options...)
	if err != nil {
		return nil, err
	}

	var (
		image    = fmt.Sprintf("ghcr.io/caarlos0/svu:%s", cfg.Version)
		svuFlags = cfg.toArgs()
	)

	container, err := containers.CustomizedContainerFromImage(ctx, runtime, image, true)
	if err != nil {
		return nil, err
	}

	container = container.WithExec(append([]string{cfg.Command}, svuFlags...))

	version, err := container.Stdout(ctx)
	if err != nil {
		return nil, err
	}

	svuFlags = append(svuFlags, "--strip-prefix")
	container = container.WithExec(append([]string{cfg.Command}, svuFlags...))

	versionWithoutPrefix, err := container.Stdout(ctx)
	if err != nil {
		return nil, err
	}

	return &Output{
		Version:              strings.TrimSpace(version),
		VersionWithoutPrefix: strings.TrimSpace(versionWithoutPrefix),
	}, nil
}
