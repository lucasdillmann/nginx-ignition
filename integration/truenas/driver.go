package truenas

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/integration/truenas/client"
)

type Driver struct {
	client        *client.Client
	cacheDuration int
}

func newDriver(configuration *configuration.Configuration) (*Driver, error) {
	cacheDuration, err := configuration.GetInt("nginx-ignition.integration.truenas.api-cache-timeout-seconds")
	if err != nil {
		return nil, err
	}

	return &Driver{cacheDuration: cacheDuration}, nil
}

func (a *Driver) ID() string {
	return "TRUENAS"
}

func (a *Driver) Name() string {
	return "TrueNAS"
}

func (a *Driver) Description() string {
	return "TrueNAS allows, alongside many other things, to run your favorite apps under Docker containers. With this " +
		"integration enabled, you will be able to easily pick any app exposing a service in your TrueNAS as a " +
		"target for your nginx ignition's host routes."
}

func (a *Driver) ConfigurationFields() []dynamicfields.DynamicField {
	return []dynamicfields.DynamicField{
		urlField,
		proxyUrlField,
		usernameField,
		passwordField,
	}
}

func (a *Driver) GetAvailableOptions(
	_ context.Context,
	parameters map[string]any,
	_, _ int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[integration.DriverOption], error) {
	apps, err := a.getAvailableApps(parameters)
	if err != nil {
		return nil, err
	}

	options := a.buildOptions(apps, tcpOnly)

	if searchTerms != nil {
		filteredOptions := make([]integration.DriverOption, 0)
		for _, option := range options {
			if strings.Contains(strings.ToLower(option.Name), strings.ToLower(*searchTerms)) {
				filteredOptions = append(filteredOptions, option)
			}
		}

		options = filteredOptions
	}

	resultSize := len(options)
	return pagination.New(0, resultSize, resultSize, options), nil
}

func (a *Driver) GetAvailableOptionById(
	_ context.Context,
	parameters map[string]any,
	id string,
) (*integration.DriverOption, error) {
	parts := strings.Split(id, ":")
	appId := parts[0]
	containerPort := parts[1]

	app, port, err := a.getWorkloadPort(parameters, appId, containerPort)
	if err != nil {
		return nil, err
	}

	if app == nil || port == nil {
		return nil, nil
	}

	return &integration.DriverOption{
		ID:       id,
		Name:     app.Name,
		Port:     port.HostPorts[0].HostPort,
		Protocol: integration.Protocol(strings.ToUpper(port.Protocol)),
	}, nil
}

func (a *Driver) GetOptionProxyURL(
	_ context.Context,
	parameters map[string]any,
	id string,
) (*string, []string, error) {
	baseUrl := parameters[urlField.ID].(string)
	proxyUrl := parameters[proxyUrlField.ID].(string)
	parts := strings.Split(id, ":")
	appId := parts[0]
	containerPort := parts[1]

	_, port, err := a.getWorkloadPort(parameters, appId, containerPort)
	if err != nil {
		return nil, nil, err
	}

	if port == nil || len(port.HostPorts) == 0 {
		return nil, nil, fmt.Errorf("unable to resolve proxy URL for %s: service is probably offline/stopped", id)
	}

	hostPort := port.HostPorts[0].HostPort
	hostIp := port.HostPorts[0].HostIp

	var endpoint string
	switch {
	case proxyUrl != "":
		parseResult, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, nil, err
		}

		endpoint = parseResult.Host

	case hostIp == "0.0.0.0":
		parseResult, err := url.Parse(baseUrl)
		if err != nil {
			return nil, nil, err
		}

		endpoint = parseResult.Host

	default:
		endpoint = hostIp
	}

	output := fmt.Sprintf("http://%s:%d", endpoint, hostPort)
	return &output, nil, nil
}

func (a *Driver) getWorkloadPort(
	parameters map[string]any,
	appId, containerPort string,
) (*client.AvailableAppDTO, *client.WorkloadPortDTO, error) {
	apps, err := a.getAvailableApps(parameters)
	if err != nil {
		return nil, nil, err
	}

	for _, app := range apps {
		if app.ID == appId {
			for _, port := range app.ActiveWorkloads.UsedPorts {
				if strconv.Itoa(port.ContainerPort) == containerPort {
					return &app, &port, nil
				}
			}
		}
	}

	return nil, nil, nil
}

func (a *Driver) buildOptions(apps []client.AvailableAppDTO, tcpOnly bool) []integration.DriverOption {
	options := make([]integration.DriverOption, 0)

	for _, app := range apps {
		for _, port := range app.ActiveWorkloads.UsedPorts {
			for _, hostPort := range port.HostPorts {
				if strings.Contains(hostPort.HostIp, ":") {
					continue
				}

				protocol := strings.ToUpper(port.Protocol)
				if tcpOnly && protocol != "TCP" {
					continue
				}

				options = append(options, integration.DriverOption{
					ID:       fmt.Sprintf("%s:%d", app.ID, port.ContainerPort),
					Name:     app.Name,
					Port:     hostPort.HostPort,
					Protocol: integration.Protocol(protocol),
				})
			}
		}
	}

	return options
}

func (a *Driver) getAvailableApps(parameters map[string]any) ([]client.AvailableAppDTO, error) {
	baseUrl := parameters[urlField.ID].(string)
	username := parameters[usernameField.ID].(string)
	password := parameters[passwordField.ID].(string)

	if a.client == nil {
		a.client = client.New(baseUrl, username, password, a.cacheDuration)
	} else {
		a.client.UpdateCredentials(baseUrl, username, password)
	}

	return a.client.GetAvailableApps()
}
