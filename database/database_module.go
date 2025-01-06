package database

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/user"
	"dillmann.com.br/nginx-ignition/database/access_list_repository"
	"dillmann.com.br/nginx-ignition/database/host_repository"
	"dillmann.com.br/nginx-ignition/database/user_repository"
	"go.uber.org/dig"
)

func InstallBeans(container *dig.Container) error {
	accessListRepository := access_list_repository.New(nil)
	hostRepository := host_repository.New(nil)
	userRepository := user_repository.New(nil)

	_ = container.Provide(func() *access_list.Repository { return &accessListRepository })
	_ = container.Provide(func() *host.Repository { return &hostRepository })
	_ = container.Provide(func() *user.Repository { return &userRepository })

	return nil
}
