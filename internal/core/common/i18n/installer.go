package i18n

import (
	"nginx-ignition/internal/core/common/container"
)

func Install() error {
	return container.Provide(newCommands)
}
