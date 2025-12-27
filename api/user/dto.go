package user

import (
	"github.com/google/uuid"
)

type userLoginRequestDto struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type userLoginResponseDto struct {
	Token string `json:"token"`
}

type userOnboardingStatusResponseDto struct {
	Finished bool `json:"finished"`
}

type userPasswordUpdateRequestDto struct {
	CurrentPassword *string `json:"currentPassword"`
	NewPassword     *string `json:"newPassword"`
}

type userRequestDto struct {
	Enabled     *bool              `json:"enabled"`
	Name        *string            `json:"name"`
	Username    *string            `json:"username"`
	Password    *string            `json:"password,omitempty"`
	Permissions userPermissionsDto `json:"permissions"`
}

type userResponseDto struct {
	Permissions userPermissionsDto `json:"permissions"`
	Name        string             `json:"name"`
	Username    string             `json:"username"`
	ID          uuid.UUID          `json:"id"`
	Enabled     bool               `json:"enabled"`
}

type userPermissionsDto struct {
	Hosts        string `json:"hosts"`
	Streams      string `json:"streams"`
	Certificates string `json:"certificates"`
	Logs         string `json:"logs"`
	Integrations string `json:"integrations"`
	AccessLists  string `json:"accessLists"`
	Settings     string `json:"settings"`
	Users        string `json:"users"`
	NginxServer  string `json:"nginxServer"`
	ExportData   string `json:"exportData"`
	VPNs         string `json:"vpns"`
	Caches       string `json:"caches"`
}
