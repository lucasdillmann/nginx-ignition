package authorization

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

const testJwtSecret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

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
	commands := user.NewMockedCommands(controller)

	configOverrides := map[string]string{
		"nginx-ignition.security.jwt.secret": testJwtSecret,
	}
	for key, value := range overrides {
		configOverrides[key] = value
	}

	authorizer, err := New(configuration.NewWithOverrides(configOverrides), commands)
	require.NoError(t, err)

	return authorizer, commands
}
