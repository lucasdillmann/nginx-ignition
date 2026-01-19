package tailscale

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
	return "TAILSCALE"
}

func (d Driver) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.VpnTailscaleName)
}

func (d Driver) ImportantInstructions(ctx context.Context) []*i18n.Message {
	return importantInstructions(ctx)
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

	tEndpoint, ok := value.(*tailnetEndpoint)
	if !ok {
		return errors.New("invalid endpoint type in state")
	}

	tEndpoint.stop(ctx)

	return nil
}

func (d Driver) doStart(
	ctx context.Context,
	configDir string,
	endpoint vpn.Endpoint,
	parameters map[string]any,
) error {
	authKey, ok := parameters[authKeyFieldName].(string)
	if !ok || authKey == "" {
		return coreerror.New(i18n.M(ctx, i18n.K.VpnTailscaleAuthKeyRequired), true)
	}

	var serverURL string
	if value, casted := parameters[coordinatorURLFieldName].(string); casted {
		serverURL = value
	}

	tEndpoint := &tailnetEndpoint{
		authKey:   authKey,
		configDir: configDir,
		endpoint:  endpoint,
		serverURL: serverURL,
	}

	if err := tEndpoint.start(ctx); err != nil {
		return err
	}

	state.Store(endpoint.Hash(), tEndpoint)

	return nil
}
