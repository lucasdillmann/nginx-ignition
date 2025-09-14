package authorization

import (
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

type PermissionResolver func(permissions user.Permissions) user.AccessLevel

type ABAC struct {
	configuration           *configuration.Configuration
	anonymousPaths          []string
	allowedForAllUsers      []string
	pathPermissionResolvers map[string]PermissionResolver
	jwt                     *Jwt
}

func New(configuration *configuration.Configuration, repository user.Repository) (*ABAC, error) {
	jwt, err := newJwt(configuration, repository)
	if err != nil {
		return nil, err
	}

	return &ABAC{
		configuration:           configuration,
		anonymousPaths:          []string{},
		pathPermissionResolvers: map[string]PermissionResolver{},
		jwt:                     jwt,
	}, nil
}

func (m *ABAC) Jwt() *Jwt {
	return m.jwt
}

func (m *ABAC) AllowAnonymous(method, path string) {
	m.anonymousPaths = append(m.anonymousPaths, method+":"+path)
}

func (m *ABAC) AllowAllUsers(method, path string) {
	m.allowedForAllUsers = append(m.allowedForAllUsers, method+":"+path)
}

func (m *ABAC) ConfigureGroup(
	router *gin.Engine,
	path string,
	permissionResolver PermissionResolver,
) *gin.RouterGroup {
	m.pathPermissionResolvers[path] = permissionResolver
	return router.Group(path)
}

func (m *ABAC) isAnonymous(method, path string) bool {
	for _, p := range m.anonymousPaths {
		if p == method+":"+path {
			return true
		}
	}

	return false
}

func (m *ABAC) isAllowedForAllUsers(method, path string) bool {
	for _, p := range m.allowedForAllUsers {
		if p == method+":"+path {
			return true
		}
	}

	return false
}

func (m *ABAC) isAccessGranted(method, path string, permissions *user.Permissions) bool {
	currentAccessLevel := user.NoAccessAccessLevel
	for basePath, resolver := range m.pathPermissionResolvers {
		if strings.HasPrefix(path, basePath) {
			currentAccessLevel = resolver(*permissions)
			break
		}
	}

	switch method {
	case "GET":
		return currentAccessLevel == user.ReadOnlyAccessLevel || currentAccessLevel == user.ReadWriteAccessLevel
	case "POST", "PUT", "DELETE", "PATCH":
		return currentAccessLevel == user.ReadWriteAccessLevel
	default:
		return false
	}
}
