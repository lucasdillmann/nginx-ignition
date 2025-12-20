package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/integration"
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
	return []*dynamicfields.DynamicField{
		&connectionModeField,
		&socketPathField,
		&hostUrlField,
		&swarmModeField,
		&swarmServiceMeshField,
		&swarmDNSResolverField,
		&proxyUrlField,
	}
}

func (a *Driver) GetAvailableOptions(
	ctx context.Context,
	parameters map[string]any,
	_, _ int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[*integration.DriverOption], error) {
	optionResolver, err := startOptionResolver(parameters)
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
	optionResolver, err := startOptionResolver(parameters)
	if err != nil {
		return nil, err
	}

	option, err := optionResolver.ResolveOptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return option.DriverOption, nil
}

func (a *Driver) GetOptionProxyURL(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*string, *[]string, error) {
	optionResolver, err := startOptionResolver(parameters)
	if err != nil {
		return nil, nil, err
	}

	option, err := optionResolver.ResolveOptionByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	publicUrl, _ := parameters[proxyUrlField.ID].(string)

	url, err := option.URL(ctx, publicUrl)
	return url, option.DNSResolvers, err
}

func startOptionResolver(parameters map[string]any) (resolver.Resolver, error) {
	var connectionString string
	switch parameters[connectionModeField.ID].(string) {
	case "SOCKET":
		socketPath := parameters[socketPathField.ID].(string)
		connectionString = "unix://" + socketPath
	case "TCP":
		hostUrl := parameters[hostUrlField.ID].(string)
		connectionString = hostUrl
	default:
		return nil, coreerror.New("Invalid connection mode", false)
	}

	dockerClient, err := client.NewClientWithOpts(
		client.WithHost(connectionString),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	swarmEnabled, useServiceMesh, dnsResolvers := extractSwarmParams(parameters)
	if swarmEnabled {
		return resolver.FromServices(dockerClient, useServiceMesh, dnsResolvers), nil
	}

	return resolver.FromContainers(dockerClient), nil
}

func extractSwarmParams(parameters map[string]any) (bool, bool, *[]string) {
	swarmMode := false
	useServiceMesh := false
	var dnsResolvers *[]string

	if rawValue, exists := parameters[swarmModeField.ID]; exists {
		swarmMode = rawValue.(bool)
	}

	if rawValue, exists := parameters[swarmServiceMeshField.ID]; exists {
		useServiceMesh = rawValue.(bool)
	}

	if rawValue, exists := parameters[swarmDNSResolverField.ID]; exists {
		textValue := rawValue.(string)
		dnsResolvers = ptr.Of(strings.Split(textValue, "\n"))

		for index, value := range *dnsResolvers {
			(*dnsResolvers)[index] = strings.TrimSpace(value)
		}
	}

	return swarmMode, useServiceMesh, dnsResolvers
}
