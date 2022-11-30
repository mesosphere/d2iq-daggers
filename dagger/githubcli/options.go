package githubcli

import (
	"github.com/mesosphere/daggers/daggers"
	"github.com/mesosphere/daggers/daggers/containers"
)

type Config struct {
	GoImageRepo      string   `env:"GO_IMAGE_REPO,notEmpty" envDefault:"docker.io/golang"`
	GoImageTag       string   `env:"GO_IMAGE_TAG,notEmpty" envDefault:"1.19"`
	GithubCliVersion string   `env:"GH_VERSION,notEmpty" envDefault:"2.20.2"`
	Extensions       []string `env:"GH_EXTENSIONS" envDefault:""`
	Args             []string `env:"GH_ARGS" envDefault:""  envSeparator:" "`
	MountWorkDir     bool     `env:"GH_MOUNT_WORKDIR" envDefault:"true"`

	Env                  map[string]string
	ContainerCustomizers []containers.ContainerCustomizerFn
}

// WithGoImageRepo sets whether to enable go module caching. Optional, defaults to docker.io/golang.
func WithGoImageRepo(repo string) daggers.Option[Config] {
	return func(c Config) Config {
		c.GoImageRepo = repo
		return c
	}
}

// WithGoImageTag sets the go image tag to use for the container. Optional, defaults to 1.19.
func WithGoImageTag(tag string) daggers.Option[Config] {
	return func(c Config) Config {
		c.GoImageTag = tag
		return c
	}
}

// WithGithubCliVersion sets the github cli version to use for the container.
func WithGithubCliVersion(version string) daggers.Option[Config] {
	return func(c Config) Config {
		c.GithubCliVersion = version
		return c
	}
}

// WithExtensions sets the extensions to install for github cli.
func WithExtensions(extensions ...string) daggers.Option[Config] {
	return func(c Config) Config {
		c.Extensions = extensions
		return c
	}
}

// WithArgs sets the arguments to pass to github cli.
func WithArgs(args ...string) daggers.Option[Config] {
	return func(c Config) Config {
		c.Args = args
		return c
	}
}

// WithMountWorkDir sets whether to mount runtime workdir to the container. Optional, defaults to true.
func WithMountWorkDir(mount bool) daggers.Option[Config] {
	return func(c Config) Config {
		c.MountWorkDir = mount
		return c
	}
}

// WithEnv sets the environment variables to pass to go.
func WithEnv(envMap map[string]string) daggers.Option[Config] {
	return func(c Config) Config {
		c.Env = envMap
		return c
	}
}

// WithContainerCustomizers adds the container customizers to use for the container.
func WithContainerCustomizers(customizers ...containers.ContainerCustomizerFn) daggers.Option[Config] {
	return func(c Config) Config {
		c.ContainerCustomizers = append(c.ContainerCustomizers, customizers...)
		return c
	}
}
