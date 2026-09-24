package authorization

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	user2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func Test_Jwt_GenerateToken(t *testing.T) {
	t.Run("returns a signed JWT with the expected claims", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()

		token, err := authorizer.Jwt().GenerateToken(usr)
		require.NoError(t, err)
		require.NotNil(t, token)
		require.NotEmpty(t, *token)

		claims := parseJwtClaims(t, authorizer, *token)
		assert.Equal(t, uniqueIdentifier, claims["iss"])
		assert.Equal(t, uniqueIdentifier, claims["aud"])
		assert.Equal(t, usr.ID.String(), claims["sub"])
		assert.NotEmpty(t, claims["jti"])
	})
}

func Test_Jwt_ValidateToken(t *testing.T) {
	t.Run("returns the subject for a valid token", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		require.NoError(t, err)
		require.NotNil(t, subject)
		assert.Equal(t, usr.ID, subject.User.ID)
		assert.NotEmpty(t, subject.TokenID)
	})

	t.Run("rejects a malformed token", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), "not-a-jwt")

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("rejects a token signed with a different key", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"aud": uniqueIdentifier,
			"iss": uniqueIdentifier,
			"sub": usr.ID.String(),
			"jti": uuid.New().String(),
		})
		raw, err := token.SignedString([]byte("different-key"))
		require.NoError(t, err)

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), raw)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("rejects a revoked token", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr)
		tokenID := parseJwtClaims(t, authorizer, *token)["jti"].(string)
		authorizer.Jwt().RevokeToken(tokenID)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("rejects a disabled user and revokes the token", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		usr.Enabled = false
		token, _ := authorizer.Jwt().GenerateToken(usr)
		tokenID := parseJwtClaims(t, authorizer, *token)["jti"].(string)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.Error(t, err)
		assert.Nil(t, subject)
		assert.True(t, authorizer.Jwt().isRevoked(tokenID))
	})

	t.Run("rejects a token for an unknown user", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(nil, errors.New("user not found"))
		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})
}

func Test_Jwt_RefreshToken(t *testing.T) {
	t.Run("returns no token when the token is far from expiry", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr)
		subject := subjectFromToken(t, authorizer, *token, usr)

		refreshed, err := authorizer.Jwt().RefreshToken(subject)

		require.NoError(t, err)
		assert.Nil(t, refreshed)
	})

	t.Run("returns a new token when within the renewal window", func(t *testing.T) {
		authorizer, _ := newAuthorizerWithOverrides(t, map[string]string{
			"nginx-ignition.security.jwt.renew-window-seconds": "30",
		})
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr)
		subject := subjectFromToken(t, authorizer, *token, usr)
		originalTokenID := subject.TokenID

		refreshed, err := authorizer.Jwt().RefreshToken(subject)
		require.NoError(t, err)
		require.NotNil(t, refreshed)

		claims := parseJwtClaims(t, authorizer, *refreshed)
		assert.Equal(t, usr.ID.String(), claims["sub"])
		assert.NotEqual(t, originalTokenID, claims["jti"])
	})
}

func Test_Jwt_RevokeToken(t *testing.T) {
	t.Run("marks a token as revoked", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		jwtInstance := authorizer.Jwt()

		assert.False(t, jwtInstance.isRevoked("token-id"))
		jwtInstance.RevokeToken("token-id")
		assert.True(t, jwtInstance.isRevoked("token-id"))
		assert.False(t, jwtInstance.isRevoked("other-token-id"))
	})
}

func Test_New(t *testing.T) {
	t.Run("returns an error for an invalid JWT secret length", func(t *testing.T) {
		controller := gomock.NewController(t)
		commands := user2.NewMockedCommands(controller)
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.secret": "too-short",
		})

		authorizer, err := New(cfg, commands)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
	})

	t.Run("returns an error when ttl-seconds is less than 30", func(t *testing.T) {
		controller := gomock.NewController(t)
		commands := user2.NewMockedCommands(controller)
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.secret":      testJwtSecret,
			"nginx-ignition.security.jwt.ttl-seconds": "10",
		})

		authorizer, err := New(cfg, commands)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
		assert.Contains(t, err.Error(), "ttl-seconds cannot be less than 30")
	})

	t.Run("returns an error when clock-skew-seconds is negative", func(t *testing.T) {
		controller := gomock.NewController(t)
		commands := user2.NewMockedCommands(controller)
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.secret":             testJwtSecret,
			"nginx-ignition.security.jwt.clock-skew-seconds": "-1",
		})

		authorizer, err := New(cfg, commands)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
		assert.Contains(t, err.Error(), "clock-skew-seconds cannot be negative")
	})

	t.Run("returns an error when renew-window-seconds is negative", func(t *testing.T) {
		controller := gomock.NewController(t)
		commands := user2.NewMockedCommands(controller)
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.secret":               testJwtSecret,
			"nginx-ignition.security.jwt.renew-window-seconds": "-1",
		})

		authorizer, err := New(cfg, commands)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
		assert.Contains(t, err.Error(), "renew-window-seconds cannot be negative")
	})

	t.Run("returns an error when renew-window-seconds > ttl-seconds", func(t *testing.T) {
		controller := gomock.NewController(t)
		commands := user2.NewMockedCommands(controller)
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.secret":               testJwtSecret,
			"nginx-ignition.security.jwt.ttl-seconds":          "60",
			"nginx-ignition.security.jwt.renew-window-seconds": "120",
		})

		authorizer, err := New(cfg, commands)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
		assert.Contains(t, err.Error(), "renew-window-seconds cannot be bigger than ttl-seconds")
	})
}

func parseJwtClaims(t *testing.T, authorizer *ABAC, raw string) jwt.MapClaims {
	t.Helper()
	token, err := jwt.Parse(raw, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return authorizer.jwt.secretKey, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	return claims
}

func subjectFromToken(t *testing.T, authorizer *ABAC, raw string, usr *user2.User) *Subject {
	t.Helper()
	claims := parseJwtClaims(t, authorizer, raw)

	return &Subject{
		User:    usr,
		TokenID: claims["jti"].(string),
		claims:  &claims,
	}
}
