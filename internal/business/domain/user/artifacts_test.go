package user

import (
	"github.com/google/uuid"
)

func newUser() *User {
	return &User{
		ID:       uuid.New(),
		Username: "testuser",
		Name:     "Test User",
		Enabled:  true,
		Permissions: Permissions{
			Hosts:        NoAccessAccessLevel,
			Streams:      NoAccessAccessLevel,
			Certificates: NoAccessAccessLevel,
			Logs:         NoAccessAccessLevel,
			Integrations: NoAccessAccessLevel,
			AccessLists:  NoAccessAccessLevel,
			Settings:     NoAccessAccessLevel,
			Users:        NoAccessAccessLevel,
			NginxServer:  ReadOnlyAccessLevel,
			ExportData:   NoAccessAccessLevel,
			VPNs:         NoAccessAccessLevel,
			Caches:       NoAccessAccessLevel,
			TrafficStats: NoAccessAccessLevel,
		},
	}
}

func newSaveRequest() *SaveRequest {
	return &SaveRequest{
		ID:         uuid.New(),
		Username:   "testuser",
		Name:       "Test User",
		Enabled:    true,
		Password:   new("password123"),
		RemoveTOTP: false,
		Permissions: Permissions{
			Hosts:        NoAccessAccessLevel,
			Streams:      NoAccessAccessLevel,
			Certificates: NoAccessAccessLevel,
			Logs:         NoAccessAccessLevel,
			Integrations: NoAccessAccessLevel,
			AccessLists:  NoAccessAccessLevel,
			Settings:     NoAccessAccessLevel,
			Users:        NoAccessAccessLevel,
			NginxServer:  ReadOnlyAccessLevel,
			ExportData:   NoAccessAccessLevel,
			VPNs:         NoAccessAccessLevel,
			Caches:       NoAccessAccessLevel,
			TrafficStats: NoAccessAccessLevel,
		},
	}
}
