package release

import (
	"time"
)

type config struct {
	env               map[string]string
	autoSnapshot      bool
	config            string
	parallelism       string
	rmDist            bool
	releaseFooter     string
	releaseFooterTmpl string
	releaseHeader     string
	releaseHeaderTmpl string
	releaseNotes      string
	releaseNotesTmpl  string
	skipAnnounce      bool
	skipPublish       bool
	skipSbom          bool
	skipSign          bool
	skipValidate      bool
	snapshot          bool
	timeout           string
}

// Option is a function that configures the goreleaser release options.
type Option func(config config) config

// WithEnv append extra env variables to goreleaser build process.
func WithEnv(env map[string]string) Option {
	return func(config config) config {
		config.env = env
		return config
	}
}

// WithAutoSnapshot sets --snapshot flag.
func WithAutoSnapshot(autoSnapshot bool) Option {
	return func(config config) config {
		config.autoSnapshot = autoSnapshot
		return config
	}
}

// WithConfig sets  --config flag.
func WithConfig(configPath string) Option {
	return func(config config) config {
		config.config = configPath
		return config
	}
}

// WithParallelism sets --parallelism flag.
func WithParallelism(parallelism string) Option {
	return func(config config) config {
		config.parallelism = parallelism
		return config
	}
}

// WithRmDist sets --rm-dist flag.
func WithRmDist(rmDist bool) Option {
	return func(config config) config {
		config.rmDist = rmDist
		return config
	}
}

// WithReleaseFooter sets --release-footer flag.
func WithReleaseFooter(releaseFooter string) Option {
	return func(config config) config {
		config.releaseFooter = releaseFooter
		return config
	}
}

// WithReleaseFooterTmpl sets --release-footer-tmpl flag.
func WithReleaseFooterTmpl(releaseFooterTmpl string) Option {
	return func(config config) config {
		config.releaseFooterTmpl = releaseFooterTmpl
		return config
	}
}

// WithReleaseHeader sets --release-header flag.
func WithReleaseHeader(releaseHeader string) Option {
	return func(config config) config {
		config.releaseHeader = releaseHeader
		return config
	}
}

// WithReleaseHeaderTmpl sets --release-header-tmpl flag.
func WithReleaseHeaderTmpl(releaseHeaderTmpl string) Option {
	return func(config config) config {
		config.releaseHeaderTmpl = releaseHeaderTmpl
		return config
	}
}

// WithReleaseNotes sets --release-notes flag.
func WithReleaseNotes(releaseNotes string) Option {
	return func(config config) config {
		config.releaseNotes = releaseNotes
		return config
	}
}

// WithReleaseNotesTmpl sets --release-notes-tmpl flag.
func WithReleaseNotesTmpl(releaseNotesTmpl string) Option {
	return func(config config) config {
		config.releaseNotesTmpl = releaseNotesTmpl
		return config
	}
}

// SkipAnnounce sets --skip-announce flag.
func SkipAnnounce(skipAnnounce bool) Option {
	return func(config config) config {
		config.skipAnnounce = skipAnnounce
		return config
	}
}

// SkipPublish sets --skip-publish flag.
func SkipPublish(skipPublish bool) Option {
	return func(config config) config {
		config.skipPublish = skipPublish
		return config
	}
}

// SkipSbom sets --skip-sbom flag.
func SkipSbom(skipSbom bool) Option {
	return func(config config) config {
		config.skipSbom = skipSbom
		return config
	}
}

// SkipSign sets --skip-sign flag.
func SkipSign(skipSign bool) Option {
	return func(config config) config {
		config.skipSign = skipSign
		return config
	}
}

// SkipValidate sets --skip-validate flag.
func SkipValidate(skipValidate bool) Option {
	return func(config config) config {
		config.skipValidate = skipValidate
		return config
	}
}

// WithSnapshot sets --snapshot flag.
func WithSnapshot(snapshot bool) Option {
	return func(config config) config {
		config.snapshot = snapshot
		return config
	}
}

// WithTimeout sets --timeout flag.
func WithTimeout(timeout time.Duration) Option {
	return func(config config) config {
		config.timeout = timeout.String()
		return config
	}
}

func (c *config) toArgs() []string {
	var args []string

	args = appendNonEmptyStringVal(args, "--config", c.config)
	args = appendBoolVal(args, "--snapshot", c.autoSnapshot)
	args = appendNonEmptyStringVal(args, "--parallelism", c.parallelism)
	args = appendBoolVal(args, "--rm-dist", c.rmDist)
	args = appendNonEmptyStringVal(args, "--release-footer", c.releaseFooter)
	args = appendNonEmptyStringVal(args, "--release-footer-tmpl", c.releaseFooterTmpl)
	args = appendNonEmptyStringVal(args, "--release-header", c.releaseHeader)
	args = appendNonEmptyStringVal(args, "--release-header-tmpl", c.releaseHeaderTmpl)
	args = appendNonEmptyStringVal(args, "--release-notes", c.releaseNotes)
	args = appendNonEmptyStringVal(args, "--release-notes-tmpl", c.releaseNotesTmpl)
	args = appendBoolVal(args, "--skip-announce", c.skipAnnounce)
	args = appendBoolVal(args, "--skip-publish", c.skipPublish)
	args = appendBoolVal(args, "--skip-sbom", c.skipSbom)
	args = appendBoolVal(args, "--skip-sign", c.skipSign)
	args = appendBoolVal(args, "--skip-validate", c.skipValidate)
	args = appendBoolVal(args, "--snapshot", c.snapshot)
	args = appendNonEmptyStringVal(args, "--timeout", c.timeout)

	return args
}

func appendNonEmptyStringVal(slice []string, flag, val string) []string {
	// if val is empty return slice as it's
	if val == "" {
		return slice
	}

	return append(slice, flag, val)
}

func appendBoolVal(slice []string, flag string, val bool) []string {
	// if val is false, no need to append the flag
	if !val {
		return slice
	}

	return append(slice, flag)
}
