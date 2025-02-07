package integration

import (
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type GetByIdOutput struct {
	ID                  string
	Name                string
	Description         string
	Enabled             bool
	ConfigurationFields []*dynamic_fields.DynamicField
	Parameters          map[string]any
}

type ListOutput struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
}

type ConfigureByIdCommand func(
	id string,
	enabled bool,
	parameters map[string]any,
) error

type GetByIdCommand func(id string) (*GetByIdOutput, error)

type GetOptionByIdCommand func(integrationId, optionId string) (*AdapterOption, error)

type GetOptionUrlByIdCommand func(integrationId, optionId string) (*string, error)

type ListOptionsCommand func(
	integrationId string,
	pageNumber, pageSize int,
	searchTerms *string,
) (*pagination.Page[*AdapterOption], error)

type ListCommand func() ([]*ListOutput, error)
