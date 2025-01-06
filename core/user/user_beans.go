package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"go.uber.org/dig"
)

func InstallBeans(container *dig.Container) error {
	return container.Invoke(func(repository *Repository, configuration *configuration.Configuration) {
		serviceInstance := &service{repository, configuration}

		_ = container.Provide(func() AuthenticateCommand { return serviceInstance.authenticate })
		_ = container.Provide(func() DeleteByIdCommand { return serviceInstance.deleteById })
		_ = container.Provide(func() GetByIdCommand { return serviceInstance.getById })
		_ = container.Provide(func() GetCountCommand { return serviceInstance.count })
		_ = container.Provide(func() GetStatusCommand { return serviceInstance.isEnabled })
		_ = container.Provide(func() ListCommand { return serviceInstance.list })
		_ = container.Provide(func() SaveCommand { return serviceInstance.save })
		_ = container.Provide(func() ChangePasswordCommand { return serviceInstance.changePassword })
	})
}
