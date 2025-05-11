package integration

import (
	"context"
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

type Commands struct {
	GetById          func(ctx context.Context, id string) (*GetByIdOutput, error)
	GetOptionById    func(ctx context.Context, integrationId, optionId string) (*AdapterOption, error)
	GetOptionUrlById func(ctx context.Context, integrationId, optionId string) (*string, error)
	List             func(ctx context.Context) ([]*ListOutput, error)
	ConfigureById    func(
		ctx context.Context,
		id string,
		enabled bool,
		parameters map[string]any,
	) error
	ListOptions func(
		ctx context.Context,
		integrationId string,
		pageNumber, pageSize int,
		searchTerms *string,
		tcpOnly bool,
	) (*pagination.Page[*AdapterOption], error)
}
