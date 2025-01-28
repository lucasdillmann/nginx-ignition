package docker

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"net/url"
	"strings"
)

type containerMetadata struct {
	name      string
	container *types.Container
	port      *types.Port
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
	parameters map[string]interface{},
	_, _ int,
	searchTerms *string,
) (*pagination.Page[*integration.AdapterOption], error) {
	options, err := a.resolveAvailableOptions(parameters, searchTerms)
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
	parameters map[string]interface{},
	id string,
) (*integration.AdapterOption, error) {
	option, err := a.resolveAvailableOptionById(parameters, id)
	if err != nil || option == nil {
		return nil, err
	}

	return toAdapterOption(option), nil
}

func (a *Adapter) GetOptionProxyUrl(
	parameters map[string]interface{},
	id string,
) (*string, error) {
	option, err := a.resolveAvailableOptionById(parameters, id)
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
	parameters map[string]interface{},
	id string,
) (*containerMetadata, error) {
	options, err := a.resolveAvailableOptions(parameters, nil)
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
	parameters map[string]interface{},
	searchTerms *string,
) ([]*containerMetadata, error) {
	dockerClient, err := startClient(parameters)
	if err != nil {
		return nil, err
	}

	containers, err := dockerClient.ContainerList(
		context.Background(),
		container.ListOptions{
			All:     true,
			Filters: filters.Args{},
		},
	)
	if err != nil {
		return nil, err
	}

	var options []*containerMetadata
	for _, item := range containers {
		for _, port := range item.Ports {
			if port.Type == "tcp" && port.PublicPort != 0 {
				option := &containerMetadata{
					name:      strings.TrimPrefix(item.Names[0], "/"),
					container: &item,
					port:      &port,
				}
				options = append(options, option)
			}
		}
	}

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

func toAdapterOption(option *containerMetadata) *integration.AdapterOption {
	return &integration.AdapterOption{
		ID:   fmt.Sprintf("%s:%d", option.container.ID, option.port.PublicPort),
		Name: fmt.Sprintf("%s (%d HTTP)", option.name, option.port.PublicPort),
	}
}

func startClient(parameters map[string]interface{}) (*client.Client, error) {
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
	)
}
