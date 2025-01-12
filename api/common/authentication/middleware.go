package authentication

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/user"
)

type Middleware struct {
	configuration     *configuration.Configuration
	anonymousPaths    []string
	roleRequiredPaths map[string]user.Role
	jwt               *Jwt
}

func New(configuration *configuration.Configuration, repository *user.Repository) (*Middleware, error) {
	jwt, err := newJwt(configuration, repository)
	if err != nil {
		return nil, err
	}

	return &Middleware{
		configuration:     configuration,
		anonymousPaths:    []string{},
		roleRequiredPaths: map[string]user.Role{},
		jwt:               jwt,
	}, nil
}

func (m *Middleware) Jwt() *Jwt {
	return m.jwt
}

func (m *Middleware) AllowAnonymous(path string) {
	m.anonymousPaths = append(m.anonymousPaths, path)
}

func (m *Middleware) RequireRole(path string, role user.Role) {
	m.roleRequiredPaths[path] = role
}

func (m *Middleware) isAnonymous(path string) bool {
	for _, p := range m.anonymousPaths {
		if p == path {
			return true
		}
	}

	return false
}

func (m *Middleware) findRequiredRole(path string) *user.Role {
	role, exists := m.roleRequiredPaths[path]
	if !exists {
		return nil
	}

	return &role
}
