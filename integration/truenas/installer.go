package truenas

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(newDriver)
}
