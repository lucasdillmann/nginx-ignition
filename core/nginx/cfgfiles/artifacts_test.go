package cfgfiles

import (
	"testing"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
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

func newProviderContext(t *testing.T) *providerContext {
	return &providerContext{
		context: t.Context(),
		paths:   newPaths(),
		supportedFeatures: &SupportedFeatures{
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
