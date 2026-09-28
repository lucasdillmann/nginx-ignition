package cfgfiles

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx/cfgfiles/provider"
)

func Install() error {
	if err := container.Run(provider.Install); err != nil {
		return err
	}

	return container.Provide(newFacade)
}
