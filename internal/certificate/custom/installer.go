package custom

import (
	"nginx-ignition/internal/core/common/container"
)

func Install() error {
	return container.Provide(New)
}
