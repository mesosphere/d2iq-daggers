package asdf

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Install installs all plugins and versions specified in local and global .tool-versions file.
func Install() error {
	plugins, err := ListPlugins()
	if err != nil {
		return err
	}

	tools, err := ParseToolVersions()
	if err != nil {
		return err
	}

	for plugin := range tools {
		ensurePluginExist(plugins, plugin)
	}

	// Install all plugins and versions
	return sh.RunV("asdf", "install")
}

// InstallPlugins installs given plugins and version specified in local .tool-versions file or latest version
// if not specified.
func InstallPlugins(pluginsToInstall ...string) error {
	plugins, err := ListPlugins()
	if err != nil {
		return err
	}

	tools, err := ParseToolVersions()
	if err != nil {
		return err
	}

	for _, plugin := range pluginsToInstall {
		version, ok := tools[plugin]
		if !ok {
			return fmt.Errorf("plugin %s is not specified in .tool-versions", plugin)
		}

		ensurePluginExist(plugins, plugin)

		err = installPackage(plugin, version)
		if err != nil {
			return err
		}
	}

	return nil
}

// Upgrade upgrades all plugins and versions specified in local .tool-versions file.
//
//nolint:revive // Disable cognitive-complexity check. There is not enough gain from reducing cognitive-complexity,
func Upgrade() error {
	plugins, err := ListPlugins()
	if err != nil {
		return err
	}

	tools, err := ParseToolVersions()
	if err != nil {
		return err
	}

	for plugin, version := range tools {
		ensurePluginExist(plugins, plugin)

		// If plugin is frozen, we'll not upgrade it
		if version.VersionFreeze {
			fmt.Printf("plugin %s is frozen, with version %s, skipping upgrade\n", plugin, version.Version)
			continue
		}

		latest, err := getLatestVersion(plugin)
		if err != nil {
			return err
		}

		// if latest version is the same as the one specified in .tool-versions, no need to do anything
		if latest == version.Version {
			fmt.Printf("plugin %s is already up to date with version %s\n", plugin, version.Version)

			continue
		}

		fmt.Printf("upgrading %s from %s to %s\n", plugin, version.Version, latest)

		err = sh.RunV("asdf", "install", plugin, latest)
		if err != nil {
			return fmt.Errorf("failed to upgrade %s: %w", plugin, err)
		}

		err = sh.RunV("asdf", "local", plugin, latest)
		if err != nil {
			return fmt.Errorf("failed to set %s to %s: %w", plugin, latest, err)
		}
	}

	return nil
}

func getLatestVersion(plugin string) (string, error) {
	allVersions, err := ListPluginVersions(plugin)
	if err != nil {
		return "", fmt.Errorf("failed to list versions for plugin %s: %w", plugin, err)
	}

	// last of slice would be latest version singe versions are asc order
	return allVersions[len(allVersions)-1], nil
}

func ensurePluginExist(plugins Plugins, plugin string) {
	if !plugins[plugin] {
		fmt.Printf("Add plugin %s\n", plugin)
		// ignore error if plugin is already installed.
		// TODO: improve this check
		_ = sh.RunV("asdf", "plugin", "add", plugin)
	}
}

func installPackage(plugin string, version Version) error {
	fmt.Printf("Installing plugin %s with version %s\n", plugin, version.Version)

	err := sh.RunV("asdf", "install", plugin, version.Version)
	if err != nil {
		return fmt.Errorf("failed to install %s: %w", plugin, err)
	}

	return nil
}
