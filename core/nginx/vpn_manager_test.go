package nginx

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_endpointAdapter(t *testing.T) {
	t.Run("Hash", func(t *testing.T) {
		id := uuid.New()
		name := "test"
		domain := "example.com"

		t.Run("generates consistent hash", func(t *testing.T) {
			adapter := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  &domain,
				enableHTTPS: false,
			}
			assert.Equal(t, id.String()+name+domain+"false", adapter.Hash())
		})

		t.Run("handles nil domain", func(t *testing.T) {
			adapter := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  nil,
				enableHTTPS: false,
			}
			assert.Equal(t, id.String()+name+"false", adapter.Hash())
		})

		t.Run("generates different hash when enableHTTPS changes", func(t *testing.T) {
			adapterHTTPSOn := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  &domain,
				enableHTTPS: true,
			}
			adapterHTTPSOff := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  &domain,
				enableHTTPS: false,
			}
			assert.NotEqual(t, adapterHTTPSOn.Hash(), adapterHTTPSOff.Hash())
		})

		t.Run("generates different hash when certificate changes", func(t *testing.T) {
			certID1 := uuid.New()
			certID2 := uuid.New()
			adapterCert1 := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  &domain,
				enableHTTPS: true,
				certDetails: &certificate.Certificate{ID: certID1},
			}
			adapterCert2 := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  &domain,
				enableHTTPS: true,
				certDetails: &certificate.Certificate{ID: certID2},
			}
			assert.NotEqual(t, adapterCert1.Hash(), adapterCert2.Hash())
		})

		t.Run("generates same hash when nothing changes", func(t *testing.T) {
			certID := uuid.New()
			adapter1 := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  &domain,
				enableHTTPS: true,
				certDetails: &certificate.Certificate{ID: certID},
			}
			adapter2 := &endpointAdapter{
				vpnID:       id,
				name:        name,
				domainName:  &domain,
				enableHTTPS: true,
				certDetails: &certificate.Certificate{ID: certID},
			}
			assert.Equal(t, adapter1.Hash(), adapter2.Hash())
		})
	})

	t.Run("Targets", func(t *testing.T) {
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
			adapter := &endpointAdapter{
				domainName:  &domain,
				bindings:    bindings,
				enableHTTPS: true,
			}
			targets := adapter.Targets()

			assert.Len(t, targets, 2)
			assert.Equal(t, vpn.EndpointTarget{
				Host:  domain,
				IP:    "127.0.0.1",
				Port:  80,
				HTTPS: vpn.EndpointHTTPS{},
			}, targets[0])
			assert.Equal(t, vpn.EndpointTarget{
				Host: domain,
				IP:   "127.0.0.1",
				Port: 443,
				HTTPS: vpn.EndpointHTTPS{
					Enabled: true,
				},
			}, targets[1])
		})
	})
}

func Test_vpnManager(t *testing.T) {
	t.Run("buildEndpoints", func(t *testing.T) {
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

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(t.Context()).AnyTimes().Return(&settings.Settings{
			GlobalBindings: globalBindings,
		}, nil)

		certCommands := certificate.NewMockedCommands(ctrl)

		manager := newVpnManager(nil, settingsCmds, certCommands)

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

			endpoints, err := manager.buildEndpoints(t.Context(), hosts)
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

			endpoints, err := manager.buildEndpoints(t.Context(), hosts)
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

			endpoints, err := manager.buildEndpoints(t.Context(), hosts)
			assert.NoError(t, err)
			assert.Equal(t, "fallback.com", *endpoints[0].(*endpointAdapter).domainName)
		})
	})

	t.Run("stopObsoleteEndpoints", func(t *testing.T) {
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpnCmds := vpn.NewMockedCommands(ctrl)
			vpnCmds.EXPECT().Stop(t.Context(), ep2).Return(nil)

			manager := &vpnManager{
				vpnCommands:      vpnCmds,
				currentEndpoints: []vpn.Endpoint{ep1, ep2},
			}

			err := manager.stopObsoleteEndpoints(t.Context(), []vpn.Endpoint{ep1})
			assert.NoError(t, err)
		})
	})

	t.Run("startNewEndpoints", func(t *testing.T) {
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpnCmds := vpn.NewMockedCommands(ctrl)
			vpnCmds.EXPECT().Start(t.Context(), ep2).Return(nil)

			manager := &vpnManager{
				vpnCommands:      vpnCmds,
				currentEndpoints: []vpn.Endpoint{ep1},
			}

			err := manager.startNewEndpoints(t.Context(), []vpn.Endpoint{ep1, ep2})
			assert.NoError(t, err)
		})
	})
}
