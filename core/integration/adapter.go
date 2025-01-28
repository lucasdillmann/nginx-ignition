package integration

import (
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Adapter interface {
	ID() string
	Name() string
	Priority() int
	Description() string
	ConfigurationFields() []*dynamic_fields.DynamicField
	GetAvailableOptions(
		parameters map[string]interface{},
		pageNumber, pageSize int,
		searchTerms *string,
	) (*pagination.Page[*AdapterOption], error)
	GetAvailableOptionById(
		parameters map[string]interface{},
		id string,
	) (*AdapterOption, error)
	GetOptionProxyUrl(
		parameters map[string]interface{},
		id string,
	) (*string, error)
}

type AdapterOption struct {
	ID   string
	Name string
}
