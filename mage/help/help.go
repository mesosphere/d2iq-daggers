package help

import "fmt"

// Precommit shows the help for precommit.
func Precommit() {
	fmt.Println(`Usage: mage <namespace:>precommit (e.g. mage lint:precommit or mage precommit)

Configurable via the following environment variables:
	PRECOMMIT_BASE_IMAGE:    The base image to run pre-commit in.`)
}

// Svu shows the help for svu.
func Svu() {
	fmt.Println(`Usage: mage <namespace:><command> (e.g. mage svu:current or mage next)

Commands:
	current:  Print the current version.
	next:     Print the next version.
	major:    Print the next major version.
	minor:    Print the next minor version.
	patch:    Print the next patch version.

Configurable via the following environment variables:
	SVU_VERSION:    SVUVersion specifies the version of svu to use.
	SVU_METADATA:   Controls whether to include pre-release and build metadata in the version. Defaults to true.
	SVU_PATTERN:    Sets the pattern to use when searching for tags. Defaults to "*".
	SVU_PRERELEASE: Controls whether to include pre-release metadata in the version. Defaults to true.
	SVU_BUILD:      Controls whether to include build metadata in the version. Defaults to true.
	SVU_PREFIX:     Sets the prefix to use when searching for tags. Defaults to "v".
	SVU_SUFFIX:     Sets the suffix to use when searching for tags. Defaults to "".
	SVU_TAG_MODE:   Sets the tag mode to use when searching for tags. Defaults to "all-branches".`)
}
