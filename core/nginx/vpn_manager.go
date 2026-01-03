package nginx

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type endpointAdapter struct {
	domainName *string
	name       string
	bindings   []binding.Binding
	vpnID      uuid.UUID
}

type vpnManager struct {
	vpnCommands      vpn.Commands
	settingsCommands settings.Commands
	currentEndpoints []vpn.Endpoint
}

func newVpnManager(vpnCommands vpn.Commands, settingsCommands settings.Commands) *vpnManager {
	return &vpnManager{
		vpnCommands:      vpnCommands,
		settingsCommands: settingsCommands,
		currentEndpoints: make([]vpn.Endpoint, 0),
	}
}

func (m *vpnManager) start(ctx context.Context, hosts []host.Host) error {
	endpoints, err := m.buildEndpoints(ctx, hosts)
	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		if err := m.vpnCommands.Start(ctx, endpoint); err != nil {
			return err
		}
	}

	m.currentEndpoints = endpoints
	return nil
}

func (m *vpnManager) reload(ctx context.Context, hosts []host.Host) error {
	newEndpoints, err := m.buildEndpoints(ctx, hosts)
	if err != nil {
		return err
	}

	if err := m.stopObsoleteEndpoints(ctx, newEndpoints); err != nil {
		return err
	}

	if err := m.startNewEndpoints(ctx, newEndpoints); err != nil {
		return err
	}

	m.currentEndpoints = newEndpoints
	return nil
}

func (m *vpnManager) stopObsoleteEndpoints(ctx context.Context, newEndpoints []vpn.Endpoint) error {
	for _, oldEndpoint := range m.currentEndpoints {
		found := false
		for _, newEndpoint := range newEndpoints {
			if oldEndpoint.Hash() == newEndpoint.Hash() {
				found = true
				break
			}
		}

		if !found {
			if err := m.vpnCommands.Stop(ctx, oldEndpoint); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *vpnManager) startNewEndpoints(ctx context.Context, newEndpoints []vpn.Endpoint) error {
	for _, newEndpoint := range newEndpoints {
		found := false
		for _, oldDest := range m.currentEndpoints {
			if oldDest.Hash() == newEndpoint.Hash() {
				found = true
				break
			}
		}

		if !found {
			if err := m.vpnCommands.Start(ctx, newEndpoint); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *vpnManager) buildEndpoints(
	ctx context.Context,
	hosts []host.Host,
) ([]vpn.Endpoint, error) {
	setts, err := m.settingsCommands.Get(ctx)
	if err != nil {
		return nil, err
	}

	globalBindings := setts.GlobalBindings
	endpoints := make([]vpn.Endpoint, 0)

	for _, h := range hosts {
		for _, vpnEntry := range h.VPNs {
			bindings := h.Bindings
			if h.UseGlobalBindings {
				bindings = globalBindings
			}

			domainName := vpnEntry.Host
			if (domainName == nil || *domainName == "") && len(h.DomainNames) > 0 {
				domainName = &h.DomainNames[0]
			}

			endpoints = append(endpoints, &endpointAdapter{
				vpnID:      vpnEntry.VPNID,
				name:       vpnEntry.Name,
				domainName: domainName,
				bindings:   bindings,
			})
		}
	}

	return endpoints, nil
}

func (m *vpnManager) stop(ctx context.Context) error {
	for _, endpoint := range m.currentEndpoints {
		if err := m.vpnCommands.Stop(ctx, endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (a *endpointAdapter) Hash() string {
	var domainNameStr string
	if a.domainName != nil {
		domainNameStr = *a.domainName
	}

	return a.vpnID.String() + a.name + domainNameStr
}

func (a *endpointAdapter) VPNID() uuid.UUID {
	return a.vpnID
}

func (a *endpointAdapter) SourceName() string {
	return a.name
}

func (a *endpointAdapter) Targets() []vpn.EndpointTarget {
	var targetHost string
	if a.domainName != nil {
		targetHost = *a.domainName
	}

	output := make([]vpn.EndpointTarget, len(a.bindings))
	for index, b := range a.bindings {
		output[index].Host = targetHost
		output[index].IP = b.IP
		output[index].Port = b.Port
		output[index].HTTPS = b.Type == binding.HTTPSBindingType
	}

	return output
}
