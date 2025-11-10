package nginx

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type destinationAdapter struct {
	vpnID      uuid.UUID
	name       string
	domainName *string
	binding    *host.Binding
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

func (m *vpnManager) start(ctx context.Context, hosts []*host.Host) error {
	destinations, err := m.buildDestinations(ctx, hosts)
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

func (m *vpnManager) reload(ctx context.Context, hosts []*host.Host) error {
	newDestinations, err := m.buildDestinations(ctx, hosts)
	if err != nil {
		return err
	}

	for _, oldDest := range m.currentData {
		found := false
		for _, newDest := range newDestinations {
			if oldDest.Hash() == newDest.Hash() {
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
			if oldDest.Hash() == newDest.Hash() {
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

func (m *vpnManager) buildDestinations(ctx context.Context, hosts []*host.Host) ([]vpn.Destination, error) {
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

			domainName := vpnEntry.Host
			if (domainName == nil || *domainName == "") && len(h.DomainNames) > 0 && h.DomainNames[0] != nil {
				domainName = h.DomainNames[0]
			}

			for _, binding := range bindings {
				destinations = append(destinations, &destinationAdapter{
					vpnID:      vpnEntry.VPNID,
					name:       vpnEntry.Name,
					domainName: domainName,
					binding:    binding,
				})
			}
		}
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

func (a *destinationAdapter) Hash() string {
	var domainNameStr string
	if a.domainName != nil {
		domainNameStr = *a.domainName
	}

	return a.vpnID.String() + a.name + domainNameStr
}

func (a *destinationAdapter) VPNID() uuid.UUID {
	return a.vpnID
}

func (a *destinationAdapter) SourceName() string {
	return a.name
}

func (a *destinationAdapter) TargetHost() string {
	if a.domainName == nil {
		return ""
	}

	return *a.domainName
}

func (a *destinationAdapter) IP() string {
	return a.binding.IP
}

func (a *destinationAdapter) Port() int {
	return a.binding.Port
}

func (a *destinationAdapter) HTTPS() bool {
	return a.binding.Type == host.HttpsBindingType
}
