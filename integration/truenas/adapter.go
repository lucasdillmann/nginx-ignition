package truenas

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
	return "TRUENAS_SCALE"
}

func (a *Adapter) Name() string {
	return "TrueNAS Scale"
}

func (a *Adapter) Priority() int {
	return 2
}

func (a *Adapter) Description() string {
	return "TrueNAS allows, alongside many other things, to run your favorite apps under Docker containers. With this " +
		"integration enabled, you will be able to easily pick any app exposing a service in your TrueNAS as a " +
		"target for your nginx ignition's host routes."
}

func (a *Adapter) ConfigurationFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{
		&url,
		&proxyUrl,
		&username,
		&password,
	}
}

func (a *Adapter) GetAvailableOptions(
	parameters map[string]interface{},
	pageNumber, pageSize int,
	searchTerms *string,
) (*pagination.Page[*integration.AdapterOption], error) {
	return nil, core_error.New("not implemented", false)
}

func (a *Adapter) GetAvailableOptionById(
	parameters map[string]interface{},
	id string,
) (*integration.AdapterOption, error) {
	return nil, core_error.New("not implemented", false)
}

func (a *Adapter) GetOptionProxyUrl(
	parameters map[string]interface{},
	id string,
) (*string, error) {
	return nil, core_error.New("not implemented", false)
}
