package user

import (
	"dillmann.com.br/nginx-ignition/core/user"
)

func toDomain(model *userModel) user.User {
	return user.User{
		ID:           model.ID,
		Enabled:      model.Enabled,
		Name:         model.Name,
		Username:     model.Username,
		PasswordHash: model.PasswordHash,
		PasswordSalt: model.PasswordSalt,
		Permissions: user.Permissions{
			Hosts:        user.AccessLevel(model.HostsAccessLevel),
			Streams:      user.AccessLevel(model.StreamsAccessLevel),
			Certificates: user.AccessLevel(model.CertificatesAccessLevel),
			Logs:         user.AccessLevel(model.LogsAccessLevel),
			Integrations: user.AccessLevel(model.IntegrationsAccessLevel),
			AccessLists:  user.AccessLevel(model.AccessListsAccessLevel),
			Settings:     user.AccessLevel(model.SettingsAccessLevel),
			Users:        user.AccessLevel(model.UsersAccessLevel),
			NginxServer:  user.AccessLevel(model.NginxServerAccessLevel),
			ExportData:   user.AccessLevel(model.ExportDataAccessLevel),
			VPNs:         user.AccessLevel(model.VPNsAccessLevel),
			Caches:       user.AccessLevel(model.CachesAccessLevel),
			TrafficStats: user.AccessLevel(model.TrafficStatsAccessLevel),
		},
	}
}

func toModel(domain *user.User) userModel {
	return userModel{
		ID:                      domain.ID,
		Enabled:                 domain.Enabled,
		Name:                    domain.Name,
		Username:                domain.Username,
		PasswordHash:            domain.PasswordHash,
		PasswordSalt:            domain.PasswordSalt,
		HostsAccessLevel:        string(domain.Permissions.Hosts),
		StreamsAccessLevel:      string(domain.Permissions.Streams),
		CertificatesAccessLevel: string(domain.Permissions.Certificates),
		LogsAccessLevel:         string(domain.Permissions.Logs),
		IntegrationsAccessLevel: string(domain.Permissions.Integrations),
		AccessListsAccessLevel:  string(domain.Permissions.AccessLists),
		SettingsAccessLevel:     string(domain.Permissions.Settings),
		UsersAccessLevel:        string(domain.Permissions.Users),
		NginxServerAccessLevel:  string(domain.Permissions.NginxServer),
		ExportDataAccessLevel:   string(domain.Permissions.ExportData),
		VPNsAccessLevel:         string(domain.Permissions.VPNs),
		CachesAccessLevel:       string(domain.Permissions.Caches),
		TrafficStatsAccessLevel: string(domain.Permissions.TrafficStats),
	}
}
