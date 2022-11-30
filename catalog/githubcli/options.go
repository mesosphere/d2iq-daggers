package githubcli

import "github.com/mesosphere/daggers/daggers"

type config struct {
	GoBaseImage      string   `env:"GO_BASE_IMAGE,notEmpty" envDefault:"docker.io/golang"`
	GoVersion        string   `env:"GO_VERSION,notEmpty" envDefault:"1.19"`
	GithubCliVersion string   `env:"GH_VERSION,notEmpty" envDefault:"2.20.2"`
	Extensions       []string `env:"GH_EXTENSIONS" envDefault:""`
	Args             []string `env:"GH_ARGS" envDefault:""  envSeparator:" "`
}

// WithGoBaseImage sets the go base image to use for the container.
func WithGoBaseImage(image string) daggers.Option[config] {
	return func(c config) config {
		c.GoBaseImage = image
		return c
	}
}

// WithGoVersion sets the go version to use for the container.
func WithGoVersion(version string) daggers.Option[config] {
	return func(c config) config {
		c.GoVersion = version
		return c
	}
}

// WithGithubCliVersion sets the github cli version to use for the container.
func WithGithubCliVersion(version string) daggers.Option[config] {
	return func(c config) config {
		c.GithubCliVersion = version
		return c
	}
}

// WithExtensions sets the extensions to install for github cli.
func WithExtensions(extensions ...string) daggers.Option[config] {
	return func(c config) config {
		c.Extensions = extensions
		return c
	}
}

// WithArgs sets the arguments to pass to github cli.
func WithArgs(args ...string) daggers.Option[config] {
	return func(c config) config {
		c.Args = args
		return c
	}
}
