package cfgfiles

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx/cfgfiles/provider"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
)

type Facade struct {
	hostCommands     host.Commands
	streamCommands   stream.Commands
	cacheCommands    cache.Commands
	settingsCommands settings.Commands
	providers        []provider.Provider
}

func newFacade(
	providers []provider.Provider,
	hostCommands host.Commands,
	streamCommands stream.Commands,
	cacheCommands cache.Commands,
	settingsCommands settings.Commands,
) *Facade {
	return &Facade{
		hostCommands:     hostCommands,
		streamCommands:   streamCommands,
		cacheCommands:    cacheCommands,
		settingsCommands: settingsCommands,
		providers:        providers,
	}
}

func (f *Facade) GetConfigurationFiles(
	ctx context.Context,
	paths *provider.Paths,
	supportedFeatures *provider.SupportedFeatures,
) (
	configFiles []provider.File,
	hosts []host.Host,
	streams []stream.Stream,
	err error,
) {
	enabledHosts, err := f.hostCommands.GetAllEnabled(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	enabledStreams, err := f.streamCommands.GetAllEnabled(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	enabledCaches, err := f.cacheCommands.GetAllInUse(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	cfg, err := f.settingsCommands.Get(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	providerCtx := &provider.Context{
		Context:           ctx,
		Paths:             paths,
		Hosts:             enabledHosts,
		Streams:           enabledStreams,
		Caches:            enabledCaches,
		SupportedFeatures: supportedFeatures,
		Cfg:               cfg,
	}

	configFiles = make([]provider.File, 0)
	for _, p := range f.providers {
		files, err := p.Provide(providerCtx)
		if err != nil {
			return nil, nil, nil, err
		}

		configFiles = append(configFiles, files...)
	}

	return configFiles, enabledHosts, enabledStreams, nil
}
