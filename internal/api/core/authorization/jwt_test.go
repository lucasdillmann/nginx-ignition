package authorization

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func Test_Jwt_GenerateToken(t *testing.T) {
	t.Run("returns a signed JWT with the expected claims", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()

		token, err := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)
		require.NoError(t, err)
		require.NotNil(t, token)
		require.NotEmpty(t, *token)

		claims := parseJwtClaims(t, authorizer, *token)
		assert.Equal(t, uniqueIdentifier, claims["iss"])
		assert.Equal(t, uniqueIdentifier, claims["aud"])
		assert.Equal(t, usr.ID.String(), claims["sub"])
		assert.Equal(t, SessionKind, TokenKind(claims[tokenKindClaim].(float64)))
		assert.NotEmpty(t, claims["jti"])
		assert.NotEmpty(t, claims["exp"])
	})

	t.Run("issues an API token with the provided ID and no expiration", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()

		token, err := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, nil)
		require.NoError(t, err)

		claims := parseJwtClaims(t, authorizer, *token)
		assert.Equal(t, APIKind, TokenKind(claims[tokenKindClaim].(float64)))
		assert.Equal(t, tokenID.String(), claims["jti"])
		assert.Empty(t, claims["exp"])
	})

	t.Run("issues an API token with the provided expiration", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		expiration := time.Now().Add(time.Hour * 24).Truncate(time.Second)

		token, err := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, &expiration)
		require.NoError(t, err)

		claims := parseJwtClaims(t, authorizer, *token)
		assert.Equal(t, tokenID.String(), claims["jti"])
		assert.Equal(t, expiration.Unix(), int64(claims["exp"].(float64)))
	})

	t.Run("ignores the provided expiration for session tokens", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		expiration := time.Now().Add(time.Hour * 24 * 365)

		token, err := authorizer.Jwt().GenerateToken(usr, SessionKind, &tokenID, &expiration)
		require.NoError(t, err)

		claims := parseJwtClaims(t, authorizer, *token)
		assert.Less(t, int64(claims["exp"].(float64)), expiration.Unix())
	})
}

func Test_Jwt_ValidateToken(t *testing.T) {
	t.Run("returns the subject for a valid token", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)

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
		token, _ := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)
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
		token, _ := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)
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
		token, _ := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(nil, errors.New("user not found"))
		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("rejects a session token when the user no longer exists", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)
		tokenID := parseJwtClaims(t, authorizer, *token)["jti"].(string)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(nil, nil)

		assert.NotPanics(t, func() {
			subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

			assert.Nil(t, subject)
			requireUnauthorized(t, err)
		})
		assert.False(t, authorizer.Jwt().isRevoked(tokenID))
	})

	t.Run("rejects an API token when the user no longer exists", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		token, _ := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, nil)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(nil, nil)

		assert.NotPanics(t, func() {
			subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

			assert.Nil(t, subject)
			requireUnauthorized(t, err)
		})
		assert.False(t, authorizer.Jwt().isRevoked(tokenID.String()))
	})

	t.Run("rejects a token without a kind claim", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()

		raw := signRawToken(t, authorizer, jwt.MapClaims{
			"aud": uniqueIdentifier,
			"iss": uniqueIdentifier,
			"sub": usr.ID.String(),
			"jti": uuid.New().String(),
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), raw)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("rejects a token with an unknown kind claim", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()

		raw := signRawToken(t, authorizer, jwt.MapClaims{
			"aud":          uniqueIdentifier,
			"iss":          uniqueIdentifier,
			"sub":          usr.ID.String(),
			"jti":          uuid.New().String(),
			"exp":          time.Now().Add(time.Hour).Unix(),
			tokenKindClaim: 42,
		})

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), raw)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("accepts an existing API token", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		token, _ := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, nil)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
		commands.EXPECT().
			FindAPIToken(gomock.Any(), usr.ID, tokenID).
			Return(&user.APIToken{ID: tokenID, UserID: usr.ID}, nil)

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		require.NoError(t, err)
		require.NotNil(t, subject)
		assert.Equal(t, APIKind, subject.Kind)
		assert.Equal(t, tokenID.String(), subject.TokenID)
	})

	t.Run("rejects an API token that was revoked", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		token, _ := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, nil)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
		commands.EXPECT().FindAPIToken(gomock.Any(), usr.ID, tokenID).Return(nil, nil)

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("rejects an API token whose database expiration already passed", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		expiration := time.Now().Add(time.Hour * 24)
		token, _ := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, &expiration)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
		commands.EXPECT().
			FindAPIToken(gomock.Any(), usr.ID, tokenID).
			Return(&user.APIToken{ID: tokenID, UserID: usr.ID, Expiration: &expiration}, nil)

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.NoError(t, err)
		assert.NotNil(t, subject)
	})

	t.Run("rejects an API token whose database expiration is in the past", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		issuedExpiration := time.Now().Add(time.Hour * 24)
		revokedExpiration := time.Now().Add(-time.Hour)
		token, _ := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, &issuedExpiration)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
		commands.EXPECT().
			FindAPIToken(gomock.Any(), usr.ID, tokenID).
			Return(&user.APIToken{ID: tokenID, UserID: usr.ID, Expiration: &revokedExpiration}, nil)

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.Error(t, err)
		assert.Nil(t, subject)
	})

	t.Run("rejects a disabled user without caching an API token as revoked", func(t *testing.T) {
		authorizer, commands := newAuthorizer(t)
		usr := newUser()
		usr.Enabled = false
		tokenID := uuid.New()
		token, _ := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, nil)

		commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

		subject, err := authorizer.Jwt().ValidateToken(context.Background(), *token)

		assert.Error(t, err)
		assert.Nil(t, subject)
		assert.False(t, authorizer.Jwt().isRevoked(tokenID.String()))
	})
}

func Test_Jwt_RefreshToken(t *testing.T) {
	t.Run("returns no token when the token is far from expiry", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()
		token, _ := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)
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
		token, _ := authorizer.Jwt().GenerateToken(usr, SessionKind, nil, nil)
		subject := subjectFromToken(t, authorizer, *token, usr)
		originalTokenID := subject.TokenID

		refreshed, err := authorizer.Jwt().RefreshToken(subject)
		require.NoError(t, err)
		require.NotNil(t, refreshed)

		claims := parseJwtClaims(t, authorizer, *refreshed)
		assert.Equal(t, usr.ID.String(), claims["sub"])
		assert.NotEqual(t, originalTokenID, claims["jti"])
	})

	t.Run("returns an error for an API token", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)
		usr := newUser()
		tokenID := uuid.New()
		token, _ := authorizer.Jwt().GenerateToken(usr, APIKind, &tokenID, nil)
		subject := subjectFromToken(t, authorizer, *token, usr)
		subject.Kind = APIKind

		refreshed, err := authorizer.Jwt().RefreshToken(subject)

		assert.EqualError(t, err, "token cannot be refreshed")
		assert.Nil(t, refreshed)
	})

	t.Run("returns an error for a nil subject", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)

		refreshed, err := authorizer.Jwt().RefreshToken(nil)

		assert.EqualError(t, err, "token cannot be refreshed")
		assert.Nil(t, refreshed)
	})

	t.Run("returns an error for a subject without claims", func(t *testing.T) {
		authorizer, _ := newAuthorizer(t)

		refreshed, err := authorizer.Jwt().RefreshToken(&Subject{Kind: SessionKind})

		assert.EqualError(t, err, "token cannot be refreshed")
		assert.Nil(t, refreshed)
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
	newAuthorizer := func(cfg *configuration.Configuration) (*ABAC, error) {
		controller := gomock.NewController(t)
		userCommands := user.NewMockedCommands(controller)
		authorizationCommands := authorization.NewMockedCommands(controller)
		authorizationCommands.EXPECT().JwtSecret().Return(testJwtSecret).AnyTimes()

		return New(cfg, userCommands, authorizationCommands)
	}

	t.Run("returns an error when ttl-seconds is less than 30", func(t *testing.T) {
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.ttl-seconds": "10",
		})

		authorizer, err := newAuthorizer(cfg)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
		assert.Contains(t, err.Error(), "ttl-seconds cannot be less than 30")
	})

	t.Run("returns an error when clock-skew-seconds is negative", func(t *testing.T) {
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.clock-skew-seconds": "-1",
		})

		authorizer, err := newAuthorizer(cfg)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
		assert.Contains(t, err.Error(), "clock-skew-seconds cannot be negative")
	})

	t.Run("returns an error when renew-window-seconds is negative", func(t *testing.T) {
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.renew-window-seconds": "-1",
		})

		authorizer, err := newAuthorizer(cfg)

		assert.Error(t, err)
		assert.Nil(t, authorizer)
		assert.Contains(t, err.Error(), "renew-window-seconds cannot be negative")
	})

	t.Run("returns an error when renew-window-seconds > ttl-seconds", func(t *testing.T) {
		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.security.jwt.ttl-seconds":          "60",
			"nginx-ignition.security.jwt.renew-window-seconds": "120",
		})

		authorizer, err := newAuthorizer(cfg)

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

		return authorizer.jwt.secretKey(), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	return claims
}

func signRawToken(t *testing.T, authorizer *ABAC, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	raw, err := token.SignedString(authorizer.jwt.secretKey())
	require.NoError(t, err)

	return raw
}

func subjectFromToken(t *testing.T, authorizer *ABAC, raw string, usr *user.User) *Subject {
	t.Helper()
	claims := parseJwtClaims(t, authorizer, raw)

	return &Subject{
		Kind:    SessionKind,
		User:    usr,
		TokenID: claims["jti"].(string),
		claims:  &claims,
	}
}
