package resolver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/integration"
)

type simpleAdapter struct {
	client      *client.Client
	publicURL   string
	useNameAsID bool
}

func (s *simpleAdapter) ResolveOptionByID(ctx context.Context, id string) (*Option, error) {
	availableOptions, err := s.ResolveOptions(ctx, false, nil)
	if err != nil {
		return nil, err
	}

	for _, item := range availableOptions {
		if item.ID == id {
			return &item, nil
		}
	}

	return nil, nil
}

func (s *simpleAdapter) ResolveOptions(
	ctx context.Context,
	tcpOnly bool,
	searchTerms *string,
) ([]Option, error) {
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

	return s.buildOptions(ctx, containers, tcpOnly), nil
}

func (s *simpleAdapter) buildOptions(
	ctx context.Context,
	containers []container.Summary,
	tcpOnly bool,
) []Option {
	optionIDs := make(map[string]bool)
	options := make([]Option, 0, len(containers))

	for _, item := range containers {
		for _, port := range item.Ports {
			if tcpOnly && strings.ToUpper(port.Type) != "TCP" {
				continue
			}

			if option := s.buildOption(
				ctx,
				&port,
				&item,
				true,
			); option != nil && !optionIDs[option.ID] {
				options = append(options, *option)
				optionIDs[option.ID] = true
			}

			if option := s.buildOption(
				ctx,
				&port,
				&item,
				false,
			); option != nil && !optionIDs[option.ID] {
				options = append(options, *option)
				optionIDs[option.ID] = true
			}
		}
	}

	return options
}

func (s *simpleAdapter) buildOption(
	ctx context.Context,
	port *container.Port,
	item *container.Summary,
	usePublicPort bool,
) *Option {
	portNumber := port.PrivatePort
	qualifierType := containerQualifier

	if usePublicPort {
		portNumber = port.PublicPort
		qualifierType = hostQualifier
	}

	if portNumber == 0 {
		return nil
	}

	itemID := item.ID
	itemName := strings.TrimPrefix(item.Names[0], "/")

	if s.useNameAsID {
		itemID = normalizeContainerName(itemName, itemID)
	}

	return &Option{
		DriverOption: integration.DriverOption{
			ID:        fmt.Sprintf("%s:%d:%s", itemID, portNumber, qualifierType),
			Name:      itemName,
			Port:      int(portNumber),
			Qualifier: new(qualifierType),
			Protocol:  integration.Protocol(port.Type),
		},
		urlResolver: func(_ context.Context, option *Option) (*string, []string, error) {
			return s.buildOptionURL(ctx, option, item)
		},
	}
}

func (s *simpleAdapter) buildOptionURL(
	ctx context.Context,
	option *Option,
	summary *container.Summary,
) (*string, []string, error) {
	var targetHost string
	if s.publicURL != "" && *option.Qualifier == hostQualifier {
		uri, err := url.Parse(s.publicURL)
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
			return nil, nil, coreerror.New(
				i18n.M(ctx, i18n.K.IntegrationDockerResolverNoNetwork).V("id", option.ID),
				false,
			)
		}
	}

	return new(fmt.Sprintf(httpURLTemplate, targetHost, option.Port)), nil, nil
}

func normalizeContainerName(name, containerID string) string {
	if strings.TrimSpace(name) == "" {
		return containerID
	}

	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = containerNameGeneralNormalizationRegex.ReplaceAllString(name, "_")
	name = containerNameUnderscoreNormalizationRegex.ReplaceAllString(name, "_")

	if name == "" {
		return containerID
	}

	return name
}
