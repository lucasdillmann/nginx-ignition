package user

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/user"
)

func newUser() *user.User {
	return &user.User{
		ID:       uuid.New(),
		Name:     "Test User",
		Username: "testuser",
		Enabled:  true,
		Permissions: user.Permissions{
			Hosts:        user.ReadWriteAccessLevel,
			Streams:      user.ReadWriteAccessLevel,
			Certificates: user.ReadWriteAccessLevel,
			Logs:         user.ReadWriteAccessLevel,
			Integrations: user.ReadWriteAccessLevel,
			AccessLists:  user.ReadWriteAccessLevel,
			Settings:     user.ReadWriteAccessLevel,
			Users:        user.ReadWriteAccessLevel,
			NginxServer:  user.ReadWriteAccessLevel,
			ExportData:   user.ReadWriteAccessLevel,
			VPNs:         user.ReadWriteAccessLevel,
			Caches:       user.ReadWriteAccessLevel,
		},
	}
}

func newUserRequest() userRequestDTO {
	return userRequestDTO{
		Name:       new("Test User"),
		Username:   new("testuser"),
		Password:   new("password123"),
		Enabled:    new(true),
		RemoveTOTP: new(false),
		Permissions: userPermissionsDTO{
			Hosts:        string(user.ReadWriteAccessLevel),
			Streams:      string(user.ReadWriteAccessLevel),
			Certificates: string(user.ReadWriteAccessLevel),
			Logs:         string(user.ReadWriteAccessLevel),
			Integrations: string(user.ReadWriteAccessLevel),
			AccessLists:  string(user.ReadWriteAccessLevel),
			Settings:     string(user.ReadWriteAccessLevel),
			Users:        string(user.ReadWriteAccessLevel),
			NginxServer:  string(user.ReadWriteAccessLevel),
			ExportData:   string(user.ReadWriteAccessLevel),
			VPNs:         string(user.ReadWriteAccessLevel),
			Caches:       string(user.ReadWriteAccessLevel),
		},
	}
}

func newUserPage() *pagination.Page[user.User] {
	return pagination.Of([]user.User{
		*newUser(),
	})
}
