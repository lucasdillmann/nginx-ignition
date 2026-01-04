package user

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
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
		},
	}
}

func newSaveRequest() *SaveRequest {
	return &SaveRequest{
		ID:       uuid.New(),
		Username: "testuser",
		Name:     "Test User",
		Enabled:  true,
		Password: ptr.Of("password123"),
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
		},
	}
}
