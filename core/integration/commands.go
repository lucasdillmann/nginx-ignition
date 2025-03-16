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

type ConfigureByIdCommand func(
	ctx context.Context,
	id string,
	enabled bool,
	parameters map[string]any,
) error

type GetByIdCommand func(ctx context.Context, id string) (*GetByIdOutput, error)

type GetOptionByIdCommand func(ctx context.Context, integrationId, optionId string) (*AdapterOption, error)

type GetOptionUrlByIdCommand func(ctx context.Context, integrationId, optionId string) (*string, error)

type ListOptionsCommand func(
	ctx context.Context,
	integrationId string,
	pageNumber, pageSize int,
	searchTerms *string,
) (*pagination.Page[*AdapterOption], error)

type ListCommand func(ctx context.Context) ([]*ListOutput, error)
