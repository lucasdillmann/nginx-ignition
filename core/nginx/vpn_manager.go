package nginx

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/host"
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
	commands    *vpn.Commands
	currentData []vpn.Destination
}

func newVpnManager(commands *vpn.Commands) *vpnManager {
	return &vpnManager{
		commands:    commands,
		currentData: make([]vpn.Destination, 0),
	}
}

func (m *vpnManager) start(ctx context.Context, hosts []*host.Host, streams []*stream.Stream) error {
	// TODO: Implement this
	return nil
}

func (m *vpnManager) reload(ctx context.Context, hosts []*host.Host, streams []*stream.Stream) error {
	// TODO: Implement this
	return nil
}

func (m *vpnManager) stop(ctx context.Context) error {
	for _, destination := range m.currentData {
		if err := m.commands.Stop(ctx, destination); err != nil {
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
