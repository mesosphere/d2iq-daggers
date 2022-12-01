package asdf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
)

// Plugins is a simple custom type to hold plugin name to quickly check if a plugin is already added.
type Plugins map[string]bool

// Version is a simple struct to holds plugin version information.
type Version struct {
	// Plugin version
	Version string

	// If true, plugin version should be ignored during upgrade
	VersionFreeze bool
}

// PluginVersions is a map of plugin name and version.
type PluginVersions map[string]Version

// ListPlugins lists all plugins installed in asdf.
func ListPlugins() (Plugins, error) {
	plugins := make(Plugins)

	// List all plugins
	output, err := sh.Output("asdf", "plugin", "list")
	if err != nil {
		return nil, err
	}

	plugin := bufio.NewScanner(strings.NewReader(output))

	for plugin.Scan() {
		plugins[strings.TrimSpace(plugin.Text())] = true
	}

	return plugins, nil
}

// ListPluginVersions lists all installed versions of software for the specified plugin in ascending order.
func ListPluginVersions(plugin string) ([]string, error) {
	output, err := sh.Output("asdf", "list", "all", plugin)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0)

	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		versions = append(versions, strings.TrimSpace(scanner.Text()))
	}

	return versions, nil
}

// ParseToolVersions parses .tool-versions file and returns a map of plugin name and version.
func ParseToolVersions() (PluginVersions, error) {
	plugins := make(PluginVersions)

	// Check if .tool-versions file exists, if not, return empty map
	if _, err := os.Stat(".tool-versions"); os.IsNotExist(err) {
		fmt.Println("no .tool-versions file found in current directory, skipping")
		return plugins, nil
	}

	tools, err := os.Open(".tool-versions")
	if err != nil {
		return nil, fmt.Errorf("failed to read .tools-versions: %w", err)
	}

	traverseNonCommentLines(tools, func(line string) {
		// Split the line into tokens. Normally, token array should have 2 elements. However, this file can have
		// comments which are prefixed with a #. We're using these comments to freeze versions of plugins.
		// If a plugin is frozen,
		// we'll not try to upgrade it.
		tokens := strings.Split(line, " ")

		version := Version{
			Version:       tokens[1],
			VersionFreeze: false,
		}

		// If next two tokens are # and FREEZE, then we'll freeze the version of the plugin
		if len(tokens) > 3 {
			version.VersionFreeze = tokens[2] == "#" && tokens[3] == "FREEZE"
		}

		plugins[tokens[0]] = version
	})

	return plugins, nil
}

func traverseNonCommentLines(r io.Reader, visit func(line string)) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		// Trim whitespace
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		visit(line)
	}
}
