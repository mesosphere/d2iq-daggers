package svu

// Command is represents the svu sub-command.
type Command string

const (
	// CommandNext is the svu next sub-command.
	CommandNext Command = "next"
	// CommandMajor is the svu major sub-command.
	CommandMajor Command = "major"
	// CommandMinor is the svu minor sub-command.
	CommandMinor Command = "minor"
	// CommandPatch is the svu patch sub-command.
	CommandPatch Command = "patch"
	// CommandCurrent is the svu pre sub-command.
	CommandCurrent Command = "current"
)

// TagMode is a custom type representing the possible values for the --tag-mode flag.
type TagMode string

const (
	// TagModeAllBranches is the value for the --tag-mode flag that will use all branches.
	TagModeAllBranches TagMode = "all-branches"
	// TagModeCurrentBranch is the value for the --tag-mode flag that will only use current branch tags.
	TagModeCurrentBranch TagMode = "current-branch"
)

type config struct {
	Version    string
	Metadata   bool
	Prerelease bool
	Build      bool
	Command    Command
	Pattern    string
	Prefix     string
	Suffix     string
	TagMode    TagMode
}

func defaultConfig() config {
	return config{
		Version:    "v1.9.0",
		Metadata:   true,
		Prerelease: true,
		Build:      true,
		Command:    CommandNext,
		Pattern:    "*",
		Prefix:     "v",
		Suffix:     "",
		TagMode:    TagModeAllBranches,
	}
}

// Option is a function that configures the svu action.
type Option func(config) config

// SVUVersion specifies the version of svu to use. Defaults to v1.9.0. This should be one of the
// released image tags - see https://github.com/caarlos0/svu/pkgs/container/svu for available
// tags.
//
//nolint:revive // Disable stuttering check.
func SVUVersion(v string) Option {
	return func(c config) config {
		c.Version = v
		return c
	}
}

// WithMetadata controls whether to include pre-release and build metadata in the version. Defaults to true.
func WithMetadata(b bool) Option {
	return func(c config) config {
		c.Metadata = b
		return c
	}
}

// WithPreRelease controls whether to include pre-release metadata in the version. Defaults to true.
func WithPreRelease(b bool) Option {
	return func(c config) config {
		c.Prerelease = b
		return c
	}
}

// WithBuild controls whether to include build metadata in the version. Defaults to true.
func WithBuild(b bool) Option {
	return func(c config) config {
		c.Build = b
		return c
	}
}

// WithCommand sets the svu sub-command to run. Defaults to "next".
func WithCommand(cmd Command) Option {
	return func(c config) config {
		c.Command = cmd
		return c
	}
}

// WithPattern sets the pattern to use when searching for tags. Defaults to "*".
func WithPattern(pattern string) Option {
	return func(c config) config {
		c.Pattern = pattern
		return c
	}
}

// WithPrefix sets the prefix to use when searching for tags. Defaults to "v".
func WithPrefix(prefix string) Option {
	return func(c config) config {
		c.Prefix = prefix
		return c
	}
}

// WithSuffix sets the suffix to use when searching for tags. Defaults to "".
func WithSuffix(suffix string) Option {
	return func(c config) config {
		c.Suffix = suffix
		return c
	}
}

// WithTagMode sets the tag mode to use when searching for tags. Defaults to TagModeAllBranches.
func WithTagMode(tagMode TagMode) Option {
	return func(c config) config {
		c.TagMode = tagMode
		return c
	}
}
