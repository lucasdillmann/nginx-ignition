package settings

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type settingsDTO struct {
	Nginx                *nginxSettingsDTO                `json:"nginx"`
	LogRotation          *logRotationSettingsDTO          `json:"logRotation"`
	CertificateAutoRenew *certificateAutoRenewSettingsDTO `json:"certificateAutoRenew"`
	GlobalBindings       []bindingDTO                     `json:"globalBindings"`
}

type nginxSettingsDTO struct {
	Logs                *nginxLogsSettingsDTO     `json:"logs"`
	Timeouts            *nginxTimeoutsSettingsDTO `json:"timeouts"`
	Buffers             *nginxBuffersSettingsDTO  `json:"buffers"`
	Stats               *nginxStatsSettingsDTO    `json:"stats"`
	WorkerProcesses     *int                      `json:"workerProcesses"`
	WorkerConnections   *int                      `json:"workerConnections"`
	DefaultContentType  *string                   `json:"defaultContentType"`
	ServerTokensEnabled *bool                     `json:"serverTokensEnabled"`
	MaximumBodySizeMb   *int                      `json:"maximumBodySizeMb"`
	SendfileEnabled     *bool                     `json:"sendfileEnabled"`
	GzipEnabled         *bool                     `json:"gzipEnabled"`
	TCPNoDelayEnabled   *bool                     `json:"tcpNoDelayEnabled"`
	RuntimeUser         *string                   `json:"runtimeUser"`
	Custom              *string                   `json:"custom"`
}

type logRotationSettingsDTO struct {
	Enabled           *bool              `json:"enabled"`
	MaximumLines      *int               `json:"maximumLines"`
	IntervalUnit      *settings.TimeUnit `json:"intervalUnit"`
	IntervalUnitCount *int               `json:"intervalUnitCount"`
}

type certificateAutoRenewSettingsDTO struct {
	Enabled           *bool              `json:"enabled"`
	IntervalUnit      *settings.TimeUnit `json:"intervalUnit"`
	IntervalUnitCount *int               `json:"intervalUnitCount"`
}

type nginxTimeoutsSettingsDTO struct {
	Read       *int `json:"read"`
	Connect    *int `json:"connect"`
	Send       *int `json:"send"`
	Keepalive  *int `json:"keepalive"`
	ClientBody *int `json:"clientBody"`
}

type nginxBuffersSettingsDTO struct {
	ClientBodyKb      *int                `json:"clientBodyKb"`
	ClientHeaderKb    *int                `json:"clientHeaderKb"`
	LargeClientHeader *nginxBufferSizeDTO `json:"largeClientHeader"`
	Output            *nginxBufferSizeDTO `json:"output"`
}

type nginxBufferSizeDTO struct {
	SizeKb *int `json:"sizeKb"`
	Amount *int `json:"amount"`
}

type nginxLogsSettingsDTO struct {
	ServerLogsEnabled *bool              `json:"serverLogsEnabled"`
	ServerLogsLevel   *settings.LogLevel `json:"serverLogsLevel"`
	AccessLogsEnabled *bool              `json:"accessLogsEnabled"`
	ErrorLogsEnabled  *bool              `json:"errorLogsEnabled"`
	ErrorLogsLevel    *settings.LogLevel `json:"errorLogsLevel"`
}

type nginxStatsSettingsDTO struct {
	Enabled          *bool   `json:"enabled"`
	Persistent       *bool   `json:"persistent"`
	AllHosts         *bool   `json:"allHosts"`
	MaximumSizeMB    *int    `json:"maximumSizeMb"`
	DatabaseLocation *string `json:"databaseLocation"`
}

type bindingDTO struct {
	Type          *binding.Type `json:"type"`
	IP            *string       `json:"ip"`
	Port          *int          `json:"port"`
	CertificateID *uuid.UUID    `json:"certificateId"`
}
