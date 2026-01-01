package nginx

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func TestEndpointAdapter_Hash(t *testing.T) {
	id := uuid.New()
	name := "test"
	domain := "example.com"

	t.Run("generates consistent hash", func(t *testing.T) {
		a := &endpointAdapter{
			vpnID:      id,
			name:       name,
			domainName: &domain,
		}
		assert.Equal(t, id.String()+name+domain, a.Hash())
	})

	t.Run("handles nil domain", func(t *testing.T) {
		a := &endpointAdapter{
			vpnID:      id,
			name:       name,
			domainName: nil,
		}
		assert.Equal(t, id.String()+name, a.Hash())
	})
}

func TestEndpointAdapter_Targets(t *testing.T) {
	domain := "example.com"
	bindings := []binding.Binding{
		{
			IP:   "127.0.0.1",
			Port: 80,
			Type: binding.HTTPBindingType,
		},
		{
			IP:   "127.0.0.1",
			Port: 443,
			Type: binding.HTTPSBindingType,
		},
	}

	t.Run("maps bindings to targets correctly", func(t *testing.T) {
		a := &endpointAdapter{
			domainName: &domain,
			bindings:   bindings,
		}
		targets := a.Targets()

		assert.Len(t, targets, 2)
		assert.Equal(t, vpn.EndpointTarget{
			Host:  domain,
			IP:    "127.0.0.1",
			Port:  80,
			HTTPS: false,
		}, targets[0])
		assert.Equal(t, vpn.EndpointTarget{
			Host:  domain,
			IP:    "127.0.0.1",
			Port:  443,
			HTTPS: true,
		}, targets[1])
	})
}

func TestVpnManager_BuildEndpoints(t *testing.T) {
	ctx := context.Background()
	vpnID := uuid.New()
	globalBindings := []binding.Binding{
		{
			IP:   "10.0.0.1",
			Port: 80,
		},
	}
	hostBindings := []binding.Binding{
		{
			IP:   "192.168.1.1",
			Port: 80,
		},
	}

	settingsCmds := &settings.Commands{
		Get: func(_ context.Context) (*settings.Settings, error) {
			return &settings.Settings{
				GlobalBindings: globalBindings,
			}, nil
		},
	}

	m := newVpnManager(nil, settingsCmds)

	t.Run("uses host bindings when UseGlobalBindings is false", func(t *testing.T) {
		hosts := []host.Host{
			{
				UseGlobalBindings: false,
				Bindings:          hostBindings,
				VPNs: []host.VPN{
					{
						VPNID: vpnID,
						Name:  "test",
						Host:  nil,
					},
				},
				DomainNames: []string{"host.com"},
			},
		}

		endpoints, err := m.buildEndpoints(ctx, hosts)
		assert.NoError(t, err)
		assert.Len(t, endpoints, 1)
		assert.Equal(t, hostBindings, endpoints[0].(*endpointAdapter).bindings)
	})

	t.Run("uses global bindings when UseGlobalBindings is true", func(t *testing.T) {
		hosts := []host.Host{
			{
				UseGlobalBindings: true,
				Bindings:          hostBindings,
				VPNs: []host.VPN{
					{
						VPNID: vpnID,
						Name:  "test",
						Host:  nil,
					},
				},
				DomainNames: []string{"host.com"},
			},
		}

		endpoints, err := m.buildEndpoints(ctx, hosts)
		assert.NoError(t, err)
		assert.Len(t, endpoints, 1)
		assert.Equal(t, globalBindings, endpoints[0].(*endpointAdapter).bindings)
	})

	t.Run("falls back to host domain name when vpn host is missing", func(t *testing.T) {
		hosts := []host.Host{
			{
				DomainNames: []string{"fallback.com"},
				VPNs: []host.VPN{
					{
						VPNID: vpnID,
						Name:  "test",
						Host:  nil,
					},
				},
			},
		}

		endpoints, err := m.buildEndpoints(ctx, hosts)
		assert.NoError(t, err)
		assert.Equal(t, "fallback.com", *endpoints[0].(*endpointAdapter).domainName)
	})
}

func TestVpnManager_StopObsoleteEndpoints(t *testing.T) {
	ctx := context.Background()
	vpnID := uuid.New()
	ep1 := &endpointAdapter{
		vpnID: vpnID,
		name:  "ep1",
	}
	ep2 := &endpointAdapter{
		vpnID: vpnID,
		name:  "ep2",
	}

	t.Run("stops endpoints not present in the new list", func(t *testing.T) {
		stopped := make([]string, 0)
		vpnCmds := &vpn.Commands{
			Stop: func(_ context.Context, e vpn.Endpoint) error {
				stopped = append(stopped, e.SourceName())
				return nil
			},
		}

		m := &vpnManager{
			vpnCommands:      vpnCmds,
			currentEndpoints: []vpn.Endpoint{ep1, ep2},
		}

		err := m.stopObsoleteEndpoints(ctx, []vpn.Endpoint{ep1})
		assert.NoError(t, err)
		assert.Equal(t, []string{"ep2"}, stopped)
	})
}

func TestVpnManager_StartNewEndpoints(t *testing.T) {
	ctx := context.Background()
	vpnID := uuid.New()
	ep1 := &endpointAdapter{
		vpnID: vpnID,
		name:  "ep1",
	}
	ep2 := &endpointAdapter{
		vpnID: vpnID,
		name:  "ep2",
	}

	t.Run("starts only endpoints not currently running", func(t *testing.T) {
		started := make([]string, 0)
		vpnCmds := &vpn.Commands{
			Start: func(_ context.Context, e vpn.Endpoint) error {
				started = append(started, e.SourceName())
				return nil
			},
		}

		m := &vpnManager{
			vpnCommands:      vpnCmds,
			currentEndpoints: []vpn.Endpoint{ep1},
		}

		err := m.startNewEndpoints(ctx, []vpn.Endpoint{ep1, ep2})
		assert.NoError(t, err)
		assert.Equal(t, []string{"ep2"}, started)
	})
}
