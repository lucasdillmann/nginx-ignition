package user

import (
	"github.com/google/uuid"
)

type userLoginRequestDto struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type userLoginResponseDto struct {
	Token string `json:"token" validate:"required"`
}

type userOnboardingStatusResponseDto struct {
	Finished bool `json:"finished" validate:"required"`
}

type userPasswordUpdateRequestDto struct {
	CurrentPassword *string `json:"currentPassword" validate:"required"`
	NewPassword     *string `json:"newPassword" validate:"required"`
}

type userRequestDto struct {
	Enabled     *bool               `json:"enabled" validate:"required"`
	Name        *string             `json:"name" validate:"required"`
	Username    *string             `json:"username" validate:"required"`
	Password    *string             `json:"password,omitempty"`
	Permissions *userPermissionsDto `json:"permissions" validate:"required"`
}

type userResponseDto struct {
	ID          uuid.UUID          `json:"id" validate:"required"`
	Enabled     bool               `json:"enabled" validate:"required"`
	Name        string             `json:"name" validate:"required"`
	Username    string             `json:"username" validate:"required"`
	Permissions userPermissionsDto `json:"permissions" validate:"required"`
}

type userPermissionsDto struct {
	Hosts        *string `json:"hosts" validate:"required"`
	Streams      *string `json:"streams" validate:"required"`
	Certificates *string `json:"certificates" validate:"required"`
	Logs         *string `json:"logs" validate:"required"`
	Integrations *string `json:"integrations" validate:"required"`
	AccessLists  *string `json:"accessLists" validate:"required"`
	Settings     *string `json:"settings" validate:"required"`
	Users        *string `json:"users" validate:"required"`
	NginxServer  *string `json:"nginxServer" validate:"required"`
	ExportData   *string `json:"exportData" validate:"required"`
}
