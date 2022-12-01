package githubcli

import "github.com/mesosphere/daggers/daggers"

type config struct {
	GoImageRepo      string   `env:"GO_IMAGE_REPO,notEmpty" envDefault:"docker.io/golang"`
	GoImageTag       string   `env:"GO_IMAGE_TAG,notEmpty" envDefault:"1.19"`
	GithubCliVersion string   `env:"GH_VERSION,notEmpty" envDefault:"2.20.2"`
	Extensions       []string `env:"GH_EXTENSIONS" envDefault:""`
	Args             []string `env:"GH_ARGS" envDefault:""  envSeparator:" "`
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
