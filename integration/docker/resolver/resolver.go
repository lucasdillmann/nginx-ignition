package resolver

import (
	"context"
	"strings"

	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/integration/docker/fields"
)

type Resolver interface {
	ResolveOptions(ctx context.Context, tcpOnly bool, searchTerms *string) ([]Option, error)
	ResolveOptionByID(ctx context.Context, optionID string) (*Option, error)
}

func For(ctx context.Context, parameters map[string]any) (Resolver, error) {
	var connectionURL string
	switch parameters[fields.ConnectionModeFieldID].(string) {
	case fields.SocketConnectionMode:
		socketPath := parameters[fields.SocketPathFieldID].(string)
		switch {
		case strings.Contains(socketPath, "://"):
			connectionURL = socketPath
		case strings.HasPrefix(socketPath, "//./pipe/"),
			strings.HasPrefix(socketPath, "\\\\.\\pipe\\"):
			connectionURL = "npipe://" + socketPath
		default:
			connectionURL = "unix://" + socketPath
		}
	case fields.TCPConnectionMode:
		hostURL := parameters[fields.HostURLFieldID].(string)
		connectionURL = hostURL
	default:
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.IntegrationDockerResolverInvalidConnectionMode),
			false,
		)
	}

	dockerClient, err := client.NewClientWithOpts(
		client.WithHost(connectionURL),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	publicURL, _ := parameters[fields.ProxyURLFieldID].(string)
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
	if rawValue, exists := parameters[fields.UseContainerNameAsIDFieldID]; exists {
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
	if rawValue, exists := parameters[fields.SwarmModeFieldID]; exists {
		swarmMode = rawValue.(bool)
	}

	if rawValue, exists := parameters[fields.SwarmServiceMeshFieldID]; exists {
		useServiceMesh = rawValue.(bool)
	}

	if useServiceMesh {
		if rawValue, exists := parameters[fields.SwarmDNSResolversFieldID]; exists {
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
