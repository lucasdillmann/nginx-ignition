package truenas

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/integration/truenas/client"
)

type Driver struct {
	client        *client.Client
	cacheDuration int
}

func newDriver(cfg *configuration.Configuration) (*Driver, error) {
	cacheDuration, err := cfg.GetInt("nginx-ignition.integration.truenas.api-cache-timeout-seconds")
	if err != nil {
		return nil, err
	}

	return &Driver{cacheDuration: cacheDuration}, nil
}

func (a *Driver) ID() string {
	return "TRUENAS"
}

func (a *Driver) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.TruenasCommonName)
}

func (a *Driver) Description(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.TruenasCommonDescription)
}

func (a *Driver) ConfigurationFields(ctx context.Context) []dynamicfields.DynamicField {
	return dynamicFields(ctx)
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

func (a *Driver) GetAvailableOptionByID(
	_ context.Context,
	parameters map[string]any,
	id string,
) (*integration.DriverOption, error) {
	parts := strings.Split(id, ":")
	appID := parts[0]
	containerPort := parts[1]

	app, port, err := a.getWorkloadPort(parameters, appID, containerPort)
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
	ctx context.Context,
	parameters map[string]any,
	id string,
) (*string, []string, error) {
	baseURL := parameters[urlFieldID].(string)
	proxyURL := parameters[proxyURLFieldID].(string)
	parts := strings.Split(id, ":")
	appID := parts[0]
	containerPort := parts[1]

	_, port, err := a.getWorkloadPort(parameters, appID, containerPort)
	if err != nil {
		return nil, nil, err
	}

	if port == nil || len(port.HostPorts) == 0 {
		return nil, nil, coreerror.New(
			i18n.M(ctx, i18n.K.TruenasErrorProxyUrl).V("id", id),
			false,
		)
	}

	hostPort := port.HostPorts[0].HostPort
	hostIP := port.HostPorts[0].HostIP

	var endpoint string
	switch {
	case proxyURL != "":
		parseResult, err := url.Parse(proxyURL)
		if err != nil {
			return nil, nil, err
		}

		endpoint = parseResult.Host

	case hostIP == "0.0.0.0":
		parseResult, err := url.Parse(baseURL)
		if err != nil {
			return nil, nil, err
		}

		endpoint = parseResult.Host

	default:
		endpoint = hostIP
	}

	output := fmt.Sprintf("http://%s:%d", endpoint, hostPort)
	return &output, nil, nil
}

func (a *Driver) getWorkloadPort(
	parameters map[string]any,
	appID, containerPort string,
) (*client.AvailableAppDTO, *client.WorkloadPortDTO, error) {
	apps, err := a.getAvailableApps(parameters)
	if err != nil {
		return nil, nil, err
	}

	for _, app := range apps {
		if app.ID == appID {
			for _, port := range app.ActiveWorkloads.UsedPorts {
				if strconv.Itoa(port.ContainerPort) == containerPort {
					return &app, &port, nil
				}
			}
		}
	}

	return nil, nil, nil
}

func (a *Driver) buildOptions(
	apps []client.AvailableAppDTO,
	tcpOnly bool,
) []integration.DriverOption {
	options := make([]integration.DriverOption, 0)

	for _, app := range apps {
		for _, port := range app.ActiveWorkloads.UsedPorts {
			for _, hostPort := range port.HostPorts {
				if strings.Contains(hostPort.HostIP, ":") {
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
	baseURL := parameters[urlFieldID].(string)
	username := parameters[usernameFieldID].(string)
	password := parameters[passwordFieldID].(string)

	if a.client == nil {
		a.client = client.New(baseURL, username, password, a.cacheDuration)
	} else {
		a.client.UpdateCredentials(baseURL, username, password)
	}

	return a.client.GetAvailableApps()
}
