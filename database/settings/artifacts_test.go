package settings

import (
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func newSettings() *settings.Settings {
	return &settings.Settings{
		Nginx: &settings.NginxSettings{
			Timeouts: &settings.NginxTimeoutsSettings{
				Read:       60,
				Connect:    60,
				Send:       60,
				Keepalive:  75,
				ClientBody: 60,
			},
			Buffers: &settings.NginxBuffersSettings{
				LargeClientHeader: &settings.NginxBufferSize{
					SizeKb: 8,
					Amount: 4,
				},
				Output: &settings.NginxBufferSize{
					SizeKb: 32,
					Amount: 4,
				},
				ClientBodyKb:   16,
				ClientHeaderKb: 1,
			},
			Logs: &settings.NginxLogsSettings{
				ServerLogsLevel:   settings.WarnLogLevel,
				ErrorLogsLevel:    settings.ErrorLogLevel,
				ServerLogsEnabled: true,
				AccessLogsEnabled: true,
				ErrorLogsEnabled:  true,
			},
			Custom:             ptr.Of("# Custom Nginx Config"),
			RuntimeUser:        "nginx",
			DefaultContentType: "text/html",
			WorkerProcesses:    1,
			WorkerConnections:  1024,
			MaximumBodySizeMb:  10,
			Stats: &settings.NginxStatsSettings{
				Enabled:       false,
				Persistent:    false,
				AllHosts:      false,
				MaximumSizeMB: 16,
			},
			ServerTokensEnabled: false,
			TCPNoDelayEnabled:   true,
			GzipEnabled:         true,
			SendfileEnabled:     true,
		},
		LogRotation: &settings.LogRotationSettings{
			IntervalUnit:      settings.DaysTimeUnit,
			MaximumLines:      10000,
			IntervalUnitCount: 7,
			Enabled:           true,
		},
		CertificateAutoRenew: &settings.CertificateAutoRenewSettings{
			IntervalUnit:      settings.DaysTimeUnit,
			IntervalUnitCount: 30,
			Enabled:           true,
		},
		GlobalBindings: []binding.Binding{
			{
				Type: binding.HTTPBindingType,
				IP:   "0.0.0.0",
				Port: 80,
			},
		},
	}
}
