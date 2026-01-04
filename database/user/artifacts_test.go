package user

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/user"
)

func newUser() *user.User {
	return &user.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Username:     "testuser-" + uuid.New().String(),
		PasswordHash: "hash",
		PasswordSalt: "salt",
		Permissions: user.Permissions{
			Hosts:        user.ReadWriteAccessLevel,
			Streams:      user.ReadWriteAccessLevel,
			Certificates: user.ReadWriteAccessLevel,
			Logs:         user.ReadOnlyAccessLevel,
			Integrations: user.ReadWriteAccessLevel,
			AccessLists:  user.ReadWriteAccessLevel,
			Settings:     user.ReadWriteAccessLevel,
			Users:        user.ReadWriteAccessLevel,
			NginxServer:  user.ReadWriteAccessLevel,
			ExportData:   user.ReadOnlyAccessLevel,
			VPNs:         user.ReadWriteAccessLevel,
			Caches:       user.ReadWriteAccessLevel,
		},
		Enabled: true,
	}
}
