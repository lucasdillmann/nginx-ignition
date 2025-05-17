package docker

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"net/url"
	"strings"
)

type containerMetadata struct {
	name      string
	container *container.Summary
	port      *container.Port
}

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
		&connectionModeField,
		&socketPathField,
		&hostUrlField,
		&proxyUrlField,
	}
}

func (a *Adapter) GetAvailableOptions(
	ctx context.Context,
	parameters map[string]any,
	_, _ int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[*integration.AdapterOption], error) {
	options, err := a.resolveAvailableOptions(ctx, parameters, searchTerms, tcpOnly)
	if err != nil {
		return nil, err
	}

	totalItems := len(options)
	adapterOptions := make([]*integration.AdapterOption, totalItems)
	for i, option := range options {
		adapterOptions[i] = toAdapterOption(option)
	}

	return pagination.New(0, totalItems, totalItems, adapterOptions), nil
}

func (a *Adapter) GetAvailableOptionById(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*integration.AdapterOption, error) {
	option, err := a.resolveAvailableOptionById(ctx, parameters, id)
	if err != nil || option == nil {
		return nil, err
	}

	return toAdapterOption(option), nil
}

func (a *Adapter) GetOptionProxyUrl(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*string, error) {
	option, err := a.resolveAvailableOptionById(ctx, parameters, id)
	if err != nil || option == nil {
		return nil, err
	}

	publicUrl, _ := parameters[proxyUrlField.ID].(string)
	var targetHost string
	if publicUrl != "" {
		uri, err := url.Parse(publicUrl)
		if err != nil {
			return nil, err
		}

		targetHost = uri.Hostname()
	} else {
		if len(option.container.NetworkSettings.Networks) > 0 {
			for _, network := range option.container.NetworkSettings.Networks {
				targetHost = network.IPAddress
				break
			}
		}

		if targetHost == "" {
			return nil, fmt.Errorf("no network or IP address found for the container with ID %s", id)
		}
	}

	targetPort := option.port.PrivatePort
	if publicUrl != "" {
		targetPort = option.port.PublicPort
	}

	result := fmt.Sprintf("http://%s:%d", targetHost, targetPort)
	return &result, nil
}

func (a *Adapter) resolveAvailableOptionById(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*containerMetadata, error) {
	options, err := a.resolveAvailableOptions(ctx, parameters, nil, false)
	if err != nil {
		return nil, err
	}

	idParts := strings.Split(id, ":")
	if len(idParts) != 2 {
		return nil, core_error.New("Invalid option ID", true)
	}

	containerId := idParts[0]
	portNumber := idParts[1]

	for _, option := range options {
		if option.container.ID == containerId && fmt.Sprintf("%d", option.port.PublicPort) == portNumber {
			return option, nil
		}
	}

	return nil, nil
}

func (a *Adapter) resolveAvailableOptions(
	ctx context.Context,
	parameters map[string]any,
	searchTerms *string,
	tcpOnly bool,
) ([]*containerMetadata, error) {
	dockerClient, err := startClient(parameters)
	if err != nil {
		return nil, err
	}

	containers, err := dockerClient.ContainerList(
		ctx,
		container.ListOptions{
			All:     true,
			Filters: filters.Args{},
		},
	)
	if err != nil {
		return nil, err
	}

	options := a.buildOptions(containers, tcpOnly)

	if searchTerms != nil {
		filteredOptions := make([]*containerMetadata, 0)
		for _, option := range options {
			if strings.Contains(strings.ToLower(option.name), strings.ToLower(*searchTerms)) {
				filteredOptions = append(filteredOptions, option)
			}
		}

		options = filteredOptions
	}

	return options, nil
}

func (a *Adapter) buildOptions(containers []container.Summary, tcpOnly bool) []*containerMetadata {
	var options []*containerMetadata
	for _, item := range containers {
		for _, port := range item.Ports {
			if tcpOnly && strings.ToUpper(port.Type) != "TCP" {
				continue
			}

			if port.PublicPort != 0 {
				option := &containerMetadata{
					name:      strings.TrimPrefix(item.Names[0], "/"),
					container: &item,
					port:      &port,
				}
				options = append(options, option)
			}
		}
	}
	return options
}

func toAdapterOption(option *containerMetadata) *integration.AdapterOption {
	port := option.port
	protocol := strings.ToUpper(port.Type)

	return &integration.AdapterOption{
		ID:       fmt.Sprintf("%s:%d", option.container.ID, port.PublicPort),
		Name:     option.name,
		Port:     int(port.PublicPort),
		Protocol: integration.Protocol(protocol),
	}
}

func startClient(parameters map[string]any) (*client.Client, error) {
	connectionMode := parameters[connectionModeField.ID].(string)

	var connectionString string
	switch connectionMode {
	case "SOCKET":
		socketPath := parameters[socketPathField.ID].(string)
		connectionString = "unix://" + socketPath
	case "TCP":
		hostUrl := parameters[hostUrlField.ID].(string)
		connectionString = hostUrl
	default:
		return nil, core_error.New("Invalid connection mode", false)
	}

	return client.NewClientWithOpts(
		client.WithHost(connectionString),
		client.WithAPIVersionNegotiation(),
	)
}
