package docker

import (
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
)

type Adapter struct {
}

func newAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) ID() string {
	return "DOCKER"
}

func (a *Adapter) Name() string {
	return "Docker"
}

func (a *Adapter) Priority() int {
	return 1
}

func (a *Adapter) Description() string {
	return "Enables easy pick of a Docker container with ports exposing a service as a target for your nginx " +
		"ignition's host routes."
}

func (a *Adapter) ConfigurationFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{
		&connectionMode,
		&socketPath,
		&hostUrl,
		&proxyUrl,
	}
}

func (a *Adapter) GetAvailableOptions(
	parameters map[string]interface{},
	pageNumber, pageSize int,
	searchTerms *string,
) (*pagination.Page[*integration.AdapterOption], error) {
	// TODO: Implement this
	return nil, core_error.New("not implemented", false)
}

func (a *Adapter) GetAvailableOptionById(
	parameters map[string]interface{},
	id string,
) (*integration.AdapterOption, error) {
	// TODO: Implement this
	return nil, core_error.New("not implemented", false)
}

func (a *Adapter) GetOptionProxyUrl(
	parameters map[string]interface{},
	id string,
) (*string, error) {
	// TODO: Implement this
	return nil, core_error.New("not implemented", false)
}
