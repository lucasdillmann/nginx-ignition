package settings

import (
	"dillmann.com.br/nginx-ignition/core/settings"
)

func newSettings() *settings.Settings {
	return &settings.Settings{
		Nginx: &settings.NginxSettings{
			WorkerProcesses:    1,
			WorkerConnections:  1024,
			DefaultContentType: "text/plain",
			Stats: &settings.NginxStatsSettings{
				Enabled:       true,
				Persistent:    true,
				AllHosts:      true,
				MaximumSizeMB: 16,
			},
			ServerTokensEnabled: true,
			MaximumBodySizeMb:   10,
			SendfileEnabled:     true,
			GzipEnabled:         true,
			TCPNoDelayEnabled:   true,
			RuntimeUser:         "nginx",
			Logs: &settings.NginxLogsSettings{
				ServerLogsEnabled: true,
				ServerLogsLevel:   settings.WarnLogLevel,
				AccessLogsEnabled: true,
				ErrorLogsEnabled:  true,
				ErrorLogsLevel:    settings.ErrorLogLevel,
			},
			Timeouts: &settings.NginxTimeoutsSettings{
				Read:       60,
				Connect:    60,
				Send:       60,
				Keepalive:  60,
				ClientBody: 60,
			},
			Buffers: &settings.NginxBuffersSettings{
				ClientBodyKb:   8,
				ClientHeaderKb: 1,
				LargeClientHeader: &settings.NginxBufferSize{
					SizeKb: 8,
					Amount: 4,
				},
				Output: &settings.NginxBufferSize{
					SizeKb: 32,
					Amount: 1,
				},
			},
		},
		LogRotation: &settings.LogRotationSettings{
			Enabled:           true,
			MaximumLines:      1000,
			IntervalUnit:      settings.DaysTimeUnit,
			IntervalUnitCount: 1,
		},
		CertificateAutoRenew: &settings.CertificateAutoRenewSettings{
			Enabled:           true,
			IntervalUnit:      settings.DaysTimeUnit,
			IntervalUnitCount: 30,
		},
	}
}

func newSettingsDTO() *settingsDTO {
	return &settingsDTO{
		Nginx: &nginxSettingsDTO{
			WorkerProcesses:    new(1),
			WorkerConnections:  new(1024),
			DefaultContentType: new("text/plain"),
			Stats: &nginxStatsSettingsDTO{
				Enabled:       new(true),
				Persistent:    new(true),
				AllHosts:      new(true),
				MaximumSizeMB: new(16),
			},
			ServerTokensEnabled: new(true),
			MaximumBodySizeMb:   new(10),
			SendfileEnabled:     new(true),
			GzipEnabled:         new(true),
			TCPNoDelayEnabled:   new(true),
			RuntimeUser:         new("nginx"),
			Logs: &nginxLogsSettingsDTO{
				ServerLogsEnabled: new(true),
				ServerLogsLevel:   new(settings.WarnLogLevel),
				AccessLogsEnabled: new(true),
				ErrorLogsEnabled:  new(true),
				ErrorLogsLevel:    new(settings.ErrorLogLevel),
			},
			Timeouts: &nginxTimeoutsSettingsDTO{
				Read:       new(60),
				Connect:    new(60),
				Send:       new(60),
				Keepalive:  new(60),
				ClientBody: new(60),
			},
			Buffers: &nginxBuffersSettingsDTO{
				ClientBodyKb:   new(8),
				ClientHeaderKb: new(1),
				LargeClientHeader: &nginxBufferSizeDTO{
					SizeKb: new(8),
					Amount: new(4),
				},
				Output: &nginxBufferSizeDTO{
					SizeKb: new(32),
					Amount: new(1),
				},
			},
		},
		LogRotation: &logRotationSettingsDTO{
			Enabled:           new(true),
			MaximumLines:      new(1000),
			IntervalUnit:      new(settings.DaysTimeUnit),
			IntervalUnitCount: new(1),
		},
		CertificateAutoRenew: &certificateAutoRenewSettingsDTO{
			Enabled:           new(true),
			IntervalUnit:      new(settings.DaysTimeUnit),
			IntervalUnitCount: new(30),
		},
	}
}
