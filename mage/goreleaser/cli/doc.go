// Package cli provides a wrapper around goreleaser cli to execute goreleaser commands using mage/sh
//
// Our goreleaser flow is contains docker image build and push and this is not possible to do with dagger at the
// moment. We will need to add this feature to dagger after https://github.com/dagger/dagger/issues/3712 resolved.
// Currently, we are using pure mage to execute goreleaser commands.
package cli
