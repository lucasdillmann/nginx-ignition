package cfgfiles

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func TestHostCertificateFileProvider_Provide(t *testing.T) {
	p := &hostCertificateFileProvider{}
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

	p.settingsCommands = &settings.Commands{
		Get: func(_ context.Context) (*settings.Settings, error) {
			return &settings.Settings{
				Nginx: &settings.NginxSettings{
					Logs: &settings.NginxLogsSettings{},
				},
			}, nil
		},
	}

	p.certificateCommands = &certificate.Commands{
		Get: func(_ context.Context, _ uuid.UUID) (*certificate.Certificate, error) {
			return &certificate.Certificate{
				ID:         certID,
				PublicKey:  base64.StdEncoding.EncodeToString([]byte("cert-data")),
				PrivateKey: base64.StdEncoding.EncodeToString([]byte("key-data")),
				CertificationChain: []string{
					base64.StdEncoding.EncodeToString([]byte("chain-data")),
				},
			}, nil
		},
	}

	files, err := p.provide(ctx)
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	assert.Equal(t, fmt.Sprintf("certificate-%s.pem", certID), files[0].Name)
	assert.Contains(t, files[0].Contents, "-----BEGIN CERTIFICATE-----")
	assert.Contains(t, files[0].Contents, base64.StdEncoding.EncodeToString([]byte("cert-data")))
	assert.Contains(t, files[0].Contents, base64.StdEncoding.EncodeToString([]byte("chain-data")))
	assert.Contains(t, files[0].Contents, "-----BEGIN PRIVATE KEY-----")
	assert.Contains(t, files[0].Contents, base64.StdEncoding.EncodeToString([]byte("key-data")))
}

func TestHostCertificateFileProvider_PemEncoding(t *testing.T) {
	t.Run("convertToPemEncodedCertificateString wraps raw bytes in PEM", func(t *testing.T) {
		raw := []byte("fake-cert")
		encoded := convertToPemEncodedCertificateString(raw)
		assert.Contains(t, encoded, "-----BEGIN CERTIFICATE-----")
		assert.Contains(t, encoded, "-----END CERTIFICATE-----")
	})

	t.Run("convertToPemEncodedCertificateString returns existing PEM as is", func(t *testing.T) {
		existing := "-----BEGIN CERTIFICATE-----\nstuff\n-----END CERTIFICATE-----"
		encoded := convertToPemEncodedCertificateString([]byte(existing))
		assert.Equal(t, existing, encoded)
	})

	t.Run("convertToPemEncodedPrivateKeyString wraps raw bytes in PEM", func(t *testing.T) {
		raw := []byte("fake-key")
		encoded := convertToPemEncodedPrivateKeyString(raw)
		assert.Contains(t, encoded, "-----BEGIN PRIVATE KEY-----")
		assert.Contains(t, encoded, "-----END PRIVATE KEY-----")
	})
}
