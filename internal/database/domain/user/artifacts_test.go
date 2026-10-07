package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
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
			TrafficStats: user.ReadOnlyAccessLevel,
		},
		Enabled: true,
		TOTP: user.TOTP{
			Secret:    nil,
			Validated: false,
		},
	}
}

func newAPIToken(usr *user.User) *user.APIToken {
	createdAt := time.Now().UTC().Truncate(time.Second)
	expiration := createdAt.Add(time.Hour * 24)

	return &user.APIToken{
		ID:         uuid.New(),
		UserID:     usr.ID,
		Name:       "test-token-" + uuid.New().String(),
		Expiration: &expiration,
		CreatedAt:  createdAt,
	}
}
