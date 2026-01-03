package resolver

import (
	"context"
	"strings"

	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/integration/docker/fields"
)

type Resolver interface {
	ResolveOptions(ctx context.Context, tcpOnly bool, searchTerms *string) ([]Option, error)
	ResolveOptionByID(ctx context.Context, optionID string) (*Option, error)
}

func For(parameters map[string]any) (Resolver, error) {
	var connectionURL string
	switch parameters[fields.ConnectionMode.ID].(string) {
	case fields.SocketConnectionMode:
		socketPath := parameters[fields.SocketPath.ID].(string)
		connectionURL = "unix://" + socketPath
	case fields.TCPConnectionMode:
		hostURL := parameters[fields.HostURL.ID].(string)
		connectionURL = hostURL
	default:
		return nil, coreerror.New("Invalid connection mode", false)
	}

	dockerClient, err := client.NewClientWithOpts(
		client.WithHost(connectionURL),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	publicURL, _ := parameters[fields.ProxyURL.ID].(string)
	swarmEnabled, useServiceMesh, dnsResolvers := extractSwarmParams(parameters)

	if swarmEnabled {
		return &swarmAdapter{
			client:         dockerClient,
			useServiceMesh: useServiceMesh,
			dnsResolvers:   dnsResolvers,
			publicURL:      publicURL,
		}, nil
	}

	var useNameAsID bool
	if rawValue, exists := parameters[fields.UseContainerNameAsID.ID]; exists {
		useNameAsID = rawValue.(bool)
	}

	return &simpleAdapter{
		client:      dockerClient,
		publicURL:   publicURL,
		useNameAsID: useNameAsID,
	}, nil
}

func extractSwarmParams(parameters map[string]any) (
	swarmMode bool,
	useServiceMesh bool,
	dnsResolvers []string,
) {
	if rawValue, exists := parameters[fields.SwarmMode.ID]; exists {
		swarmMode = rawValue.(bool)
	}

	if rawValue, exists := parameters[fields.SwarmServiceMesh.ID]; exists {
		useServiceMesh = rawValue.(bool)
	}

	if useServiceMesh {
		if rawValue, exists := parameters[fields.SwarmDNSResolvers.ID]; exists {
			textValue := rawValue.(string)
			result := make([]string, 0)

			for _, value := range strings.Split(textValue, "\n") {
				if normalizedValue := strings.TrimSpace(value); normalizedValue != "" {
					result = append(result, normalizedValue)
				}
			}

			if len(result) != 0 {
				dnsResolvers = result
			}
		}
	}

	return swarmMode, useServiceMesh, dnsResolvers
}
