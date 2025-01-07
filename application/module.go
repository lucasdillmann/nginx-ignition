package main

import (
	"dillmann.com.br/nginx-ignition/application/configuration_provider"
	"go.uber.org/dig"
)

func install(container *dig.Container) error {
	return container.Provide(configuration_provider.New)
}
