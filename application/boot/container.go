package boot

import (
	"context"

	"go.uber.org/dig"

	"dillmann.com.br/nginx-ignition/api"
	"dillmann.com.br/nginx-ignition/certificate/custom"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt"
	"dillmann.com.br/nginx-ignition/certificate/selfsigned"
	"dillmann.com.br/nginx-ignition/core"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/database"
	"dillmann.com.br/nginx-ignition/integration/docker"
	"dillmann.com.br/nginx-ignition/integration/truenas"
)

func startContainer(ctx context.Context) (*dig.Container, error) {
	container := dig.New()
	_ = container.Provide(func() *dig.Container {
		return container
	})

	_ = container.Provide(func() context.Context {
		return ctx
	})

	if err := installModules(container); err != nil {
		return nil, err
	}

	return container, nil
}

func installModules(container *dig.Container) error {
	if err := container.Provide(configuration.New); err != nil {
		return err
	}

	if err := container.Provide(lifecycle.New); err != nil {
		return err
	}

	if err := database.Install(container); err != nil {
		return err
	}

	if err := core.Install(container); err != nil {
		return err
	}

	if err := api.Install(container); err != nil {
		return err
	}

	if err := letsencrypt.Install(container); err != nil {
		return err
	}

	if err := selfsigned.Install(container); err != nil {
		return err
	}

	if err := custom.Install(container); err != nil {
		return err
	}

	if err := docker.Install(container); err != nil {
		return err
	}

	if err := truenas.Install(container); err != nil {
		return err
	}

	if err := container.Invoke(installCertificateProviderAggregation); err != nil {
		return err
	}

	return container.Invoke(installIntegrationAdapterAggregation)
}

func installCertificateProviderAggregation(
	container *dig.Container,
	acmeCertificateProvider *letsencrypt.Provider,
	customCertificateProvider *custom.Provider,
	selfSignedCertificateProvider *selfsigned.Provider,
) error {
	return container.Provide(func() []certificate.Provider {
		return []certificate.Provider{
			acmeCertificateProvider,
			customCertificateProvider,
			selfSignedCertificateProvider,
		}
	})
}

func installIntegrationAdapterAggregation(
	container *dig.Container,
	dockerAdapter *docker.Adapter,
	trueNasAdapter *truenas.Adapter,
) error {
	return container.Provide(func() []integration.Driver {
		return []integration.Driver{
			dockerAdapter,
			trueNasAdapter,
		}
	})
}
