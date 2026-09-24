package broadcast

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
)

func Install() error {
	return container.Run(registerShutdown)
}
