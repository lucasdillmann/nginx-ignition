package broadcast

import (
	"nginx-ignition/internal/core/common/container"
)

func Install() error {
	return container.Run(registerShutdown)
}
