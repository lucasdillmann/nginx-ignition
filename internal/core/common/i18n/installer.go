package i18n

import (
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
)

func Install() error {
	return container.Provide(newCommands)
}
