package settings_api

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/google/uuid"
)

type SettingsDto struct {
	Nginx                *NginxSettingsDto                `json:"nginx" validate:"required"`
	LogRotation          *LogRotationSettingsDto          `json:"logRotation" validate:"required"`
	CertificateAutoRenew *CertificateAutoRenewSettingsDto `json:"certificateAutoRenew" validate:"required"`
	GlobalBindings       *[]BindingDto                    `json:"globalBindings" validate:"required"`
}

type NginxSettingsDto struct {
	Logs                *NginxLogsSettingsDto     `json:"logs" validate:"required"`
	Timeouts            *NginxTimeoutsSettingsDto `json:"timeouts" validate:"required"`
	WorkerProcesses     int                       `json:"workerProcesses" validate:"required"`
	WorkerConnections   int                       `json:"workerConnections" validate:"required"`
	DefaultContentType  string                    `json:"defaultContentType" validate:"required"`
	ServerTokensEnabled *bool                     `json:"serverTokensEnabled" validate:"required"`
	MaximumBodySizeMb   int                       `json:"maximumBodySizeMb" validate:"required"`
	SendfileEnabled     *bool                     `json:"sendfileEnabled" validate:"required"`
	GzipEnabled         *bool                     `json:"gzipEnabled" validate:"required"`
}

type LogRotationSettingsDto struct {
	Enabled           *bool             `json:"enabled" validate:"required"`
	MaximumLines      int               `json:"maximumLines" validate:"required"`
	IntervalUnit      settings.TimeUnit `json:"intervalUnit" validate:"required"`
	IntervalUnitCount int               `json:"intervalUnitCount" validate:"required"`
}

type CertificateAutoRenewSettingsDto struct {
	Enabled           *bool             `json:"enabled" validate:"required"`
	IntervalUnit      settings.TimeUnit `json:"intervalUnit" validate:"required"`
	IntervalUnitCount int               `json:"intervalUnitCount" validate:"required"`
}

type NginxTimeoutsSettingsDto struct {
	Read      int `json:"read" validate:"required"`
	Connect   int `json:"connect" validate:"required"`
	Send      int `json:"send" validate:"required"`
	Keepalive int `json:"keepalive" validate:"required"`
}

type NginxLogsSettingsDto struct {
	ServerLogsEnabled *bool             `json:"serverLogsEnabled" validate:"required"`
	ServerLogsLevel   settings.LogLevel `json:"serverLogsLevel" validate:"required"`
	AccessLogsEnabled *bool             `json:"accessLogsEnabled" validate:"required"`
	ErrorLogsEnabled  *bool             `json:"errorLogsEnabled" validate:"required"`
	ErrorLogsLevel    settings.LogLevel `json:"errorLogsLevel" validate:"required"`
}

type BindingDto struct {
	Type          host.BindingType `json:"type" validate:"required"`
	IP            string           `json:"ip" validate:"required,ip"`
	Port          int              `json:"port" validate:"required"`
	CertificateID *uuid.UUID       `json:"certificateId" validate:"required"`
}
