package migrations

import "go.uber.org/dig"

func Install(container *dig.Container) error {
	if err := container.Provide(newMigrations); err != nil {
		return err
	}

	if err := container.Invoke(registerStartup); err != nil {
		return err
	}

	return nil
}
