package settings

func newSettings() *Settings {
	return &Settings{
		Nginx: &NginxSettings{
			DefaultContentType: "text/html",
			RuntimeUser:        "nginx",
			Timeouts: &NginxTimeoutsSettings{
				Read:      60,
				Send:      60,
				Connect:   60,
				Keepalive: 65,
			},
			WorkerProcesses:   1,
			WorkerConnections: 1024,
			MaximumBodySizeMb: 1,
			Stats: &NginxStatsSettings{
				Enabled:       false,
				Persistent:    false,
				MaximumSizeMB: 16,
			},
		},
		LogRotation: &LogRotationSettings{
			IntervalUnitCount: 1,
			MaximumLines:      1000,
		},
		CertificateAutoRenew: &CertificateAutoRenewSettings{
			IntervalUnitCount: 1,
		},
	}
}
