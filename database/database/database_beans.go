package database

import "go.uber.org/dig"

func Install(container *dig.Container) error {
	if err := container.Provide(New); err != nil {
		return err
	}

	if err := container.Invoke(Register); err != nil {
		return err
	}

	return nil
}
