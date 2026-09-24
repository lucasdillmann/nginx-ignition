package authorization

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	user2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

const testJwtSecret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func newPermissions() user2.Permissions {
	return user2.Permissions{
		Hosts:        user2.NoAccessAccessLevel,
		Streams:      user2.NoAccessAccessLevel,
		Certificates: user2.NoAccessAccessLevel,
		Logs:         user2.NoAccessAccessLevel,
		Integrations: user2.NoAccessAccessLevel,
		AccessLists:  user2.NoAccessAccessLevel,
		Settings:     user2.NoAccessAccessLevel,
		Users:        user2.NoAccessAccessLevel,
		NginxServer:  user2.NoAccessAccessLevel,
		ExportData:   user2.NoAccessAccessLevel,
		VPNs:         user2.NoAccessAccessLevel,
		Caches:       user2.NoAccessAccessLevel,
		TrafficStats: user2.NoAccessAccessLevel,
	}
}

func newUser() *user2.User {
	return &user2.User{
		ID:          uuid.New(),
		Username:    "testuser",
		Name:        "Test User",
		Enabled:     true,
		Permissions: newPermissions(),
	}
}

func newAuthorizer(t *testing.T) (*ABAC, *user2.MockedCommands) {
	t.Helper()
	return newAuthorizerWithOverrides(t, map[string]string{})
}

func newAuthorizerWithOverrides(
	t *testing.T,
	overrides map[string]string,
) (*ABAC, *user2.MockedCommands) {
	t.Helper()
	controller := gomock.NewController(t)
	commands := user2.NewMockedCommands(controller)

	configOverrides := map[string]string{
		"nginx-ignition.security.jwt.secret":               testJwtSecret,
		"nginx-ignition.security.jwt.clock-skew-seconds":   "0",
		"nginx-ignition.security.jwt.ttl-seconds":          "30",
		"nginx-ignition.security.jwt.renew-window-seconds": "10",
	}
	for key, value := range overrides {
		configOverrides[key] = value
	}

	authorizer, err := New(configuration.NewWithOverrides(configOverrides), commands)
	require.NoError(t, err)

	return authorizer, commands
}
