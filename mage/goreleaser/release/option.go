package release

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Env               map[string]string
	AutoSnapshot      bool   `env:"AUTO_SNAPSHOT"`
	Config            string `env:"CONFIG"`
	Parallelism       string `env:"PARALLELISM"`
	RmDist            bool   `env:"RM_DIST"`
	ReleaseFooter     string `env:"RELEASE_FOOTER"`
	ReleaseFooterTmpl string `env:"RELEASE_FOOTER_TMPL"`
	ReleaseHeader     string `env:"RELEASE_HEADER"`
	ReleaseHeaderTmpl string `env:"RELEASE_HEADER_TMPL"`
	ReleaseNotes      string `env:"RELEASE_NOTES"`
	ReleaseNotesTmpl  string `env:"RELEASE_NOTES_TMPL"`
	SkipAnnounce      bool   `env:"SKIP_ANNOUNCE"`
	SkipPublish       bool   `env:"SKIP_PUBLISH"`
	SkipSbom          bool   `env:"SKIP_SBOM"`
	SkipSign          bool   `env:"SKIP_SIGN"`
	SkipValidate      bool   `env:"SKIP_VALIDATE"`
	Snapshot          bool   `env:"SNAPSHOT"`
	Timeout           string `env:"TIMEOUT"`
}

// Option is a function that configures the goreleaser release options.
type Option func(config config) config

func loadConfigFromEnv() (config, error) {
	cfg := config{}

	if err := env.Parse(&cfg, env.Options{Prefix: "GORELEASER_RELEASE_"}); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// WithEnv append extra env variables to goreleaser build process.
func WithEnv(envMap map[string]string) Option {
	return func(config config) config {
		config.Env = envMap
		return config
	}
}

// WithAutoSnapshot sets --snapshot flag.
func WithAutoSnapshot(autoSnapshot bool) Option {
	return func(config config) config {
		config.AutoSnapshot = autoSnapshot
		return config
	}
}

// WithConfig sets  --config flag.
func WithConfig(configPath string) Option {
	return func(config config) config {
		config.Config = configPath
		return config
	}
}

// WithParallelism sets --parallelism flag.
func WithParallelism(parallelism string) Option {
	return func(config config) config {
		config.Parallelism = parallelism
		return config
	}
}

// WithRmDist sets --rm-dist flag.
func WithRmDist(rmDist bool) Option {
	return func(config config) config {
		config.RmDist = rmDist
		return config
	}
}

// WithReleaseFooter sets --release-footer flag.
func WithReleaseFooter(releaseFooter string) Option {
	return func(config config) config {
		config.ReleaseFooter = releaseFooter
		return config
	}
}

// WithReleaseFooterTmpl sets --release-footer-tmpl flag.
func WithReleaseFooterTmpl(releaseFooterTmpl string) Option {
	return func(config config) config {
		config.ReleaseFooterTmpl = releaseFooterTmpl
		return config
	}
}

// WithReleaseHeader sets --release-header flag.
func WithReleaseHeader(releaseHeader string) Option {
	return func(config config) config {
		config.ReleaseHeader = releaseHeader
		return config
	}
}

// WithReleaseHeaderTmpl sets --release-header-tmpl flag.
func WithReleaseHeaderTmpl(releaseHeaderTmpl string) Option {
	return func(config config) config {
		config.ReleaseHeaderTmpl = releaseHeaderTmpl
		return config
	}
}

// WithReleaseNotes sets --release-notes flag.
func WithReleaseNotes(releaseNotes string) Option {
	return func(config config) config {
		config.ReleaseNotes = releaseNotes
		return config
	}
}

// WithReleaseNotesTmpl sets --release-notes-tmpl flag.
func WithReleaseNotesTmpl(releaseNotesTmpl string) Option {
	return func(config config) config {
		config.ReleaseNotesTmpl = releaseNotesTmpl
		return config
	}
}

// SkipAnnounce sets --skip-announce flag.
func SkipAnnounce(skipAnnounce bool) Option {
	return func(config config) config {
		config.SkipAnnounce = skipAnnounce
		return config
	}
}

// SkipPublish sets --skip-publish flag.
func SkipPublish(skipPublish bool) Option {
	return func(config config) config {
		config.SkipPublish = skipPublish
		return config
	}
}

// SkipSbom sets --skip-sbom flag.
func SkipSbom(skipSbom bool) Option {
	return func(config config) config {
		config.SkipSbom = skipSbom
		return config
	}
}

// SkipSign sets --skip-sign flag.
func SkipSign(skipSign bool) Option {
	return func(config config) config {
		config.SkipSign = skipSign
		return config
	}
}

// SkipValidate sets --skip-validate flag.
func SkipValidate(skipValidate bool) Option {
	return func(config config) config {
		config.SkipValidate = skipValidate
		return config
	}
}

// WithSnapshot sets --snapshot flag.
func WithSnapshot(snapshot bool) Option {
	return func(config config) config {
		config.Snapshot = snapshot
		return config
	}
}

// WithTimeout sets --timeout flag.
func WithTimeout(timeout time.Duration) Option {
	return func(config config) config {
		config.Timeout = timeout.String()
		return config
	}
}

func (c *config) toArgs() []string {
	var args []string

	args = appendNonEmptyStringVal(args, "--config", c.Config)
	args = appendBoolVal(args, "--snapshot", c.AutoSnapshot)
	args = appendNonEmptyStringVal(args, "--parallelism", c.Parallelism)
	args = appendBoolVal(args, "--rm-dist", c.RmDist)
	args = appendNonEmptyStringVal(args, "--release-footer", c.ReleaseFooter)
	args = appendNonEmptyStringVal(args, "--release-footer-tmpl", c.ReleaseFooterTmpl)
	args = appendNonEmptyStringVal(args, "--release-header", c.ReleaseHeader)
	args = appendNonEmptyStringVal(args, "--release-header-tmpl", c.ReleaseHeaderTmpl)
	args = appendNonEmptyStringVal(args, "--release-notes", c.ReleaseNotes)
	args = appendNonEmptyStringVal(args, "--release-notes-tmpl", c.ReleaseNotesTmpl)
	args = appendBoolVal(args, "--skip-announce", c.SkipAnnounce)
	args = appendBoolVal(args, "--skip-publish", c.SkipPublish)
	args = appendBoolVal(args, "--skip-sbom", c.SkipSbom)
	args = appendBoolVal(args, "--skip-sign", c.SkipSign)
	args = appendBoolVal(args, "--skip-validate", c.SkipValidate)
	args = appendBoolVal(args, "--snapshot", c.Snapshot)
	args = appendNonEmptyStringVal(args, "--timeout", c.Timeout)

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
