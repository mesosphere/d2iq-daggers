// Copyright 2022 D2iQ, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package precommit

import (
	"context"
	"embed"
	"fmt"
	"io"

	"dagger.io/dagger"

	"github.com/mesosphere/daggers-for-dkp/daggers"
	"github.com/mesosphere/daggers-for-dkp/daggers/containers"
)

const (
	configFileName      = "pre-commit-config.yaml"
	cacheDir            = "/pre-commit-cache"
	precommitHomeEnvVar = "PRE_COMMIT_HOME"
	precommitVersion    = "3.2.1"
)

//go:embed pre-commit-config.yaml
var configFile embed.FS

// Run runs the precommit checks.
func Run(ctx context.Context, runtime *daggers.Runtime, opts ...daggers.Option[config]) (string, error) {
	cfg, err := daggers.InitConfig(opts...)
	if err != nil {
		return "", err
	}

	var (
		url = fmt.Sprintf(
			"https://github.com/pre-commit/pre-commit/releases/download/v%[1]s/pre-commit-%[1]s.pyz",
			precommitVersion,
		)
		dest        = fmt.Sprintf("/usr/local/bin/pre-commit-%s.pyz", precommitVersion)
		envFn       = containers.WithEnvVariables(cfg.Env)
		customizers = []containers.ContainerCustomizerFn{envFn}
	)

	customizers = append(customizers, cfg.ContainerCustomizers...)

	customizers = append(
		customizers,
		containers.DownloadFile(url, dest),
	)

	container, err := containers.CustomizedContainerFromImage(ctx, runtime, cfg.BaseImage, true, customizers...)
	if err != nil {
		return "", err
	}

	config, err := configFile.Open(configFileName)
	if err != nil {
		return "", err
	}

	configContent, err := io.ReadAll(config)
	if err != nil {
		return "", err
	}

	container = container.
		WithEnvVariable(precommitHomeEnvVar, cacheDir).
		WithNewFile("."+configFileName, dagger.ContainerWithNewFileOpts{
			Contents: string(configContent),
		}).
		WithExec(
			[]string{
				"python",
				fmt.Sprintf("/usr/local/bin/pre-commit-%s.pyz", precommitVersion),
				"run", "--all-files", "--show-diff-on-failure",
			},
		)

	// Run container and get Exit code
	return container.Stdout(ctx)
}
