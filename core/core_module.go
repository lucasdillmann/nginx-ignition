package core

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"go.uber.org/dig"
)

func RegisterCoreBeans(container *dig.Container) error {
	if err := access_list.RegisterAccessListBeans(container); err != nil {
		return err
	}

	return nil
}
