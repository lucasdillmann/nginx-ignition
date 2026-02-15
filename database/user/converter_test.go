package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_Converter(t *testing.T) {
	t.Run("toDomain", func(t *testing.T) {
		t.Run("successfully converts a complete model to domain", func(t *testing.T) {
			model := &userModel{
				ID:                      uuid.New(),
				Enabled:                 true,
				Name:                    "Name",
				Username:                "username",
				PasswordHash:            "hash",
				PasswordSalt:            "salt",
				HostsAccessLevel:        "READ_WRITE",
				StreamsAccessLevel:      "READ_WRITE",
				CertificatesAccessLevel: "READ_WRITE",
				LogsAccessLevel:         "READ_ONLY",
				IntegrationsAccessLevel: "READ_WRITE",
				AccessListsAccessLevel:  "READ_WRITE",
				SettingsAccessLevel:     "READ_WRITE",
				UsersAccessLevel:        "READ_WRITE",
				NginxServerAccessLevel:  "READ_WRITE",
				ExportDataAccessLevel:   "READ_ONLY",
				VPNsAccessLevel:         "READ_WRITE",
				CachesAccessLevel:       "READ_WRITE",
				TrafficStatsAccessLevel: "READ_ONLY",
				TotpSecret:              new("secret"),
				TotpValidated:           true,
			}

			domain := toDomain(model)

			assert.Equal(t, model.ID, domain.ID)
			assert.Equal(t, model.Enabled, domain.Enabled)
			assert.Equal(t, model.Name, domain.Name)
			assert.Equal(t, model.Username, domain.Username)
			assert.Equal(t, model.PasswordHash, domain.PasswordHash)
			assert.Equal(t, model.PasswordSalt, domain.PasswordSalt)
			assert.Equal(t, user.AccessLevel(model.HostsAccessLevel), domain.Permissions.Hosts)
			assert.Equal(t, user.AccessLevel(model.StreamsAccessLevel), domain.Permissions.Streams)
			assert.Equal(
				t,
				user.AccessLevel(model.CertificatesAccessLevel),
				domain.Permissions.Certificates,
			)
			assert.Equal(t, user.AccessLevel(model.LogsAccessLevel), domain.Permissions.Logs)
			assert.Equal(
				t,
				user.AccessLevel(model.IntegrationsAccessLevel),
				domain.Permissions.Integrations,
			)
			assert.Equal(
				t,
				user.AccessLevel(model.AccessListsAccessLevel),
				domain.Permissions.AccessLists,
			)
			assert.Equal(
				t,
				user.AccessLevel(model.SettingsAccessLevel),
				domain.Permissions.Settings,
			)
			assert.Equal(t, user.AccessLevel(model.UsersAccessLevel), domain.Permissions.Users)
			assert.Equal(
				t,
				user.AccessLevel(model.NginxServerAccessLevel),
				domain.Permissions.NginxServer,
			)
			assert.Equal(
				t,
				user.AccessLevel(model.ExportDataAccessLevel),
				domain.Permissions.ExportData,
			)
			assert.Equal(t, user.AccessLevel(model.VPNsAccessLevel), domain.Permissions.VPNs)
			assert.Equal(t, user.AccessLevel(model.CachesAccessLevel), domain.Permissions.Caches)
			assert.Equal(
				t,
				user.AccessLevel(model.TrafficStatsAccessLevel),
				domain.Permissions.TrafficStats,
			)
			assert.Equal(t, model.TotpSecret, domain.TOTP.Secret)
			assert.Equal(t, model.TotpValidated, domain.TOTP.Validated)
		})
	})

	t.Run("toModel", func(t *testing.T) {
		t.Run("converts domain to model with secret", func(t *testing.T) {
			secret := "secret"
			domain := &user.User{
				TOTP: user.TOTP{
					Secret:    &secret,
					Validated: true,
				},
			}

			model := toModel(domain)

			assert.Equal(t, &secret, model.TotpSecret)
			assert.True(t, model.TotpValidated)
		})

		t.Run("successfully converts a complete domain to model", func(t *testing.T) {
			domain := &user.User{
				ID:           uuid.New(),
				Enabled:      true,
				Name:         "Name",
				Username:     "username",
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
				TOTP: user.TOTP{
					Secret:    new("secret"),
					Validated: true,
				},
			}

			model := toModel(domain)

			assert.Equal(t, domain.ID, model.ID)
			assert.Equal(t, domain.Enabled, model.Enabled)
			assert.Equal(t, domain.Name, model.Name)
			assert.Equal(t, domain.Username, model.Username)
			assert.Equal(t, domain.PasswordHash, model.PasswordHash)
			assert.Equal(t, domain.PasswordSalt, model.PasswordSalt)
			assert.Equal(t, string(domain.Permissions.Hosts), model.HostsAccessLevel)
			assert.Equal(t, string(domain.Permissions.Streams), model.StreamsAccessLevel)
			assert.Equal(t, string(domain.Permissions.Certificates), model.CertificatesAccessLevel)
			assert.Equal(t, string(domain.Permissions.Logs), model.LogsAccessLevel)
			assert.Equal(t, string(domain.Permissions.Integrations), model.IntegrationsAccessLevel)
			assert.Equal(t, string(domain.Permissions.AccessLists), model.AccessListsAccessLevel)
			assert.Equal(t, string(domain.Permissions.Settings), model.SettingsAccessLevel)
			assert.Equal(t, string(domain.Permissions.Users), model.UsersAccessLevel)
			assert.Equal(t, string(domain.Permissions.NginxServer), model.NginxServerAccessLevel)
			assert.Equal(t, string(domain.Permissions.ExportData), model.ExportDataAccessLevel)
			assert.Equal(t, string(domain.Permissions.VPNs), model.VPNsAccessLevel)
			assert.Equal(t, string(domain.Permissions.Caches), model.CachesAccessLevel)
			assert.Equal(t, string(domain.Permissions.TrafficStats), model.TrafficStatsAccessLevel)
			assert.Equal(t, domain.TOTP.Secret, model.TotpSecret)
			assert.Equal(t, domain.TOTP.Validated, model.TotpValidated)
		})

		t.Run("converts domain to model with empty secret mapping to nil", func(t *testing.T) {
			domain := &user.User{
				TOTP: user.TOTP{
					Secret:    new(""),
					Validated: false,
				},
			}

			model := toModel(domain)

			assert.Nil(t, model.TotpSecret)
		})

		t.Run("converts domain to model with nil secret mapping to nil", func(t *testing.T) {
			domain := &user.User{
				TOTP: user.TOTP{
					Secret:    nil,
					Validated: false,
				},
			}

			model := toModel(domain)

			assert.Nil(t, model.TotpSecret)
		})

		t.Run("converts domain to model with blank secret mapping to nil", func(t *testing.T) {
			domain := &user.User{
				TOTP: user.TOTP{
					Secret:    new("   "),
					Validated: false,
				},
			}

			model := toModel(domain)

			assert.Nil(t, model.TotpSecret)
		})
	})
}
