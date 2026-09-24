package integration

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
)

func Install() error {
	return container.Provide(newCommands)
}

func newCommands(
	repository Repository,
) Commands {
	drivers := func() []Driver {
		return container.Get[[]Driver]()
	}

	return newService(repository, drivers)
}
