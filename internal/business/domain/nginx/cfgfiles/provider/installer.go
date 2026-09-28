package provider

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/integration"
)

func Install() error {
	return container.Provide(newProviders)
}

func newProviders(
	accessListCommands accesslist.Commands,
	certificateCommands certificate.Commands,
	integrationCommands integration.Commands,
	cfg *configuration.Configuration,
) []Provider {
	return []Provider{
		newAccessListProvider(accessListCommands),
		newHostCertificateProvider(certificateCommands),
		newHostConfigurationProvider(integrationCommands),
		newHostRouteStaticResponseProvider(),
		newHostRouteSourceCodeProvider(),
		newMainConfigurationProvider(cfg),
		newMimeTypesProvider(),
		newStreamProvider(),
		newGeoIPProvider(cfg),
	}
}
