package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
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

func Test_toAPITokenDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		token := newAPIToken(newUser())

		result := toAPITokenDTO(token)

		assert.NotNil(t, result)
		assert.Equal(t, token.ID, result.ID)
		assert.Equal(t, token.Name, result.Name)
		assert.Equal(t, token.Expiration, result.Expiration)
		assert.Equal(t, token.CreatedAt, result.CreatedAt)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		assert.Nil(t, toAPITokenDTO(nil))
	})
}

func Test_toAPITokenCreatedDTO(t *testing.T) {
	t.Run("converts domain object to DTO including the token", func(t *testing.T) {
		token := newAPIToken(newUser())

		result := toAPITokenCreatedDTO(token, "jwt-token")

		assert.NotNil(t, result)
		assert.Equal(t, token.ID, result.ID)
		assert.Equal(t, token.Name, result.Name)
		assert.Equal(t, token.Expiration, result.Expiration)
		assert.Equal(t, token.CreatedAt, result.CreatedAt)
		assert.Equal(t, "jwt-token", result.Token)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		assert.Nil(t, toAPITokenCreatedDTO(nil, "jwt-token"))
	})
}

func Test_toAPITokenDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		expiration := time.Now().UTC().Truncate(time.Second)
		payload := apiTokenCreateRequestDTO{
			Name:       new("automation"),
			Expiration: &expiration,
		}

		result := toAPITokenDomain(&payload)

		assert.NotNil(t, result)
		assert.Equal(t, "automation", result.Name)
		assert.Equal(t, &expiration, result.Expiration)
	})

	t.Run("converts a DTO without expiration", func(t *testing.T) {
		result := toAPITokenDomain(&apiTokenCreateRequestDTO{Name: new("automation")})

		assert.NotNil(t, result)
		assert.Nil(t, result.Expiration)
	})

	t.Run("converts a DTO without name", func(t *testing.T) {
		result := toAPITokenDomain(&apiTokenCreateRequestDTO{})

		assert.NotNil(t, result)
		assert.Empty(t, result.Name)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		assert.Nil(t, toAPITokenDomain(nil))
	})
}
