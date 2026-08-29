package application

import (
	"context"

	"nginx-ignition/internal/api"
	"nginx-ignition/internal/certificate/custom"
	"nginx-ignition/internal/certificate/external"
	"nginx-ignition/internal/certificate/letsencrypt"
	"nginx-ignition/internal/certificate/selfsigned"
	"nginx-ignition/internal/core"
	"nginx-ignition/internal/core/certificate"
	"nginx-ignition/internal/core/common/configuration"
	"nginx-ignition/internal/core/common/container"
	"nginx-ignition/internal/core/common/healthcheck"
	"nginx-ignition/internal/core/common/i18n"
	"nginx-ignition/internal/core/common/lifecycle"
	"nginx-ignition/internal/core/integration"
	"nginx-ignition/internal/core/vpn"
	"nginx-ignition/internal/database"
	"nginx-ignition/internal/integration/docker"
	"nginx-ignition/internal/integration/truenas"
	"nginx-ignition/internal/vpn/netbird"
	"nginx-ignition/internal/vpn/tailscale"
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
		i18n.Install,
		database.Install,
		core.Install,
		api.Install,
		letsencrypt.Install,
		selfsigned.Install,
		custom.Install,
		external.Install,
		docker.Install,
		truenas.Install,
		tailscale.Install,
		netbird.Install,
		installCertificateDriverAggregation,
		installIntegrationDriverAggregation,
		installVpnDriverAggregation,
	)
}

func installCertificateDriverAggregation(
	acmeCertificateProvider *letsencrypt.Provider,
	customCertificateProvider *custom.Provider,
	selfSignedCertificateProvider *selfsigned.Provider,
	externalCertificateProvider *external.Provider,
) error {
	return container.Singleton([]certificate.Provider{
		acmeCertificateProvider,
		customCertificateProvider,
		selfSignedCertificateProvider,
		externalCertificateProvider,
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
	ts *tailscale.Driver,
	nb *netbird.Driver,
) error {
	return container.Singleton([]vpn.Driver{
		ts,
		nb,
	})
}
