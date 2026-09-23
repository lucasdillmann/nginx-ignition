package authorization

import (
	"cmp"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

type PermissionResolver func(permissions user.Permissions) user.AccessLevel

type pathPermissionResolver struct {
	resolver PermissionResolver
	path     string
}

type ABAC struct {
	configuration       *configuration.Configuration
	permissionResolvers []pathPermissionResolver
	jwt                 *Jwt
	anonymousPaths      []string
	allowedForAllUsers  []string
}

func New(cfg *configuration.Configuration, commands user.Commands) (*ABAC, error) {
	jwt, err := newJwt(cfg, commands)
	if err != nil {
		return nil, err
	}

	return &ABAC{
		configuration:       cfg,
		anonymousPaths:      []string{},
		permissionResolvers: make([]pathPermissionResolver, 0),
		jwt:                 jwt,
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
	m.permissionResolvers = append(
		m.permissionResolvers,
		pathPermissionResolver{
			path:     path,
			resolver: permissionResolver,
		},
	)

	slices.SortFunc(m.permissionResolvers, func(left, right pathPermissionResolver) int {
		return cmp.Compare(len(right.path), len(left.path))
	})

	return router.Group(path)
}

func (m *ABAC) isAnonymous(method, path string) bool {
	return slices.Contains(m.anonymousPaths, method+":"+path)
}

func (m *ABAC) isAllowedForAllUsers(method, path string) bool {
	return slices.Contains(m.allowedForAllUsers, method+":"+path)
}

func (m *ABAC) isAccessGranted(method, path string, permissions *user.Permissions) bool {
	currentAccessLevel := user.NoAccessAccessLevel
	for _, item := range m.permissionResolvers {
		if strings.HasPrefix(path, item.path) {
			currentAccessLevel = item.resolver(*permissions)
			break
		}
	}

	switch method {
	case "GET":
		return currentAccessLevel == user.ReadOnlyAccessLevel ||
			currentAccessLevel == user.ReadWriteAccessLevel
	case "POST", "PUT", "DELETE", "PATCH":
		return currentAccessLevel == user.ReadWriteAccessLevel
	default:
		return false
	}
}
