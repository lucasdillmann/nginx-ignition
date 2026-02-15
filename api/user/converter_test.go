package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_toDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		subject := newUser()
		subject.TOTP = user.TOTP{
			Secret:    new("secret"),
			Validated: true,
		}
		result := toDTO(subject)

		assert.NotNil(t, result)
		assert.Equal(t, subject.ID, result.ID)
		assert.Equal(t, subject.Name, result.Name)
		assert.Equal(t, subject.Username, result.Username)
		assert.True(t, result.Enabled)
		assert.True(t, result.TOTPEnabled)
		assert.Equal(t, string(user.ReadWriteAccessLevel), result.Permissions.Hosts)
	})

	t.Run(
		"converts domain object to DTO with TOTP disabled when not validated",
		func(t *testing.T) {
			subject := newUser()
			subject.TOTP = user.TOTP{
				Secret:    new("secret"),
				Validated: false,
			}
			result := toDTO(subject)

			assert.NotNil(t, result)
			assert.False(t, result.TOTPEnabled)
		},
	)

	t.Run(
		"converts domain object to DTO with TOTP disabled when secret is nil",
		func(t *testing.T) {
			subject := newUser()
			subject.TOTP = user.TOTP{
				Secret:    nil,
				Validated: true,
			}
			result := toDTO(subject)

			assert.NotNil(t, result)
			assert.False(t, result.TOTPEnabled)
		},
	)

	t.Run(
		"converts domain object to DTO with TOTP disabled when secret is empty",
		func(t *testing.T) {
			subject := newUser()
			subject.TOTP = user.TOTP{
				Secret:    new("  "),
				Validated: true,
			}
			result := toDTO(subject)

			assert.NotNil(t, result)
			assert.False(t, result.TOTPEnabled)
		},
	)

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil)
		assert.Nil(t, result)
	})
}

func Test_toDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		payload := newUserRequest()
		result := toDomain(&payload)

		assert.NotNil(t, result)
		assert.Equal(t, *payload.Name, result.Name)
		assert.Equal(t, *payload.Username, result.Username)
		assert.True(t, result.Enabled)
		assert.False(t, result.RemoveTOTP)
		assert.Equal(t, user.ReadWriteAccessLevel, result.Permissions.Hosts)
	})
}
