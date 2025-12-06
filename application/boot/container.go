package boot

import (
	"context"

	"dillmann.com.br/nginx-ignition/api"
	"dillmann.com.br/nginx-ignition/certificate/custom"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt"
	"dillmann.com.br/nginx-ignition/certificate/selfsigned"
	"dillmann.com.br/nginx-ignition/core"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
	"dillmann.com.br/nginx-ignition/database"
	"dillmann.com.br/nginx-ignition/integration/docker"
	"dillmann.com.br/nginx-ignition/integration/truenas"
	"dillmann.com.br/nginx-ignition/vpn/tailscale"
)

func startContainer(ctx context.Context) error {
	container.Init(ctx)

	if err := container.Provide(
		configuration.New,
		lifecycle.New,
		healthcheck.New,
	); err != nil {
		return err
	}

	return container.Run(
		database.Install,
		core.Install,
		api.Install,
		letsencrypt.Install,
		selfsigned.Install,
		custom.Install,
		docker.Install,
		truenas.Install,
		tailscale.Install,
		installCertificateDriverAggregation,
		installIntegrationDriverAggregation,
		installVpnDriverAggregation,
	)
}

func installCertificateDriverAggregation(
	acmeCertificateProvider *letsencrypt.Provider,
	customCertificateProvider *custom.Provider,
	selfSignedCertificateProvider *selfsigned.Provider,
) error {
	return container.Singleton([]certificate.Provider{
		acmeCertificateProvider,
		customCertificateProvider,
		selfSignedCertificateProvider,
	})
}

func installIntegrationDriverAggregation(
	dockerAdapter *docker.Driver,
	trueNasAdapter *truenas.Driver,
) error {
	return container.Singleton([]integration.Driver{
		dockerAdapter,
		trueNasAdapter,
	})
}

func installVpnDriverAggregation(
	tailscale *tailscale.Driver,
) error {
	return container.Singleton([]vpn.Driver{
		tailscale,
	})
}
