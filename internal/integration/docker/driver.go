package docker

import (
	"context"

	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/i18n"
	"nginx-ignition/internal/core/common/pagination"
	"nginx-ignition/internal/core/integration"
	"nginx-ignition/internal/integration/docker/fields"
	"nginx-ignition/internal/integration/docker/resolver"
)

type Driver struct{}

func newDriver() *Driver {
	return &Driver{}
}

func (a *Driver) ID() string {
	return "DOCKER"
}

func (a *Driver) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CommonDocker)
}

func (a *Driver) Description(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.IntegrationDockerDescription)
}

func (a *Driver) ConfigurationFields(ctx context.Context) []dynamicfields.DynamicField {
	return fields.DynamicFields(ctx)
}

func (a *Driver) GetAvailableOptions(
	ctx context.Context,
	parameters map[string]any,
	_, _ int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[integration.DriverOption], error) {
	optionResolver, err := resolver.For(ctx, parameters)
	if err != nil {
		return nil, err
	}

	options, err := optionResolver.ResolveOptions(ctx, tcpOnly, searchTerms)
	if err != nil {
		return nil, err
	}

	totalItems := len(options)
	driverOptions := make([]integration.DriverOption, totalItems)
	for index, option := range options {
		driverOptions[index] = option.DriverOption
	}

	return pagination.New(0, totalItems, totalItems, driverOptions), nil
}

func (a *Driver) GetAvailableOptionByID(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*integration.DriverOption, error) {
	optionResolver, err := resolver.For(ctx, parameters)
	if err != nil {
		return nil, err
	}

	option, err := optionResolver.ResolveOptionByID(ctx, id)
	if err != nil || option == nil {
		return nil, err
	}

	return &option.DriverOption, nil
}

func (a *Driver) GetOptionProxyURL(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*string, []string, error) {
	optionResolver, err := resolver.For(ctx, parameters)
	if err != nil {
		return nil, nil, err
	}

	option, err := optionResolver.ResolveOptionByID(ctx, id)
	if err != nil || option == nil {
		return nil, nil, err
	}

	return option.URL(ctx)
}
