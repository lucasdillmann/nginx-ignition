package cfgfiles

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_HostCertificateFileProvider(t *testing.T) {
	t.Run("Provide", func(t *testing.T) {
		paths := &Paths{
			Config: "/etc/nginx/",
		}
		id := uuid.New()

		certID := uuid.New()

		ctx := &providerContext{
			context: context.Background(),
			paths:   paths,
			hosts: []host.Host{
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
		}

		t.Run("successfully provides certificates", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).AnyTimes().Return(&settings.Settings{
				Nginx: &settings.NginxSettings{
					Logs: &settings.NginxLogsSettings{},
				},
			}, nil)

			certificateCmds := certificate.NewMockedCommands(ctrl)
			certificateCmds.EXPECT().
				Get(gomock.Any(), certID).
				AnyTimes().
				Return(&certificate.Certificate{
					ID:         certID,
					PublicKey:  base64.StdEncoding.EncodeToString([]byte("cert-data")),
					PrivateKey: base64.StdEncoding.EncodeToString([]byte("key-data")),
					CertificationChain: []string{
						base64.StdEncoding.EncodeToString([]byte("chain-data")),
					},
				}, nil)

			p := &hostCertificateFileProvider{
				settingsCommands:    settingsCmds,
				certificateCommands: certificateCmds,
			}

			files, err := p.provide(ctx)
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

		t.Run("returns error when settingsCommands fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().Get(gomock.Any()).Return(nil, assert.AnError)

			p := &hostCertificateFileProvider{settingsCommands: settingsCmds}
			_, err := p.provide(ctx)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("returns error when certificateCommands fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().
				Get(gomock.Any()).
				Return(&settings.Settings{Nginx: &settings.NginxSettings{}}, nil)

			certificateCmds := certificate.NewMockedCommands(ctrl)
			certificateCmds.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

			p := &hostCertificateFileProvider{
				settingsCommands:    settingsCmds,
				certificateCommands: certificateCmds,
			}
			_, err := p.provide(ctx)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("deduplicates certificates", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			settingsCmds := settings.NewMockedCommands(ctrl)
			settingsCmds.EXPECT().
				Get(gomock.Any()).
				AnyTimes().
				Return(&settings.Settings{Nginx: &settings.NginxSettings{}}, nil)

			certificateCmds := certificate.NewMockedCommands(ctrl)
			certificateCmds.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(&certificate.Certificate{
					PublicKey:  "data",
					PrivateKey: "data",
				}, nil)

			p := &hostCertificateFileProvider{
				settingsCommands:    settingsCmds,
				certificateCommands: certificateCmds,
			}

			subCtx := &providerContext{
				context: context.Background(),
				paths:   paths,
				hosts: []host.Host{
					{
						Bindings: []binding.Binding{
							{Type: binding.HTTPSBindingType, CertificateID: &certID},
							{Type: binding.HTTPSBindingType, CertificateID: &certID},
						},
					},
				},
			}

			files, err := p.provide(subCtx)
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
