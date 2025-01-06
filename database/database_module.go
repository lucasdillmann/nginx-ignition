package database

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/host"
	access_list_repository "dillmann.com.br/nginx-ignition/database/access_list"
	host_repository "dillmann.com.br/nginx-ignition/database/host"
	"go.uber.org/dig"
)

func RegisterDatabaseBeans(container *dig.Container) error {
	accessListRepository := access_list_repository.New(nil)
	hostRepository := host_repository.New(nil)

	_ = container.Provide(func() *access_list.AccessListRepository { return &accessListRepository })
	_ = container.Provide(func() *host.HostRepository { return &hostRepository })

	return nil
}
