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

func (d Driver) ConfigurationFields() []*dynamic_fields.DynamicField {
	return configurationFields
}

func (d Driver) Start(
	ctx context.Context,
	name, configDir string,
	destination *vpn.Destination,
	parameters map[string]any,
) error {
	if state[name] != nil {
		_ = d.Stop(ctx, name)
	}

	authKey := parameters[authKeyFieldName].(string)

	var serverURL string
	if value, casted := parameters[coordinatorUrlFieldName].(string); casted {
		serverURL = value
	}

	state[name] = &tailnetEndpoint{
		name:        name,
		authKey:     authKey,
		configDir:   configDir,
		destination: destination,
		serverURL:   serverURL,
	}

	return state[name].Start(ctx)
}

func (d Driver) Stop(ctx context.Context, name string) error {
	endpoint := state[name]
	if endpoint == nil {
		return nil
	}

	endpoint.Stop(ctx)
	delete(state, name)

	return nil
}
