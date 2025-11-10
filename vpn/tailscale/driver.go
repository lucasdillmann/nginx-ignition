package tailscale

import (
	"context"
	"fmt"

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
	if _, exists := state.Load(destination.Hash()); exists {
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
	if _, exists := state.Load(destination.Hash()); exists {
		_ = d.Stop(ctx, destination)
	}

	return d.doStart(ctx, configDir, destination, parameters)
}

func (d Driver) Stop(ctx context.Context, destination vpn.Destination) error {
	value, exists := state.LoadAndDelete(destination.Hash())
	if !exists {
		return nil
	}

	endpoint, ok := value.(*tailnetEndpoint)
	if !ok {
		return fmt.Errorf("invalid endpoint type in state")
	}

	endpoint.Stop(ctx)

	return nil
}

func (d Driver) doStart(
	ctx context.Context,
	configDir string,
	destination vpn.Destination,
	parameters map[string]any,
) error {
	authKey, ok := parameters[authKeyFieldName].(string)
	if !ok || authKey == "" {
		return fmt.Errorf("authKey parameter is required and must be a non-empty string")
	}

	var serverURL string
	if value, casted := parameters[coordinatorUrlFieldName].(string); casted {
		serverURL = value
	}

	endpoint := &tailnetEndpoint{
		authKey:     authKey,
		configDir:   configDir,
		destination: destination,
		serverURL:   serverURL,
	}

	if err := endpoint.Start(ctx); err != nil {
		return err
	}

	state.Store(destination.Hash(), endpoint)

	return nil
}
