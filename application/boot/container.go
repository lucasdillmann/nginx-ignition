package boot

import (
	"dillmann.com.br/nginx-ignition/api"
	"dillmann.com.br/nginx-ignition/core"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/database"
	"go.uber.org/dig"
)

func startContainer() (*dig.Container, error) {
	container := dig.New()

	if err := installModules(container); err != nil {
		return nil, err
	}

	return container, nil
}

func installModules(container *dig.Container) error {
	if err := container.Provide(configuration.New); err != nil {
		return err
	}

	if err := container.Provide(lifecycle.New); err != nil {
		return err
	}

	if err := database.Install(container); err != nil {
		return err
	}

	if err := core.Install(container); err != nil {
		return err
	}

	if err := api.Install(container); err != nil {
		return err
	}

	return nil
}
