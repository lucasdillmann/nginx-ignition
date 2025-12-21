package docker

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/integration/docker/fields"
	"dillmann.com.br/nginx-ignition/integration/docker/resolver"
)

type Driver struct{}

func newDriver() *Driver {
	return &Driver{}
}

func (a *Driver) ID() string {
	return "DOCKER"
}

func (a *Driver) Name() string {
	return "Docker"
}

func (a *Driver) Description() string {
	return "Enables easy pick of a Docker container with ports exposing a service as a target for your nginx " +
		"ignition's host routes."
}

func (a *Driver) ConfigurationFields() []*dynamicfields.DynamicField {
	return fields.All
}

func (a *Driver) GetAvailableOptions(
	ctx context.Context,
	parameters map[string]any,
	_, _ int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[*integration.DriverOption], error) {
	optionResolver, err := resolver.For(parameters)
	if err != nil {
		return nil, err
	}

	options, err := optionResolver.ResolveOptions(ctx, tcpOnly, searchTerms)
	if err != nil {
		return nil, err
	}

	totalItems := len(*options)
	driverOptions := make([]*integration.DriverOption, totalItems)
	for index, option := range *options {
		driverOptions[index] = option.DriverOption
	}

	return pagination.New(0, totalItems, totalItems, driverOptions), nil
}

func (a *Driver) GetAvailableOptionById(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*integration.DriverOption, error) {
	optionResolver, err := resolver.For(parameters)
	if err != nil {
		return nil, err
	}

	option, err := optionResolver.ResolveOptionByID(ctx, id)
	if err != nil || option == nil {
		return nil, err
	}

	return option.DriverOption, nil
}

func (a *Driver) GetOptionProxyURL(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*string, *[]string, error) {
	optionResolver, err := resolver.For(parameters)
	if err != nil {
		return nil, nil, err
	}

	option, err := optionResolver.ResolveOptionByID(ctx, id)
	if err != nil || option == nil {
		return nil, nil, err
	}

	return option.URL(ctx)
}
