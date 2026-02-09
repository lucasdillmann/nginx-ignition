package nginx

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func newMetadata() *nginx.Metadata {
	return &nginx.Metadata{
		Version:      "1.21.0",
		BuildDetails: "nginx version: nginx/1.21.0",
		Modules:      []string{"http_ssl_module", "http_v2_module"},
	}
}

func newTrafficStats() *nginx.Stats {
	return &nginx.Stats{
		HostName: "test-host",
		Connections: nginx.StatsConnections{
			Active:   10,
			Reading:  2,
			Writing:  3,
			Waiting:  5,
			Accepted: 100,
			Handled:  99,
			Requests: 150,
		},
		ServerZones: map[string]nginx.StatsZoneData{
			"zone1": {
				RequestCounter: 50,
				InBytes:        1024,
				OutBytes:       2048,
			},
		},
		FilterZones:   make(map[string]map[string]nginx.StatsZoneData),
		UpstreamZones: make(map[string][]nginx.StatsUpstreamZoneData),
	}
}

func newSettings() *settings.Settings {
	return &settings.Settings{
		Nginx: &settings.NginxSettings{
			Stats: &settings.NginxStatsSettings{
				Enabled:  true,
				AllHosts: false,
			},
		},
	}
}
