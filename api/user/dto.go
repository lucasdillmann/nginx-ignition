package user

import (
	"github.com/google/uuid"
)

type userLoginRequestDTO struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	TOTP     *string `json:"totp"`
}

type userLoginResponseDTO struct {
	Token string `json:"token"`
}

type userOnboardingStatusResponseDTO struct {
	Finished bool `json:"finished"`
}

type userPasswordUpdateRequestDTO struct {
	CurrentPassword *string `json:"currentPassword"`
	NewPassword     *string `json:"newPassword"`
}

type userRequestDTO struct {
	Enabled     *bool              `json:"enabled"`
	Name        *string            `json:"name"`
	Username    *string            `json:"username"`
	Password    *string            `json:"password,omitempty"`
	Permissions userPermissionsDTO `json:"permissions"`
}

type userResponseDTO struct {
	Permissions userPermissionsDTO `json:"permissions"`
	Name        string             `json:"name"`
	Username    string             `json:"username"`
	ID          uuid.UUID          `json:"id"`
	Enabled     bool               `json:"enabled"`
}

type userPermissionsDTO struct {
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
	TrafficStats string `json:"trafficStats"`
}
