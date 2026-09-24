package authorization

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func newABAC() *ABAC {
	return &ABAC{
		permissionResolvers: make([]pathPermissionResolver, 0),
	}
}

func Test_ConfigureGroup(t *testing.T) {
	t.Run("appends a resolver and returns a router group", func(t *testing.T) {
		m := newABAC()
		group := m.ConfigureGroup(gin.New(), "/api/nginx/logs", nil)

		require.NotNil(t, group)
		require.Len(t, m.permissionResolvers, 1)
		assert.Equal(t, "/api/nginx/logs", m.permissionResolvers[0].path)
	})

	t.Run("sorts resolvers by path length in descending order", func(t *testing.T) {
		m := newABAC()
		m.ConfigureGroup(gin.New(), "/api/nginx", nil)
		m.ConfigureGroup(gin.New(), "/api", nil)
		m.ConfigureGroup(gin.New(), "/api/nginx/logs", nil)

		assert.Equal(t, []string{"/api/nginx/logs", "/api/nginx", "/api"}, resolverPaths(m))
	})

	t.Run("does not depend on registration order", func(t *testing.T) {
		m := newABAC()
		m.ConfigureGroup(gin.New(), "/api/nginx/logs", nil)
		m.ConfigureGroup(gin.New(), "/api/nginx", nil)
		m.ConfigureGroup(gin.New(), "/api", nil)

		assert.Equal(t, []string{"/api/nginx/logs", "/api/nginx", "/api"}, resolverPaths(m))
	})
}

func Test_isAnonymous(t *testing.T) {
	m := newABAC()
	m.AllowAnonymous("GET", "/api/anonymous")

	assert.True(t, m.isAnonymous("GET", "/api/anonymous"))
	assert.False(t, m.isAnonymous("POST", "/api/anonymous"))
	assert.False(t, m.isAnonymous("GET", "/api/other"))
}

func Test_isAllowedForAllUsers(t *testing.T) {
	m := newABAC()
	m.AllowAllUsers("GET", "/api/all-users")

	assert.True(t, m.isAllowedForAllUsers("GET", "/api/all-users"))
	assert.False(t, m.isAllowedForAllUsers("POST", "/api/all-users"))
	assert.False(t, m.isAllowedForAllUsers("GET", "/api/other"))
}

func Test_isAccessGranted(t *testing.T) {
	m := newABAC()
	m.ConfigureGroup(
		gin.New(),
		"/api",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Settings },
	)
	m.ConfigureGroup(
		gin.New(),
		"/api/nginx",
		func(permissions user.Permissions) user.AccessLevel { return permissions.NginxServer },
	)
	m.ConfigureGroup(
		gin.New(),
		"/api/nginx/logs",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Logs },
	)
	m.ConfigureGroup(
		gin.New(),
		"/api/hosts",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Hosts },
	)
	m.ConfigureGroup(
		gin.New(),
		"/api/hosts/:id/logs",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Logs },
	)

	permissions := newPermissions()
	permissions.NginxServer = user.ReadOnlyAccessLevel
	permissions.Logs = user.ReadWriteAccessLevel

	t.Run("resolves the most specific matching path prefix", func(t *testing.T) {
		cases := []struct {
			method   string
			path     string
			expected bool
		}{
			{method: "GET", path: "/api/nginx/status", expected: true},
			{method: "POST", path: "/api/nginx/status", expected: false},
			{method: "GET", path: "/api/nginx/logs", expected: true},
			{method: "POST", path: "/api/nginx/logs", expected: true},
			{method: "GET", path: "/api/hosts/:id", expected: false},
			{method: "GET", path: "/api/hosts/:id/logs", expected: true},
			{method: "GET", path: "/api/settings", expected: false},
			{method: "GET", path: "/api", expected: false},
		}

		for _, tc := range cases {
			t.Run(tc.method+" "+tc.path, func(t *testing.T) {
				assert.Equal(t, tc.expected, m.isAccessGranted(tc.method, tc.path, &permissions))
			})
		}
	})

	t.Run(
		"denies access to unregistered subpaths through the closest resolver",
		func(t *testing.T) {
			t.Run("GET /api/nginx/config uses the nginx resolver", func(t *testing.T) {
				assert.True(t, m.isAccessGranted("GET", "/api/nginx/config", &permissions))
			})
		},
	)

	t.Run("rejects unknown HTTP methods", func(t *testing.T) {
		assert.False(t, m.isAccessGranted("HEAD", "/api/nginx/status", &permissions))
	})

	t.Run("denies when a more specific resolver grants no access", func(t *testing.T) {
		restricted := newPermissions()
		restricted.NginxServer = user.ReadOnlyAccessLevel

		assert.False(t, m.isAccessGranted("POST", "/api/nginx/logs", &restricted))
	})

	t.Run("denies read-write methods with read-only access", func(t *testing.T) {
		readOnly := newPermissions()
		readOnly.NginxServer = user.ReadOnlyAccessLevel
		readOnly.Hosts = user.ReadOnlyAccessLevel

		for _, method := range []string{"POST", "PUT", "DELETE", "PATCH"} {
			t.Run(method, func(t *testing.T) {
				assert.False(t, m.isAccessGranted(method, "/api/nginx/status", &readOnly))
				assert.False(t, m.isAccessGranted(method, "/api/hosts/:id", &readOnly))
			})
		}
	})
}

func resolverPaths(m *ABAC) []string {
	paths := make([]string, 0, len(m.permissionResolvers))
	for _, item := range m.permissionResolvers {
		paths = append(paths, item.path)
	}

	return paths
}
