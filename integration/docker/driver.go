package docker

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/integration"
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
	options, err := a.resolveAvailableOptions(ctx, parameters, searchTerms, tcpOnly)
	if err != nil {
		return nil, err
	}

	totalItems := len(options)
	driverOptions := make([]*integration.DriverOption, totalItems)
	for index, option := range options {
		driverOptions[index] = toDriverOption(option)
	}

	return pagination.New(0, totalItems, totalItems, driverOptions), nil
}

func (a *Driver) GetAvailableOptionById(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*integration.DriverOption, error) {
	option, err := a.resolveAvailableOptionById(ctx, parameters, id)
	if err != nil || option == nil {
		return nil, err
	}

	return toDriverOption(option), nil
}

func (a *Driver) GetOptionProxyURL(
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
	if publicUrl != "" && option.qualifier == hostQualifier {
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

	result := fmt.Sprintf("http://%s:%d", targetHost, option.portNumber)
	return &result, nil
}

func (a *Driver) resolveAvailableOptionById(
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*containerMetadata, error) {
	options, err := a.resolveAvailableOptions(ctx, parameters, nil, false)
	if err != nil {
		return nil, err
	}

	idParts := strings.Split(id, ":")
	if len(idParts) != 3 {
		return nil, coreerror.New("Invalid option ID", true)
	}

	for _, option := range options {
		if option.id == id {
			return option, nil
		}
	}

	return nil, nil
}

func (a *Driver) resolveAvailableOptions(
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

func (a *Driver) buildOptions(containers []container.Summary, tcpOnly bool) []*containerMetadata {
	optionIDs := make(map[string]bool)
	options := make([]*containerMetadata, 0, len(containers))

	for _, item := range containers {
		for _, port := range item.Ports {
			if tcpOnly && strings.ToUpper(port.Type) != "TCP" {
				continue
			}

			if metadata := buildOption(&port, &item, true); metadata != nil && !optionIDs[metadata.id] {
				options = append(options, metadata)
				optionIDs[metadata.id] = true
			}

			if metadata := buildOption(&port, &item, false); metadata != nil && !optionIDs[metadata.id] {
				options = append(options, metadata)
				optionIDs[metadata.id] = true
			}
		}
	}

	return options
}

func buildOption(port *container.Port, item *container.Summary, usePublicPort bool) *containerMetadata {
	portNumber := port.PrivatePort
	qualifierType := containerQualifier

	if usePublicPort {
		portNumber = port.PublicPort
		qualifierType = hostQualifier
	}

	if portNumber == 0 {
		return nil
	}

	return &containerMetadata{
		id:         fmt.Sprintf("%s:%d:%s", item.ID, portNumber, qualifierType),
		name:       strings.TrimPrefix(item.Names[0], "/"),
		container:  item,
		portNumber: int(portNumber),
		qualifier:  qualifierType,
		protocol:   port.Type,
	}
}

func toDriverOption(option *containerMetadata) *integration.DriverOption {
	return &integration.DriverOption{
		ID:        option.id,
		Name:      option.name,
		Port:      option.portNumber,
		Qualifier: ptr.Of(string(option.qualifier)),
		Protocol:  integration.Protocol(option.protocol),
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
		return nil, coreerror.New("Invalid connection mode", false)
	}

	return client.NewClientWithOpts(
		client.WithHost(connectionString),
		client.WithAPIVersionNegotiation(),
	)
}
