package resolver

import (
	"context"
	"strings"

	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/integration/docker/fields"
)

type Resolver interface {
	ResolveOptions(ctx context.Context, tcpOnly bool, searchTerms *string) (*[]Option, error)
	ResolveOptionByID(ctx context.Context, optionId string) (*Option, error)
}

func For(parameters map[string]any) (Resolver, error) {
	var connectionUrl string
	switch parameters[fields.ConnectionMode.ID].(string) {
	case "SOCKET":
		socketPath := parameters[fields.SocketPath.ID].(string)
		connectionUrl = "unix://" + socketPath
	case "TCP":
		hostUrl := parameters[fields.HostURL.ID].(string)
		connectionUrl = hostUrl
	default:
		return nil, coreerror.New("Invalid connection mode", false)
	}

	dockerClient, err := client.NewClientWithOpts(client.WithHost(connectionUrl), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	publicUrl, _ := parameters[fields.ProxyURL.ID].(string)
	swarmEnabled, useServiceMesh, dnsResolvers := extractSwarmParams(parameters)

	if swarmEnabled {
		return &swarmAdapter{
			client:         dockerClient,
			useServiceMesh: useServiceMesh,
			dnsResolvers:   dnsResolvers,
			publicUrl:      publicUrl,
		}, nil
	}

	var useNameAsID bool
	if rawValue, exists := parameters[fields.UseContainerNameAsID.ID]; exists {
		useNameAsID = rawValue.(bool)
	}

	return &simpleAdapter{
		client:      dockerClient,
		publicUrl:   publicUrl,
		useNameAsID: useNameAsID,
	}, nil
}

func extractSwarmParams(parameters map[string]any) (bool, bool, *[]string) {
	swarmMode := false
	useServiceMesh := false
	var dnsResolvers *[]string

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
				dnsResolvers = &result
			}
		}
	}

	return swarmMode, useServiceMesh, dnsResolvers
}
