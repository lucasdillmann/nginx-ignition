package tailscale

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type Driver struct{}

func newDriver() *Driver {
	return &Driver{}
}

func (d Driver) ID() string {
	return "TAILSCALE"
}

func (d Driver) Name() string {
	return "Tailscale"
}

func (d Driver) ImportantInstructions() []string {
	return importantInstructions
}

func (d Driver) ConfigurationFields() []*dynamic_fields.DynamicField {
	return configurationFields
}

func (d Driver) Start(
	ctx context.Context,
	configDir string,
	destination vpn.Destination,
	parameters map[string]any,
) error {
	if state[destination.Name()] != nil {
		return nil
	}

	return d.doStart(ctx, configDir, destination, parameters)
}

func (d Driver) Reload(
	ctx context.Context,
	configDir string,
	destination vpn.Destination,
	parameters map[string]any,
) error {
	if state[destination.Name()] != nil {
		_ = d.Stop(ctx, destination)
	}

	return d.doStart(ctx, configDir, destination, parameters)
}

func (d Driver) Stop(ctx context.Context, destination vpn.Destination) error {
	endpoint := state[destination.Name()]
	if endpoint == nil {
		return nil
	}

	endpoint.Stop(ctx)
	delete(state, destination.Name())

	return nil
}

func (d Driver) doStart(
	ctx context.Context,
	configDir string,
	destination vpn.Destination,
	parameters map[string]any,
) error {
	authKey := parameters[authKeyFieldName].(string)

	var serverURL string
	if value, casted := parameters[coordinatorUrlFieldName].(string); casted {
		serverURL = value
	}

	state[destination.Name()] = &tailnetEndpoint{
		authKey:     authKey,
		configDir:   configDir,
		destination: destination,
		serverURL:   serverURL,
	}

	return state[destination.Name()].Start(ctx)
}
