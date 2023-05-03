package gotest

import (
	"context"
	"os"

	"dagger.io/dagger"
	"github.com/magefile/mage/mg"

	"github.com/mesosphere/d2iq-daggers/catalog/golang"
	"github.com/mesosphere/d2iq-daggers/daggers"
	"github.com/mesosphere/d2iq-daggers/daggers/containers"
)

const (
	// EnvGowork env variable name for go.work.
	EnvGowork = "GOWORK"
	// EnvGoPrivate env variable name for GOPRIVATE.
	EnvGoPrivate = "GOPRIVATE"
)

// Gounit runs unit tests.
func Gounit(ctx context.Context) error {
	verbose := mg.Verbose() || mg.Debug()

	runtime, err := daggers.NewRuntime(ctx, daggers.WithVerbose(verbose))
	if err != nil {
		return err
	}

	// golang container customizer options
	customizers := golang.WithContainerCustomizers(
		containers.WithGithubAuth(ctx),
		containers.WithEnvVariables(map[string]string{
			EnvGowork:    "off",
			EnvGoPrivate: os.Getenv(EnvGoPrivate),
		}),
	)

	// create a golang container
	container, err := golang.GetContainer(ctx, runtime, customizers)
	if err != nil {
		return err
	}

	// execute the unit tests
	return runUnitTests(ctx, container)
}

// runUnitTests runs the unit tests in the container and exports the test results to .output directory.
func runUnitTests(ctx context.Context, container *dagger.Container) error {
	testContainer := container.
		WithExec([]string{"test", "-v", "-race", "-coverprofile", "coverage.txt", "-covermode", "atomic", "./..."}).
		WithExec([]string{"tool", "cover", "-html=coverage.txt", "-o", "coverage.html"})

	_, err := testContainer.ExitCode(ctx) // execute all steps and return the exit code
	if err != nil {
		return err
	}

	srcDir := testContainer.Directory("/src")

	if _, err := srcDir.File("coverage.txt").Export(ctx, ".reports/coverage.txt"); err != nil {
		return err
	}

	if _, err := srcDir.File("coverage.html").Export(ctx, ".reports/coverage.html"); err != nil {
		return err
	}

	return nil
}
