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
	publicUrl      string
	dnsResolvers   []string
	useServiceMesh bool
}

func (s *swarmAdapter) ResolveOptionByID(ctx context.Context, id string) (*Option, error) {
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

func (s *swarmAdapter) ResolveOptions(ctx context.Context, tcpOnly bool, searchTerms *string) ([]Option, error) {
	services, err := s.client.ServiceList(ctx, swarm.ServiceListOptions{})
	if err != nil {
		return nil, err
	}

	if searchTerms != nil && strings.TrimSpace(*searchTerms) != "" {
		normalizedTerms := strings.ToLower(strings.TrimSpace(*searchTerms))
		filteredResults := make([]swarm.Service, 0)

		for _, service := range services {
			if strings.Contains(strings.ToLower(service.Spec.Name), normalizedTerms) {
				filteredResults = append(filteredResults, service)
			}
		}

		services = filteredResults
	}

	return s.buildServiceOptions(services, tcpOnly), nil
}

func (s *swarmAdapter) buildServiceOptions(services []swarm.Service, tcpOnly bool) []Option {
	optionIDs := make(map[string]bool)
	options := make([]Option, 0, len(services))

	for _, service := range services {
		if service.Spec.EndpointSpec == nil || service.Spec.EndpointSpec.Ports == nil {
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

	return options
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
		DriverOption: integration.DriverOption{
			ID:           fmt.Sprintf("%s:%d:%s", service.ID, portNumber, qualifierType),
			Name:         service.Spec.Name,
			Port:         int(portNumber),
			Qualifier:    ptr.Of(qualifierType),
			Protocol:     integration.Protocol(port.Protocol),
			DNSResolvers: s.dnsResolvers,
		},
		privatePort: int(port.TargetPort),
		urlResolver: func(ctx context.Context, option *Option) (*string, []string, error) {
			return s.buildServiceOptionURL(ctx, option, service)
		},
	}
}

func (s *swarmAdapter) buildServiceOptionURL(
	ctx context.Context,
	option *Option,
	service *swarm.Service,
) (*string, []string, error) {
	if s.useServiceMesh && *option.Qualifier == ingressQualifier {
		dnsResolvers := s.dnsResolvers
		if len(dnsResolvers) == 0 {
			dnsResolvers = []string{defaultDockerDNSIP}
		}

		fullUrl := fmt.Sprintf("http://%s:%d", service.Spec.Name, option.privatePort)
		return &fullUrl, dnsResolvers, nil
	}

	if s.publicUrl != "" {
		uri, err := url.Parse(s.publicUrl)
		if err != nil {
			return nil, nil, err
		}

		fullUrl := fmt.Sprintf("http://%s:%d", uri.Hostname(), option.Port)
		return &fullUrl, nil, nil
	}

	nodeAddress, err := s.findNodeAddress(ctx, service)
	if err != nil {
		return nil, nil, err
	}

	fullUrl := fmt.Sprintf("http://%s:%d", *nodeAddress, option.Port)
	return &fullUrl, nil, err
}

func (s *swarmAdapter) findNodeAddress(ctx context.Context, service *swarm.Service) (*string, error) {
	tasks, err := s.client.TaskList(ctx, swarm.TaskListOptions{
		Filters: filters.NewArgs(filters.Arg("service", service.ID)),
	})
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("unable to resolve node IPs: no running tasks found for service %s", service.ID)
	}

	nodes, err := s.client.NodeList(ctx, swarm.NodeListOptions{})
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("unable to resolve node IPs: no nodes found")
	}

	for _, task := range tasks {
		for _, node := range nodes {
			if node.ID == task.NodeID && task.Status.State == swarm.TaskStateRunning {
				return &node.Status.Addr, nil
			}
		}
	}

	return nil, fmt.Errorf(
		"unable to resolve node IPs: unable to conciliate tasks and nodes for service %s",
		service.ID,
	)
}
