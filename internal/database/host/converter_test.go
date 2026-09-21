package host

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/core/host"
)

func Test_Converter(t *testing.T) {
	t.Run("toDomain", func(t *testing.T) {
		t.Run("successfully converts a complete model to domain", func(t *testing.T) {
			vpnID := uuid.New()
			certID := uuid.New()
			model := &hostModel{
				ID: uuid.New(),
				VPNs: []hostVpnModel{
					{
						VPNID:         vpnID,
						Name:          "VPN 1",
						Host:          new("1.2.3.4"),
						CertificateID: &certID,
						EnableHTTPS:   true,
					},
				},
			}

			domain, err := toDomain(model)

			assert.NoError(t, err)
			assert.NotNil(t, domain)
			assert.Len(t, domain.VPNs, 1)
			assert.Equal(t, vpnID, domain.VPNs[0].VPNID)
			assert.Equal(t, "VPN 1", domain.VPNs[0].Name)
			assert.Equal(t, new("1.2.3.4"), domain.VPNs[0].Host)
			assert.Equal(t, &certID, domain.VPNs[0].CertificateID)
			assert.True(t, domain.VPNs[0].EnableHTTPS)
		})

		t.Run("converts route protocol", func(t *testing.T) {
			model := &hostModel{
				ID: uuid.New(),
				Routes: []hostRouteModel{
					{
						ID:       uuid.New(),
						Protocol: "GRPC",
					},
				},
			}

			domain, err := toDomain(model)

			assert.NoError(t, err)
			assert.Len(t, domain.Routes, 1)
			assert.Equal(t, host.GRPCRouteProtocol, domain.Routes[0].Protocol)
		})

		t.Run("defaults route protocol to HTTP_1_1 when missing", func(t *testing.T) {
			model := &hostModel{
				ID: uuid.New(),
				Routes: []hostRouteModel{
					{
						ID: uuid.New(),
					},
				},
			}

			domain, err := toDomain(model)

			assert.NoError(t, err)
			assert.Len(t, domain.Routes, 1)
			assert.Equal(t, host.HTTP11RouteProtocol, domain.Routes[0].Protocol)
		})
	})

	t.Run("toModel", func(t *testing.T) {
		t.Run("successfully converts a complete domain to model", func(t *testing.T) {
			vpnID := uuid.New()
			certID := uuid.New()
			domain := &host.Host{
				ID: uuid.New(),
				VPNs: []host.VPN{
					{
						VPNID:         vpnID,
						Name:          "VPN 1",
						Host:          new("1.2.3.4"),
						CertificateID: &certID,
						EnableHTTPS:   true,
					},
				},
			}

			model, err := toModel(domain)

			assert.NoError(t, err)
			assert.NotNil(t, model)
			assert.Len(t, model.VPNs, 1)
			assert.Equal(t, vpnID, model.VPNs[0].VPNID)
			assert.Equal(t, "VPN 1", model.VPNs[0].Name)
			assert.Equal(t, new("1.2.3.4"), model.VPNs[0].Host)
			assert.Equal(t, &certID, model.VPNs[0].CertificateID)
			assert.True(t, model.VPNs[0].EnableHTTPS)
		})

		t.Run("converts route protocol", func(t *testing.T) {
			domain := &host.Host{
				ID: uuid.New(),
				Routes: []host.Route{
					{
						ID:       uuid.New(),
						Protocol: host.GRPCRouteProtocol,
					},
				},
			}

			model, err := toModel(domain)

			assert.NoError(t, err)
			assert.Len(t, model.Routes, 1)
			assert.Equal(t, "GRPC", model.Routes[0].Protocol)
		})

		t.Run("defaults route protocol to HTTP_1_1 when missing", func(t *testing.T) {
			domain := &host.Host{
				ID: uuid.New(),
				Routes: []host.Route{
					{
						ID: uuid.New(),
					},
				},
			}

			model, err := toModel(domain)

			assert.NoError(t, err)
			assert.Len(t, model.Routes, 1)
			assert.Equal(t, "HTTP_1_1", model.Routes[0].Protocol)
		})
	})
}
