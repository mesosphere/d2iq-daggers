package containers

import (
	"dagger.io/dagger"

	"github.com/mesosphere/daggers/daggers"
)

// ApplyCustomizations applies customizations to the given container.
func ApplyCustomizations(
	runtime *daggers.Runtime, container *dagger.Container, customizers ...ContainerCustomizerFn,
) (*dagger.Container, error) {
	var err error

	for _, customizer := range customizers {
		container, err = customizer(runtime, container)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

// MountRuntimeWorkdir mounts the runtime workdir to the given container with the provided path and configures the
// working directory of the container to the hardcoded /src path.
func MountRuntimeWorkdir(runtime *daggers.Runtime, container *dagger.Container) *dagger.Container {
	return container.WithMountedDirectory("/src", runtime.Workdir).WithWorkdir("/src")
}
