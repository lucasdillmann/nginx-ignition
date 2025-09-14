package settings

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type settingsDto struct {
	Nginx                *nginxSettingsDto                `json:"nginx" validate:"required"`
	LogRotation          *logRotationSettingsDto          `json:"logRotation" validate:"required"`
	CertificateAutoRenew *certificateAutoRenewSettingsDto `json:"certificateAutoRenew" validate:"required"`
	GlobalBindings       []*bindingDto                    `json:"globalBindings" validate:"required"`
}

type nginxSettingsDto struct {
	Logs                *nginxLogsSettingsDto     `json:"logs" validate:"required"`
	Timeouts            *nginxTimeoutsSettingsDto `json:"timeouts" validate:"required"`
	WorkerProcesses     *int                      `json:"workerProcesses" validate:"required"`
	WorkerConnections   *int                      `json:"workerConnections" validate:"required"`
	DefaultContentType  *string                   `json:"defaultContentType" validate:"required"`
	ServerTokensEnabled *bool                     `json:"serverTokensEnabled" validate:"required"`
	MaximumBodySizeMb   *int                      `json:"maximumBodySizeMb" validate:"required"`
	SendfileEnabled     *bool                     `json:"sendfileEnabled" validate:"required"`
	GzipEnabled         *bool                     `json:"gzipEnabled" validate:"required"`
	RuntimeUser         *string                   `json:"runtimeUser" validate:"required"`
}

type logRotationSettingsDto struct {
	Enabled           *bool              `json:"enabled" validate:"required"`
	MaximumLines      *int               `json:"maximumLines" validate:"required"`
	IntervalUnit      *settings.TimeUnit `json:"intervalUnit" validate:"required"`
	IntervalUnitCount *int               `json:"intervalUnitCount" validate:"required"`
}

type certificateAutoRenewSettingsDto struct {
	Enabled           *bool              `json:"enabled" validate:"required"`
	IntervalUnit      *settings.TimeUnit `json:"intervalUnit" validate:"required"`
	IntervalUnitCount *int               `json:"intervalUnitCount" validate:"required"`
}

type nginxTimeoutsSettingsDto struct {
	Read      *int `json:"read" validate:"required"`
	Connect   *int `json:"connect" validate:"required"`
	Send      *int `json:"send" validate:"required"`
	Keepalive *int `json:"keepalive" validate:"required"`
}

type nginxLogsSettingsDto struct {
	ServerLogsEnabled *bool              `json:"serverLogsEnabled" validate:"required"`
	ServerLogsLevel   *settings.LogLevel `json:"serverLogsLevel" validate:"required"`
	AccessLogsEnabled *bool              `json:"accessLogsEnabled" validate:"required"`
	ErrorLogsEnabled  *bool              `json:"errorLogsEnabled" validate:"required"`
	ErrorLogsLevel    *settings.LogLevel `json:"errorLogsLevel" validate:"required"`
}

type bindingDto struct {
	Type          *host.BindingType `json:"type" validate:"required"`
	IP            *string           `json:"ip" validate:"required"`
	Port          *int              `json:"port" validate:"required"`
	CertificateID *uuid.UUID        `json:"certificateId"`
}
