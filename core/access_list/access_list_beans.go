package access_list

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"go.uber.org/dig"
)

func InstallBeans(container *dig.Container) error {
	return container.Invoke(func(accessListRepository *Repository, hostRepository *host.Repository) {
		serviceInstance := &service{
			accessListRepository,
			hostRepository,
		}

		_ = container.Provide(func() GetByIdCommand { return serviceInstance.findById })
		_ = container.Provide(func() DeleteById { return serviceInstance.deleteById })
		_ = container.Provide(func() ListCommand { return serviceInstance.list })
		_ = container.Provide(func() SaveCommand { return serviceInstance.save })
	})
}
