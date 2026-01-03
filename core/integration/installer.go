package integration

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
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
