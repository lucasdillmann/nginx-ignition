package user

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/user"
)

func toDomain(dto *userRequestDTO) *user.SaveRequest {
	return &user.SaveRequest{
		ID:       uuid.New(),
		Enabled:  *dto.Enabled,
		Name:     *dto.Name,
		Username: *dto.Username,
		Password: dto.Password,
		Permissions: user.Permissions{
			Hosts:        user.AccessLevel(dto.Permissions.Hosts),
			Streams:      user.AccessLevel(dto.Permissions.Streams),
			Certificates: user.AccessLevel(dto.Permissions.Certificates),
			Logs:         user.AccessLevel(dto.Permissions.Logs),
			Integrations: user.AccessLevel(dto.Permissions.Integrations),
			AccessLists:  user.AccessLevel(dto.Permissions.AccessLists),
			Settings:     user.AccessLevel(dto.Permissions.Settings),
			Users:        user.AccessLevel(dto.Permissions.Users),
			NginxServer:  user.AccessLevel(dto.Permissions.NginxServer),
			ExportData:   user.AccessLevel(dto.Permissions.ExportData),
			VPNs:         user.AccessLevel(dto.Permissions.VPNs),
			Caches:       user.AccessLevel(dto.Permissions.Caches),
		},
	}
}

func toDTO(domain *user.User) *userResponseDTO {
	return &userResponseDTO{
		ID:       domain.ID,
		Enabled:  domain.Enabled,
		Name:     domain.Name,
		Username: domain.Username,
		Permissions: userPermissionsDTO{
			Hosts:        string(domain.Permissions.Hosts),
			Streams:      string(domain.Permissions.Streams),
			Certificates: string(domain.Permissions.Certificates),
			Logs:         string(domain.Permissions.Logs),
			Integrations: string(domain.Permissions.Integrations),
			AccessLists:  string(domain.Permissions.AccessLists),
			Settings:     string(domain.Permissions.Settings),
			Users:        string(domain.Permissions.Users),
			NginxServer:  string(domain.Permissions.NginxServer),
			ExportData:   string(domain.Permissions.ExportData),
			VPNs:         string(domain.Permissions.VPNs),
			Caches:       string(domain.Permissions.Caches),
		},
	}
}
