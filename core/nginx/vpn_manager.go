package nginx

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type destinationAdapter struct {
	vpnID         uuid.UUID
	name          string
	domainName    string
	streamBinding *stream.Address
	hostBinding   *host.Binding
}

type vpnManager struct {
	vpnCommands      *vpn.Commands
	settingsCommands *settings.Commands
	currentData      []vpn.Destination
}

func newVpnManager(vpnCommands *vpn.Commands, settingsCommands *settings.Commands) *vpnManager {
	return &vpnManager{
		vpnCommands:      vpnCommands,
		settingsCommands: settingsCommands,
		currentData:      make([]vpn.Destination, 0),
	}
}

func (m *vpnManager) start(ctx context.Context, hosts []*host.Host, streams []*stream.Stream) error {
	destinations, err := m.buildDestinations(ctx, hosts, streams)
	if err != nil {
		return err
	}

	for _, destination := range destinations {
		if err := m.vpnCommands.Start(ctx, destination); err != nil {
			return err
		}
	}

	m.currentData = destinations
	return nil
}

func (m *vpnManager) reload(ctx context.Context, hosts []*host.Host, streams []*stream.Stream) error {
	newDestinations, err := m.buildDestinations(ctx, hosts, streams)
	if err != nil {
		return err
	}

	for _, oldDest := range m.currentData {
		found := false
		for _, newDest := range newDestinations {
			if oldDest.VPNID() == newDest.VPNID() && oldDest.Name() == newDest.Name() {
				found = true
				break
			}
		}

		if !found {
			if err := m.vpnCommands.Stop(ctx, oldDest); err != nil {
				return err
			}
		}
	}

	for _, newDest := range newDestinations {
		found := false
		for _, oldDest := range m.currentData {
			if oldDest.VPNID() == newDest.VPNID() && oldDest.Name() == newDest.Name() {
				found = true
				break
			}
		}

		if !found {
			if err := m.vpnCommands.Start(ctx, newDest); err != nil {
				return err
			}
		}
	}

	m.currentData = newDestinations
	return nil
}

func (m *vpnManager) buildDestinations(ctx context.Context, hosts []*host.Host, streams []*stream.Stream) ([]vpn.Destination, error) {
	setts, err := m.settingsCommands.Get(ctx)
	if err != nil {
		return nil, err
	}

	globalBindings := setts.GlobalBindings
	destinations := make([]vpn.Destination, 0)

	for _, h := range hosts {
		for _, vpnEntry := range h.VPNs {
			bindings := h.Bindings
			if h.UseGlobalBindings {
				bindings = globalBindings
			}

			var selectedBinding *host.Binding
			for _, binding := range bindings {
				if binding.Type == host.HttpsBindingType {
					selectedBinding = binding
					break
				}
			}
			if selectedBinding == nil && len(bindings) > 0 {
				selectedBinding = bindings[0]
			}

			if selectedBinding != nil {
				var domainName string
				if len(h.DomainNames) > 0 && h.DomainNames[0] != nil {
					domainName = *h.DomainNames[0]
				}

				destinations = append(destinations, &destinationAdapter{
					vpnID:       vpnEntry.VPNID,
					name:        vpnEntry.Name,
					domainName:  domainName,
					hostBinding: selectedBinding,
				})
			}
		}
	}

	for _, s := range streams {
		destinations = append(destinations, &destinationAdapter{
			streamBinding: &s.Binding,
		})
	}

	return destinations, nil
}

func (m *vpnManager) stop(ctx context.Context) error {
	for _, destination := range m.currentData {
		if err := m.vpnCommands.Stop(ctx, destination); err != nil {
			return err
		}
	}

	return nil
}

func (a *destinationAdapter) VPNID() uuid.UUID {
	return a.vpnID
}

func (a *destinationAdapter) Name() string {
	return a.name
}

func (a *destinationAdapter) DomainName() string {
	return a.domainName
}

func (a *destinationAdapter) IP() string {
	if a.streamBinding != nil {
		return a.streamBinding.Address
	}

	return a.hostBinding.IP
}

func (a *destinationAdapter) Port() int {
	if a.streamBinding != nil {
		return *a.streamBinding.Port
	}

	return a.hostBinding.Port
}

func (a *destinationAdapter) HTTPS() bool {
	return a.hostBinding != nil && a.hostBinding.Type == host.HttpsBindingType
}
