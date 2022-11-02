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
	// TagModeCurrentBranch is the value for the --tag-mode flag that will only use current branch tags
	TagModeCurrentBranch TagMode = "current-branch"
)

type config struct {
	metadata   bool
	preRelease bool
	build      bool
	command    Command
	pattern    string
	prefix     string
	suffix     string
	tagMode    TagMode
}

func defaultConfig() config {
	return config{
		metadata:   true,
		preRelease: true,
		build:      true,
		command:    CommandNext,
		pattern:    "*",
		prefix:     "v",
		suffix:     "",
		tagMode:    TagModeAllBranches,
	}
}

type Option func(config) config

// WithMetadata controls whether to include pre-release and build metadata in the version. Defaults to true.
func WithMetadata(b bool) Option {
	return func(c config) config {
		c.metadata = b
		return c
	}
}

// WithPreRelease controls whether to include pre-release metadata in the version. Defaults to true.
func WithPreRelease(b bool) Option {
	return func(c config) config {
		c.preRelease = b
		return c
	}
}

// WithBuild controls whether to include build metadata in the version. Defaults to true.
func WithBuild(b bool) Option {
	return func(c config) config {
		c.build = b
		return c
	}
}

// WithCommand sets the svu sub-command to run. Defaults to "next".
func WithCommand(cmd Command) Option {
	return func(c config) config {
		c.command = cmd
		return c
	}
}

// WithPattern sets the pattern to use when searching for tags. Defaults to "*".
func WithPattern(pattern string) Option {
	return func(c config) config {
		c.pattern = pattern
		return c
	}
}

// WithPrefix sets the prefix to use when searching for tags. Defaults to "v".
func WithPrefix(prefix string) Option {
	return func(c config) config {
		c.prefix = prefix
		return c
	}
}

// WithSuffix sets the suffix to use when searching for tags. Defaults to "".
func WithSuffix(suffix string) Option {
	return func(c config) config {
		c.suffix = suffix
		return c
	}
}

// WithTagMode sets the tag mode to use when searching for tags. Defaults to TagModeAllBranches.
func WithTagMode(tagMode TagMode) Option {
	return func(c config) config {
		c.tagMode = tagMode
		return c
	}
}
