package access_list

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"go.uber.org/dig"
)

func RegisterAccessListBeans(container *dig.Container) error {
	return container.Invoke(func(accessListRepository *AccessListRepository, hostRepository *host.HostRepository) {
		service := &accessListService{
			accessListRepository,
			hostRepository,
		}

		_ = container.Provide(func() GetAccessListByIdCommand { return service.findById })
		_ = container.Provide(func() DeleteAccessListByIdCommand { return service.deleteById })
		_ = container.Provide(func() ListAccessListCommand { return service.list })
		_ = container.Provide(func() SaveAccessListCommand { return service.save })
	})
}
