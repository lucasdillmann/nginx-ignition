package authorization

import (
	"cmp"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	user2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type PermissionResolver func(permissions user2.Permissions) user2.AccessLevel

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

func New(cfg *configuration.Configuration, commands user2.Commands) (*ABAC, error) {
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

func (m *ABAC) isAccessGranted(method, path string, permissions *user2.Permissions) bool {
	currentAccessLevel := user2.NoAccessAccessLevel
	for _, item := range m.permissionResolvers {
		if strings.HasPrefix(path, item.path) {
			currentAccessLevel = item.resolver(*permissions)
			break
		}
	}

	switch method {
	case "GET":
		return currentAccessLevel == user2.ReadOnlyAccessLevel ||
			currentAccessLevel == user2.ReadWriteAccessLevel
	case "POST", "PUT", "DELETE", "PATCH":
		return currentAccessLevel == user2.ReadWriteAccessLevel
	default:
		return false
	}
}
