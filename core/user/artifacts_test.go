package user

import (
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

func newTestCommands(
	ctrl *gomock.Controller,
	repository Repository,
	cfg *configuration.Configuration,
) (*service, Commands) {
	i18nCommands := i18n.NewMockedCommands(ctrl)
	i18nCommands.EXPECT().Supports(gomock.Any()).Return(true).AnyTimes()

	return newCommands(repository, cfg, i18nCommands)
}

func newTestValidator(ctrl *gomock.Controller, repository Repository) *validator {
	i18nCommands := i18n.NewMockedCommands(ctrl)
	i18nCommands.EXPECT().Supports(gomock.Any()).Return(true).AnyTimes()

	return newValidator(repository, i18nCommands)
}

func newUser() *User {
	return &User{
		ID:                   uuid.New(),
		Username:             "testuser",
		Name:                 "Test User",
		NotificationLanguage: "en",
		Enabled:              true,
		Permissions: Permissions{
			Hosts:        NoAccessAccessLevel,
			Streams:      NoAccessAccessLevel,
			Certificates: NoAccessAccessLevel,
			Logs:         NoAccessAccessLevel,
			Integrations: NoAccessAccessLevel,
			AccessLists:  NoAccessAccessLevel,
			Settings:     NoAccessAccessLevel,
			Users:        NoAccessAccessLevel,
			NginxServer:  ReadOnlyAccessLevel,
			ExportData:   NoAccessAccessLevel,
			VPNs:         NoAccessAccessLevel,
			Caches:       NoAccessAccessLevel,
			TrafficStats: NoAccessAccessLevel,
		},
	}
}

func newSaveRequest() *SaveRequest {
	return &SaveRequest{
		ID:                   uuid.New(),
		Username:             "testuser",
		Name:                 "Test User",
		NotificationLanguage: "en",
		Enabled:              true,
		Password:             new("password123"),
		RemoveTOTP:           false,
		Permissions: Permissions{
			Hosts:        NoAccessAccessLevel,
			Streams:      NoAccessAccessLevel,
			Certificates: NoAccessAccessLevel,
			Logs:         NoAccessAccessLevel,
			Integrations: NoAccessAccessLevel,
			AccessLists:  NoAccessAccessLevel,
			Settings:     NoAccessAccessLevel,
			Users:        NoAccessAccessLevel,
			NginxServer:  ReadOnlyAccessLevel,
			ExportData:   NoAccessAccessLevel,
			VPNs:         NoAccessAccessLevel,
			Caches:       NoAccessAccessLevel,
			TrafficStats: NoAccessAccessLevel,
		},
	}
}
