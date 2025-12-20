package resolver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/integration"
)

type swarmAdapter struct {
	client         *client.Client
	useServiceMesh bool
	dnsResolvers   *[]string
}

func FromServices(c *client.Client, useServiceMesh bool, dnsResolvers *[]string) Resolver {
	return &swarmAdapter{
		client:         c,
		useServiceMesh: useServiceMesh,
		dnsResolvers:   dnsResolvers,
	}
}

func (s *swarmAdapter) ResolveOptionByID(ctx context.Context, id string) (*Option, error) {
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

func (s *swarmAdapter) ResolveOptions(ctx context.Context, tcpOnly bool, searchTerms *string) (*[]Option, error) {
	args := filters.NewArgs()
	if searchTerms != nil {
		q := strings.TrimSpace(*searchTerms)
		if q != "" {
			args.Add("name", q)
			args.Add("id", q)
		}
	}

	services, err := s.client.ServiceList(ctx, swarm.ServiceListOptions{
		Filters: args,
	})
	if err != nil {
		return nil, err
	}

	return s.buildServiceOptions(services, tcpOnly), nil
}

func (s *swarmAdapter) buildServiceOptions(services []swarm.Service, tcpOnly bool) *[]Option {
	optionIDs := make(map[string]bool)
	options := make([]Option, 0, len(services))

	for _, service := range services {
		if service.Spec.EndpointSpec == nil {
			continue
		}

		for _, port := range service.Spec.EndpointSpec.Ports {
			if tcpOnly && port.Protocol != swarm.PortConfigProtocolTCP {
				continue
			}

			if option := s.buildServiceOption(&port, &service); option != nil && !optionIDs[option.ID] {
				options = append(options, *option)
				optionIDs[option.ID] = true
			}
		}
	}

	return &options
}

func (s *swarmAdapter) buildServiceOption(port *swarm.PortConfig, service *swarm.Service) *Option {
	portNumber := port.PublishedPort

	if portNumber == 0 {
		return nil
	}

	qualifierType := ingressQualifier
	if port.PublishMode == swarm.PortConfigPublishModeHost {
		qualifierType = hostQualifier
	}

	return &Option{
		DriverOption: &integration.DriverOption{
			ID:           fmt.Sprintf("%s:%d:%s", service.ID, portNumber, qualifierType),
			Name:         service.Spec.Name,
			Port:         int(portNumber),
			Qualifier:    ptr.Of(qualifierType),
			Protocol:     integration.Protocol(port.Protocol),
			DNSResolvers: s.dnsResolvers,
		},
		urlResolver: func(ctx context.Context, option *Option, publicUrl string) (*string, error) {
			return s.buildServiceOptionURL(ctx, option, publicUrl, service)
		},
	}
}

func (s *swarmAdapter) buildServiceOptionURL(
	ctx context.Context,
	option *Option,
	publicUrl string,
	service *swarm.Service,
) (*string, error) {
	var targetHost string

	if s.useServiceMesh {
		targetHost = service.Spec.Name
	} else if publicUrl != "" {
		uri, err := url.Parse(publicUrl)
		if err != nil {
			return nil, err
		}

		targetHost = uri.Hostname()
	} else {
		nodes, err := s.client.NodeList(ctx, swarm.NodeListOptions{})
		if err != nil {
			return nil, err
		}

		if len(nodes) == 0 {
			return nil, fmt.Errorf("no nodes found")
		}

		var leaderNode *swarm.Node
		for _, node := range nodes {
			if node.ManagerStatus.Leader {
				leaderNode = &node
				break
			}
		}

		if leaderNode == nil {
			return nil, fmt.Errorf("no leader node found")
		}

		targetHost = leaderNode.Status.Addr
	}

	result := fmt.Sprintf("http://%s:%d", targetHost, option.Port)
	return &result, nil
}
