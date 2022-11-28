//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var env = map[string]string{"GOPRIVATE": "github.com/mesosphere"}

type (
	Test mg.Namespace
)

func (Test) Go() error {
	return sh.RunWithV(env, "go", "test", "-v", "-race", "-coverprofile", "coverage.txt", "-covermode", "atomic", "./...")
}
