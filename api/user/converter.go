package user

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/user"
)

func toDomain(dto *userRequestDto) *user.SaveRequest {
	return &user.SaveRequest{
		ID:       uuid.New(),
		Enabled:  *dto.Enabled,
		Name:     *dto.Name,
		Username: *dto.Username,
		Password: dto.Password,
		Permissions: user.Permissions{
			Hosts:        user.AccessLevel(*dto.Permissions.Hosts),
			Streams:      user.AccessLevel(*dto.Permissions.Streams),
			Certificates: user.AccessLevel(*dto.Permissions.Certificates),
			Logs:         user.AccessLevel(*dto.Permissions.Logs),
			Integrations: user.AccessLevel(*dto.Permissions.Integrations),
			AccessLists:  user.AccessLevel(*dto.Permissions.AccessLists),
			Settings:     user.AccessLevel(*dto.Permissions.Settings),
			Users:        user.AccessLevel(*dto.Permissions.Users),
			NginxServer:  user.AccessLevel(*dto.Permissions.NginxServer),
			ExportData:   user.AccessLevel(*dto.Permissions.ExportData),
			VPNs:         user.AccessLevel(*dto.Permissions.VPNs),
		},
	}
}

func toDto(domain *user.User) *userResponseDto {
	return &userResponseDto{
		ID:       domain.ID,
		Enabled:  domain.Enabled,
		Name:     domain.Name,
		Username: domain.Username,
		Permissions: userPermissionsDto{
			Hosts:        ptr.Of(string(domain.Permissions.Hosts)),
			Streams:      ptr.Of(string(domain.Permissions.Streams)),
			Certificates: ptr.Of(string(domain.Permissions.Certificates)),
			Logs:         ptr.Of(string(domain.Permissions.Logs)),
			Integrations: ptr.Of(string(domain.Permissions.Integrations)),
			AccessLists:  ptr.Of(string(domain.Permissions.AccessLists)),
			Settings:     ptr.Of(string(domain.Permissions.Settings)),
			Users:        ptr.Of(string(domain.Permissions.Users)),
			NginxServer:  ptr.Of(string(domain.Permissions.NginxServer)),
			ExportData:   ptr.Of(string(domain.Permissions.ExportData)),
			VPNs:         ptr.Of(string(domain.Permissions.VPNs)),
		},
	}
}
