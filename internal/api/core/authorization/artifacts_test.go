package authorization

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

const testJwtSecret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func requireUnauthorized(t *testing.T, err error) {
	t.Helper()

	var apiErr *apierror.APIError

	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
}

func newPermissions() user.Permissions {
	return user.Permissions{
		Hosts:        user.NoAccessAccessLevel,
		Streams:      user.NoAccessAccessLevel,
		Certificates: user.NoAccessAccessLevel,
		Logs:         user.NoAccessAccessLevel,
		Integrations: user.NoAccessAccessLevel,
		AccessLists:  user.NoAccessAccessLevel,
		Settings:     user.NoAccessAccessLevel,
		Users:        user.NoAccessAccessLevel,
		NginxServer:  user.NoAccessAccessLevel,
		ExportData:   user.NoAccessAccessLevel,
		VPNs:         user.NoAccessAccessLevel,
		Caches:       user.NoAccessAccessLevel,
		TrafficStats: user.NoAccessAccessLevel,
	}
}

func newUser() *user.User {
	return &user.User{
		ID:          uuid.New(),
		Username:    "testuser",
		Name:        "Test User",
		Enabled:     true,
		Permissions: newPermissions(),
	}
}

func newSession(usr *user.User) *user.Session {
	createdAt := time.Now().UTC().Truncate(time.Second)

	return &user.Session{
		ID:          uuid.New(),
		IterationID: uuid.New(),
		UserID:      usr.ID,
		CreatedAt:   createdAt,
		ExpiresAt:   createdAt.Add(time.Hour),
	}
}

func issueSessionToken(
	t *testing.T,
	authorizer *ABAC,
	commands *user.MockedCommands,
	usr *user.User,
	session *user.Session,
) *string {
	t.Helper()

	commands.EXPECT().
		StartSession(gomock.Any(), usr.ID, gomock.Any()).
		Return(session, nil)

	token, err := authorizer.Jwt().GenerateToken(
		context.Background(),
		usr,
		SessionKind,
		nil,
		nil,
	)
	require.NoError(t, err)

	return token
}

func newAuthorizer(t *testing.T) (*ABAC, *user.MockedCommands) {
	t.Helper()
	return newAuthorizerWithOverrides(t, map[string]string{})
}

func newAuthorizerWithOverrides(
	t *testing.T,
	overrides map[string]string,
) (*ABAC, *user.MockedCommands) {
	t.Helper()
	controller := gomock.NewController(t)
	userCommands := user.NewMockedCommands(controller)

	authorizationCommands := authorization.NewMockedCommands(controller)
	authorizationCommands.EXPECT().JwtSecret().Return(testJwtSecret).AnyTimes()

	configOverrides := map[string]string{
		"nginx-ignition.security.jwt.clock-skew-seconds":   "0",
		"nginx-ignition.security.jwt.ttl-seconds":          "30",
		"nginx-ignition.security.jwt.renew-window-seconds": "10",
	}
	for key, value := range overrides {
		configOverrides[key] = value
	}

	authorizer, err := New(
		configuration.NewWithOverrides(configOverrides),
		userCommands,
		authorizationCommands,
	)
	require.NoError(t, err)

	return authorizer, userCommands
}
