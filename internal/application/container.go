package application

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/api"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/custom"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/external"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/letsencrypt"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/selfsigned"
	"github.com/lucasdillmann/nginx-ignition/internal/core"
	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/healthcheck"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/lifecycle"
	"github.com/lucasdillmann/nginx-ignition/internal/core/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/vpn"
	"github.com/lucasdillmann/nginx-ignition/internal/database"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/docker"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/truenas"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn/netbird"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn/tailscale"
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
