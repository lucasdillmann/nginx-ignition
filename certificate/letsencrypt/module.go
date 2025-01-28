package letsencrypt

import "go.uber.org/dig"

func Install(container *dig.Container) error {
	return container.Provide(New)
}
