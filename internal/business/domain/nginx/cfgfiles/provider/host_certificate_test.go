package provider

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/binding"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
)

func Test_hostCertificateProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		paths := &Paths{
			Config: "/etc/nginx/",
		}
		id := uuid.New()

		certID := uuid.New()

		ctx := &Context{
			Context: t.Context(),
			Paths:   paths,
			Hosts: []host.Host{
				{
					ID: id,
					Bindings: []binding.Binding{
						{
							Type:          binding.HTTPSBindingType,
							CertificateID: &certID,
						},
					},
				},
			},
			Cfg: newSettings(),
		}

		t.Run("successfully provides certificates", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cert := newCertificate()
			cert.ID = certID
			cert.PublicKey = base64.StdEncoding.EncodeToString([]byte("cert-data"))
			cert.PrivateKey = base64.StdEncoding.EncodeToString([]byte("key-data"))
			cert.CertificationChain = []string{
				base64.StdEncoding.EncodeToString([]byte("chain-data")),
			}

			certificateCmds := certificate.NewMockedCommands(ctrl)
			certificateCmds.EXPECT().
				Get(gomock.Any(), certID).
				AnyTimes().
				Return(cert, nil)

			provider := &hostCertificateProvider{
				certificateCommands: certificateCmds,
			}

			files, err := provider.Provide(ctx)
			assert.NoError(t, err)
			assert.Len(t, files, 1)

			assert.Equal(t, fmt.Sprintf("certificate-%s.pem", certID), files[0].Name)
			assert.Contains(t, files[0].Contents, "-----BEGIN CERTIFICATE-----")
			assert.Contains(
				t,
				files[0].Contents,
				base64.StdEncoding.EncodeToString([]byte("cert-data")),
			)
			assert.Contains(
				t,
				files[0].Contents,
				base64.StdEncoding.EncodeToString([]byte("chain-data")),
			)
			assert.Contains(t, files[0].Contents, "-----BEGIN PRIVATE KEY-----")
			assert.Contains(
				t,
				files[0].Contents,
				base64.StdEncoding.EncodeToString([]byte("key-data")),
			)
		})

		t.Run("returns error when certificateCommands fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			certificateCmds := certificate.NewMockedCommands(ctrl)
			certificateCmds.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

			provider := &hostCertificateProvider{
				certificateCommands: certificateCmds,
			}
			_, err := provider.Provide(ctx)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("deduplicates certificates", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			certificateCmds := certificate.NewMockedCommands(ctrl)
			certificateCmds.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(newCertificate(), nil)

			provider := &hostCertificateProvider{
				certificateCommands: certificateCmds,
			}

			subCtx := &Context{
				Context: t.Context(),
				Paths:   paths,
				Hosts: []host.Host{
					{
						Bindings: []binding.Binding{
							{Type: binding.HTTPSBindingType, CertificateID: &certID},
							{Type: binding.HTTPSBindingType, CertificateID: &certID},
						},
					},
				},
				Cfg: newSettings(),
			}

			files, err := provider.Provide(subCtx)
			assert.NoError(t, err)
			assert.Len(t, files, 1)
		})
	})

	t.Run("PemEncoding", func(t *testing.T) {
		t.Run("convertToPemEncodedCertificateString wraps raw bytes in PEM", func(t *testing.T) {
			raw := []byte("fake-cert")
			encoded := convertToPemEncodedCertificateString(raw)
			assert.Contains(t, encoded, "-----BEGIN CERTIFICATE-----")
			assert.Contains(t, encoded, "-----END CERTIFICATE-----")
		})

		t.Run(
			"convertToPemEncodedCertificateString returns existing PEM as is",
			func(t *testing.T) {
				existing := "-----BEGIN CERTIFICATE-----\nstuff\n-----END CERTIFICATE-----"
				encoded := convertToPemEncodedCertificateString([]byte(existing))
				assert.Equal(t, existing, encoded)
			},
		)

		t.Run("convertToPemEncodedPrivateKeyString wraps raw bytes in PEM", func(t *testing.T) {
			raw := []byte("fake-key")
			encoded := convertToPemEncodedPrivateKeyString(raw)
			assert.Contains(t, encoded, "-----BEGIN PRIVATE KEY-----")
			assert.Contains(t, encoded, "-----END PRIVATE KEY-----")
		})
	})
}
