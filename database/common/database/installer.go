package database

import "go.uber.org/dig"

func Install(container *dig.Container) error {
	if err := container.Provide(newDatabase); err != nil {
		return err
	}

	if err := container.Invoke(registerStartup); err != nil {
		return err
	}

	if err := container.Invoke(registerShutdown); err != nil {
		return err
	}

	return nil
}
