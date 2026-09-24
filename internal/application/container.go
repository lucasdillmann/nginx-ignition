package application

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/api"
	"github.com/lucasdillmann/nginx-ignition/internal/business"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/healthcheck"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/lifecycle"
	businesscertificate "github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	businessintegration "github.com/lucasdillmann/nginx-ignition/internal/business/domain/integration"
	businessvpn "github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/custom"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/external"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/selfsigned"
	"github.com/lucasdillmann/nginx-ignition/internal/database"
	"github.com/lucasdillmann/nginx-ignition/internal/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/domain/docker"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/domain/truenas"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn/domain/netbird"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn/domain/tailscale"
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
		business.Install,
		api.Install,
		certificate.Install,
		integration.Install,
		vpn.Install,
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
	return container.Singleton([]businesscertificate.Provider{
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
	return container.Singleton([]businessintegration.Driver{
		dockerAdapter,
		trueNasAdapter,
	})
}

func installVpnDriverAggregation(
	ts *tailscale.Driver,
	nb *netbird.Driver,
) error {
	return container.Singleton([]businessvpn.Driver{
		ts,
		nb,
	})
}
