package settings

import (
	"dillmann.com.br/nginx-ignition/core/common/ptr"
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
			WorkerProcesses:    ptr.Of(1),
			WorkerConnections:  ptr.Of(1024),
			DefaultContentType: ptr.Of("text/plain"),
			Stats: &nginxStatsSettingsDTO{
				Enabled:       ptr.Of(true),
				Persistent:    ptr.Of(true),
				MaximumSizeMB: ptr.Of(16),
			},
			ServerTokensEnabled: ptr.Of(true),
			MaximumBodySizeMb:   ptr.Of(10),
			SendfileEnabled:     ptr.Of(true),
			GzipEnabled:         ptr.Of(true),
			TCPNoDelayEnabled:   ptr.Of(true),
			RuntimeUser:         ptr.Of("nginx"),
			Logs: &nginxLogsSettingsDTO{
				ServerLogsEnabled: ptr.Of(true),
				ServerLogsLevel:   ptr.Of(settings.WarnLogLevel),
				AccessLogsEnabled: ptr.Of(true),
				ErrorLogsEnabled:  ptr.Of(true),
				ErrorLogsLevel:    ptr.Of(settings.ErrorLogLevel),
			},
			Timeouts: &nginxTimeoutsSettingsDTO{
				Read:       ptr.Of(60),
				Connect:    ptr.Of(60),
				Send:       ptr.Of(60),
				Keepalive:  ptr.Of(60),
				ClientBody: ptr.Of(60),
			},
			Buffers: &nginxBuffersSettingsDTO{
				ClientBodyKb:   ptr.Of(8),
				ClientHeaderKb: ptr.Of(1),
				LargeClientHeader: &nginxBufferSizeDTO{
					SizeKb: ptr.Of(8),
					Amount: ptr.Of(4),
				},
				Output: &nginxBufferSizeDTO{
					SizeKb: ptr.Of(32),
					Amount: ptr.Of(1),
				},
			},
		},
		LogRotation: &logRotationSettingsDTO{
			Enabled:           ptr.Of(true),
			MaximumLines:      ptr.Of(1000),
			IntervalUnit:      ptr.Of(settings.DaysTimeUnit),
			IntervalUnitCount: ptr.Of(1),
		},
		CertificateAutoRenew: &certificateAutoRenewSettingsDTO{
			Enabled:           ptr.Of(true),
			IntervalUnit:      ptr.Of(settings.DaysTimeUnit),
			IntervalUnitCount: ptr.Of(30),
		},
	}
}
