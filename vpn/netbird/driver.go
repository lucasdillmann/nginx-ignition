package netbird

import (
	"context"
	"errors"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type Driver struct{}

func newDriver() *Driver {
	return &Driver{}
}

func (d Driver) ID() string {
	return "NETBIRD"
}

func (d Driver) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.VpnNetbirdName)
}

func (d Driver) ImportantInstructions(ctx context.Context) []*i18n.Message {
	return importantInstructions(ctx)
}

func (d Driver) EndpointSSLSupport(_ context.Context) vpn.EndpointSSLSupport {
	return vpn.DriverManagedEndpointSSLSupport
}

func (d Driver) ConfigurationFields(ctx context.Context) []dynamicfields.DynamicField {
	return configurationFields(ctx)
}

func (d Driver) Start(
	ctx context.Context,
	configDir string,
	endpoint vpn.Endpoint,
	parameters map[string]any,
) error {
	if _, exists := state.Load(endpoint.Hash()); exists {
		return nil
	}

	return d.doStart(ctx, configDir, endpoint, parameters)
}

func (d Driver) Reload(
	ctx context.Context,
	configDir string,
	endpoint vpn.Endpoint,
	parameters map[string]any,
) error {
	if _, exists := state.Load(endpoint.Hash()); exists {
		_ = d.Stop(ctx, endpoint)
	}

	return d.doStart(ctx, configDir, endpoint, parameters)
}

func (d Driver) Stop(ctx context.Context, endpoint vpn.Endpoint) error {
	value, exists := state.LoadAndDelete(endpoint.Hash())
	if !exists {
		return nil
	}

	nbEndpoint, ok := value.(*netbirdEndpoint)
	if !ok {
		return errors.New("invalid endpoint type in state")
	}

	nbEndpoint.stop(ctx)

	return nil
}

func (d Driver) doStart(
	ctx context.Context,
	configDir string,
	endpoint vpn.Endpoint,
	parameters map[string]any,
) error {
	setupKey, ok := parameters[setupKeyFieldName].(string)
	if !ok || setupKey == "" {
		return coreerror.New(i18n.M(ctx, i18n.K.VpnNetbirdSetupKeyRequired), true)
	}

	var managementURL string
	if value, casted := parameters[managementURLFieldName].(string); casted {
		managementURL = value
	}

	nbEndpoint := &netbirdEndpoint{
		setupKey:      setupKey,
		configDir:     configDir,
		endpoint:      endpoint,
		managementURL: managementURL,
	}

	if err := nbEndpoint.start(ctx); err != nil {
		return err
	}

	state.Store(endpoint.Hash(), nbEndpoint)

	return nil
}
