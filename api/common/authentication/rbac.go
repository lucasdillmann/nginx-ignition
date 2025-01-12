package authentication

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/user"
)

type RBAC struct {
	configuration     *configuration.Configuration
	anonymousPaths    []string
	roleRequiredPaths map[string]user.Role
	jwt               *Jwt
}

func New(configuration *configuration.Configuration, repository *user.Repository) (*RBAC, error) {
	jwt, err := newJwt(configuration, repository)
	if err != nil {
		return nil, err
	}

	return &RBAC{
		configuration:     configuration,
		anonymousPaths:    []string{},
		roleRequiredPaths: map[string]user.Role{},
		jwt:               jwt,
	}, nil
}

func (m *RBAC) Jwt() *Jwt {
	return m.jwt
}

func (m *RBAC) AllowAnonymous(path string) {
	m.anonymousPaths = append(m.anonymousPaths, path)
}

func (m *RBAC) RequireRole(method, path string, role user.Role) {
	m.roleRequiredPaths[method+":"+path] = role
}

func (m *RBAC) isAnonymous(path string) bool {
	for _, p := range m.anonymousPaths {
		if p == path {
			return true
		}
	}

	return false
}

func (m *RBAC) findRequiredRole(method, path string) *user.Role {
	role, exists := m.roleRequiredPaths[method+":"+path]
	if !exists {
		return nil
	}

	return &role
}
