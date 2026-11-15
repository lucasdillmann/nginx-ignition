package provider

import (
	"bytes"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/oschwald/maxminddb-golang"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/binding"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
)

func newPaths() *Paths {
	return &Paths{
		Base:   "/",
		Config: "/etc/nginx/",
		Logs:   "/var/log/nginx/",
		Cache:  "/var/cache/nginx/",
		Temp:   "/tmp/nginx/",
	}
}

func newProviderContext(t *testing.T) *Context {
	return &Context{
		Context: t.Context(),
		Paths:   newPaths(),
		SupportedFeatures: &SupportedFeatures{
			TLSSNI:      StaticSupportType,
			StreamType:  StaticSupportType,
			RunCodeType: StaticSupportType,
		},
	}
}

func newAccessList() accesslist.AccessList {
	return accesslist.AccessList{
		ID:             uuid.New(),
		DefaultOutcome: accesslist.DenyOutcome,
		Credentials: []accesslist.Credentials{
			{
				Username: "user",
				Password: "pwd",
			},
		},
	}
}

func newCertificate() *certificate.Certificate {
	return &certificate.Certificate{
		ID:         uuid.New(),
		PublicKey:  "cert-data",
		PrivateKey: "key-data",
		CertificationChain: []string{
			"chain-data",
		},
	}
}

func newSettings() *settings.Settings {
	return &settings.Settings{
		Nginx: &settings.NginxSettings{
			RuntimeUser:       "nginx",
			WorkerProcesses:   1,
			WorkerConnections: 1024,
			Timeouts: &settings.NginxTimeoutsSettings{
				Keepalive:  65,
				Connect:    60,
				Read:       60,
				Send:       60,
				ClientBody: 60,
			},
			Buffers: &settings.NginxBuffersSettings{
				ClientBodyKb:   16,
				ClientHeaderKb: 1,
				LargeClientHeader: &settings.NginxBufferSize{
					Amount: 4,
					SizeKb: 8,
				},
				Output: &settings.NginxBufferSize{
					Amount: 2,
					SizeKb: 32,
				},
			},
			Logs: &settings.NginxLogsSettings{
				ServerLogsEnabled: true,
				ServerLogsLevel:   settings.WarnLogLevel,
			},
			Stats: &settings.NginxStatsSettings{
				Enabled: false,
			},
		},
	}
}

func newCache() cache.Cache {
	return cache.Cache{
		ID:              uuid.New(),
		InactiveSeconds: new(3600),
		MaximumSizeMB:   new(1024),
	}
}

func newHost() host.Host {
	return host.Host{
		ID:            uuid.New(),
		Enabled:       true,
		DefaultServer: true,
		DomainNames:   []string{"example.com"},
		Bindings: []binding.Binding{
			{
				Type: binding.HTTPBindingType,
				IP:   "0.0.0.0",
				Port: 80,
			},
		},
	}
}

func newStream() stream.Stream {
	return stream.Stream{
		ID: uuid.New(),
		Binding: stream.Address{
			Protocol: stream.TCPProtocol,
			Address:  "0.0.0.0",
			Port:     new(80),
		},
		Type: stream.SimpleType,
		DefaultBackend: stream.Backend{
			Address: stream.Address{
				Protocol: stream.TCPProtocol,
				Address:  "127.0.0.1",
				Port:     new(8080),
			},
		},
	}
}

func newGeoIPDatabase(t *testing.T) []byte {
	t.Helper()

	tree, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType: "GeoLite2-Test",
		Description:  map[string]string{"en": "Test GeoIP database"},
		IPVersion:    4,
		RecordSize:   24,
		Languages:    []string{"en"},
	})
	require.NoError(t, err)

	_, network, err := net.ParseCIDR("8.8.8.0/24")
	require.NoError(t, err)

	err = tree.Insert(network, mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code": mmdbtype.String("US"),
		},
	})
	require.NoError(t, err)

	buffer := new(bytes.Buffer)
	_, err = tree.WriteTo(buffer)
	require.NoError(t, err)

	return buffer.Bytes()
}

func corruptGeoIPDataSection(t *testing.T, data []byte) []byte {
	t.Helper()

	reader, err := maxminddb.FromBytes(data)
	require.NoError(t, err)

	searchTreeSize := int(reader.Metadata.NodeCount * reader.Metadata.RecordSize / 4)
	_ = reader.Close()

	corrupted := make([]byte, len(data))
	copy(corrupted, data)

	dataSectionStart := searchTreeSize + 16
	require.Less(t, dataSectionStart+1, len(corrupted))
	corrupted[dataSectionStart] = 0x00
	corrupted[dataSectionStart+1] = 0x09

	return corrupted
}

func newGeoIPTestContext(t *testing.T) *Context {
	t.Helper()

	ctx := newProviderContext(t)
	ctx.Cfg = newSettings()
	ctx.Cfg.Nginx.Stats.Enabled = true

	return ctx
}
