package binding

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(certificateCommands *certificate.Commands) *Commands {
	serviceInstance := newService(certificateCommands)

	return &Commands{
		Validate: serviceInstance.validateBinding,
	}
}
