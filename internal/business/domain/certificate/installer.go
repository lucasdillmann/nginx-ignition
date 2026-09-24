package certificate

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
)

func Install() error {
	if err := container.Provide(newCommands); err != nil {
		return err
	}

	return container.Run(registerScheduledTask)
}

func newCommands(repository Repository) (Commands, *service) {
	providers := func() []Provider {
		return container.Get[[]Provider]()
	}

	serviceInstance := newService(repository, providers)
	return serviceInstance, serviceInstance
}
