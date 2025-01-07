package main

import (
	"dillmann.com.br/nginx-ignition/application/configuration_provider"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"go.uber.org/dig"
)

func installBeans(container *dig.Container) error {
	configurationInstance := configuration_provider.New()

	_ = container.Provide(func() *configuration.Configuration { return &configurationInstance })
	return nil
}
