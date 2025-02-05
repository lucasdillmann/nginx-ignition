package truenas

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/integration/truenas/client"
	"net/url"
	"strconv"
	"strings"
)

type Adapter struct {
	cacheDuration int
	client        *client.Client
}

func newAdapter(configuration *configuration.Configuration) (*Adapter, error) {
	cacheDuration, err := configuration.GetInt("nginx-ignition.integration.truenas.api-cache-timeout-seconds")
	if err != nil {
		return nil, err
	}

	return &Adapter{cacheDuration: cacheDuration}, nil
}

func (a *Adapter) ID() string {
	return "TRUENAS_SCALE"
}

func (a *Adapter) Name() string {
	return "TrueNAS Scale"
}

func (a *Adapter) Priority() int {
	return 2
}

func (a *Adapter) Description() string {
	return "TrueNAS allows, alongside many other things, to run your favorite apps under Docker containers. With this " +
		"integration enabled, you will be able to easily pick any app exposing a service in your TrueNAS as a " +
		"target for your nginx ignition's host routes."
}

func (a *Adapter) ConfigurationFields() []*dynamic_fields.DynamicField {
	return []*dynamic_fields.DynamicField{
		&urlField,
		&proxyUrlField,
		&usernameField,
		&passwordField,
	}
}

func (a *Adapter) GetAvailableOptions(
	parameters map[string]any,
	_, _ int,
	searchTerms *string,
) (*pagination.Page[*integration.AdapterOption], error) {
	apps, err := a.getAvailableApps(parameters)
	if err != nil {
		return nil, err
	}

	options := a.buildOptions(apps)

	if searchTerms != nil {
		var filteredOptions []*integration.AdapterOption
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

func (a *Adapter) GetAvailableOptionById(
	parameters map[string]any,
	id string,
) (*integration.AdapterOption, error) {
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

	return &integration.AdapterOption{
		ID:   id,
		Name: app.Name + " (port " + strconv.Itoa(port.HostPorts[0].HostPort) + " HTTP)",
	}, nil
}

func (a *Adapter) GetOptionProxyUrl(
	parameters map[string]any,
	id string,
) (*string, error) {
	baseUrl := parameters[urlField.ID].(string)
	proxyUrl := parameters[proxyUrlField.ID].(string)
	parts := strings.Split(id, ":")
	appId := parts[0]
	containerPort := parts[1]

	_, port, err := a.getWorkloadPort(parameters, appId, containerPort)
	if err != nil {
		return nil, err
	}

	hostPort := port.HostPorts[0].HostPort
	hostIp := port.HostPorts[0].HostIp

	var endpoint string
	if proxyUrl != "" {
		parseResult, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}

		endpoint = parseResult.Host
	} else if hostIp == "0.0.0.0" {
		parseResult, err := url.Parse(baseUrl)
		if err != nil {
			return nil, err
		}

		endpoint = parseResult.Host
	} else {
		endpoint = hostIp
	}

	output := "http://" + endpoint + ":" + strconv.Itoa(hostPort)
	return &output, nil
}

func (a *Adapter) getWorkloadPort(
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

func (a *Adapter) buildOptions(apps []client.AvailableAppDTO) []*integration.AdapterOption {
	var options []*integration.AdapterOption

	for _, app := range apps {
		for _, port := range app.ActiveWorkloads.UsedPorts {
			if strings.ToLower(port.Protocol) == "tcp" {
				for _, hostPort := range port.HostPorts {
					if strings.Contains(hostPort.HostIp, ":") {
						continue
					}

					options = append(options, &integration.AdapterOption{
						ID:   app.ID + ":" + strconv.Itoa(port.ContainerPort),
						Name: app.Name + " (port " + strconv.Itoa(hostPort.HostPort) + " HTTP)",
					})
				}
			}
		}
	}

	return options
}

func (a *Adapter) getAvailableApps(parameters map[string]any) ([]client.AvailableAppDTO, error) {
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
