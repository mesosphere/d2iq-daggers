package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/magefile/mage/sh"
)

// Command a simple wrapper type for goreleaser commands to execute.
type Command string

const (
	// CommandBuild is a custom type for `build` command.
	CommandBuild Command = "build"

	// CommandRelease is a custom type for `release` command.
	CommandRelease Command = "release"
)

// Artifact is a goreleaser type defines a single artifact.
// Copied from https://github.com/goreleaser/goreleaser/blob/main/internal/artifact/artifact.go#L159
type Artifact struct {
	Name    string         `json:"name,omitempty"`
	Path    string         `json:"path,omitempty"`
	Goos    string         `json:"goos,omitempty"`
	Goarch  string         `json:"goarch,omitempty"`
	Goarm   string         `json:"goarm,omitempty"`
	Gomips  string         `json:"gomips,omitempty"`
	Goamd64 string         `json:"goamd64,omitempty"`
	Type    string         `json:"type,omitempty"`
	Extra   map[string]any `json:"extra,omitempty"`
}

// Metadata is a goreleaser type defines metadata.json
// Copied from https://github.com/goreleaser/goreleaser/blob/main/internal/pipe/metadata/metadata.go#L62
type Metadata struct {
	ProjectName string      `json:"project_name"`
	Tag         string      `json:"tag"`
	PreviousTag string      `json:"previous_tag"`
	Version     string      `json:"version"`
	Commit      string      `json:"commit"`
	Date        time.Time   `json:"date"`
	Runtime     MetaRuntime `json:"runtime"`
}

// MetaRuntime is a goreleaser type defines runtime info in metadata.json
// Copied from https://github.com/goreleaser/goreleaser/blob/main/internal/pipe/metadata/metadata.go#L72
type MetaRuntime struct {
	Goos   string `json:"goos"`
	Goarch string `json:"goarch"`
}

// Result represents is goreleaser combined output for metadata.json and artifacts.json.
type Result struct {
	Metadata  Metadata
	Artifacts []Artifact
}

// Run executes goreleaser with given command and arguments and return results info about command.
func Run(cmd Command, debug bool, env map[string]string, args []string) (*Result, error) {
	var cliArgs []string

	if debug {
		cliArgs = append(cliArgs, "--debug")
	}

	cliArgs = append(cliArgs, string(cmd))
	cliArgs = append(cliArgs, args...)

	fmt.Printf("Running goreleaser with args: %v env: %v", cliArgs, env)

	_, err := sh.Exec(env, os.Stdout, os.Stderr, "goreleaser", cliArgs...)
	if err != nil {
		return nil, err
	}

	metadataBytes, err := os.ReadFile("dist/metadata.json")
	if err != nil {
		return nil, err
	}

	var metadata Metadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, err
	}

	artifactsBytes, err := os.ReadFile("dist/artifacts.json")
	if err != nil {
		return nil, err
	}

	var artifacts []Artifact
	err = json.Unmarshal(artifactsBytes, &artifacts)
	if err != nil {
		return nil, err
	}

	return &Result{Metadata: metadata, Artifacts: artifacts}, nil
}
