package resolver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/integration"
)

type simpleAdapter struct {
	client    *client.Client
	publicUrl string
}

func (s *simpleAdapter) ResolveOptionByID(ctx context.Context, id string) (*Option, error) {
	availableOptions, err := s.ResolveOptions(ctx, false, nil)
	if err != nil {
		return nil, err
	}

	for _, item := range *availableOptions {
		if item.ID == id {
			return &item, nil
		}
	}

	return nil, nil
}

func (s *simpleAdapter) ResolveOptions(ctx context.Context, tcpOnly bool, searchTerms *string) (*[]Option, error) {
	containers, err := s.client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	if searchTerms != nil && strings.TrimSpace(*searchTerms) != "" {
		normalizedTerms := strings.ToLower(strings.TrimSpace(*searchTerms))
		filteredResults := make([]container.Summary, 0)

		for _, item := range containers {
			matches := false

			for _, name := range item.Names {
				if strings.Contains(strings.ToLower(name), normalizedTerms) {
					matches = true
					break
				}
			}

			if matches {
				filteredResults = append(filteredResults, item)
			}
		}

		containers = filteredResults
	}

	return s.buildOptions(containers, tcpOnly), nil
}

func (s *simpleAdapter) buildOptions(containers []container.Summary, tcpOnly bool) *[]Option {
	optionIDs := make(map[string]bool)
	options := make([]Option, 0, len(containers))

	for _, item := range containers {
		for _, port := range item.Ports {
			if tcpOnly && strings.ToUpper(port.Type) != "TCP" {
				continue
			}

			if option := s.buildOption(&port, &item, true); option != nil && !optionIDs[option.ID] {
				options = append(options, *option)
				optionIDs[option.ID] = true
			}

			if option := s.buildOption(&port, &item, false); option != nil && !optionIDs[option.ID] {
				options = append(options, *option)
				optionIDs[option.ID] = true
			}
		}
	}

	return &options
}

func (s *simpleAdapter) buildOption(port *container.Port, item *container.Summary, usePublicPort bool) *Option {
	portNumber := port.PrivatePort
	qualifierType := containerQualifier

	if usePublicPort {
		portNumber = port.PublicPort
		qualifierType = hostQualifier
	}

	if portNumber == 0 {
		return nil
	}

	return &Option{
		DriverOption: &integration.DriverOption{
			ID:        fmt.Sprintf("%s:%d:%s", item.ID, portNumber, qualifierType),
			Name:      strings.TrimPrefix(item.Names[0], "/"),
			Port:      int(portNumber),
			Qualifier: ptr.Of(qualifierType),
			Protocol:  integration.Protocol(port.Type),
		},
		urlResolver: func(_ context.Context, option *Option) (*string, *[]string, error) {
			return s.buildOptionURL(option, item)
		},
	}
}

func (s *simpleAdapter) buildOptionURL(option *Option, summary *container.Summary) (*string, *[]string, error) {
	var targetHost string
	if s.publicUrl != "" && *option.Qualifier == hostQualifier {
		uri, err := url.Parse(s.publicUrl)
		if err != nil {
			return nil, nil, err
		}

		targetHost = uri.Hostname()
	} else {
		if len(summary.NetworkSettings.Networks) > 0 {
			for _, network := range summary.NetworkSettings.Networks {
				targetHost = network.IPAddress
				break
			}
		}

		if targetHost == "" {
			return nil, nil, fmt.Errorf("no network or IP address found for the container with ID %s", option.ID)
		}
	}

	result := fmt.Sprintf("http://%s:%d", targetHost, option.Port)
	return &result, nil, nil
}
