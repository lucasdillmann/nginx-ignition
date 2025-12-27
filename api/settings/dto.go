package settings

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type settingsDto struct {
	Nginx                *nginxSettingsDto                `json:"nginx"`
	LogRotation          *logRotationSettingsDto          `json:"logRotation"`
	CertificateAutoRenew *certificateAutoRenewSettingsDto `json:"certificateAutoRenew"`
	GlobalBindings       []bindingDto                     `json:"globalBindings"`
}

type nginxSettingsDto struct {
	Logs                *nginxLogsSettingsDto     `json:"logs"`
	Timeouts            *nginxTimeoutsSettingsDto `json:"timeouts"`
	Buffers             *nginxBuffersSettingsDto  `json:"buffers"`
	WorkerProcesses     *int                      `json:"workerProcesses"`
	WorkerConnections   *int                      `json:"workerConnections"`
	DefaultContentType  *string                   `json:"defaultContentType"`
	ServerTokensEnabled *bool                     `json:"serverTokensEnabled"`
	MaximumBodySizeMb   *int                      `json:"maximumBodySizeMb"`
	SendfileEnabled     *bool                     `json:"sendfileEnabled"`
	GzipEnabled         *bool                     `json:"gzipEnabled"`
	TcpNoDelayEnabled   *bool                     `json:"tcpNoDelayEnabled"`
	RuntimeUser         *string                   `json:"runtimeUser"`
	Custom              *string                   `json:"custom"`
}

type logRotationSettingsDto struct {
	Enabled           *bool              `json:"enabled"`
	MaximumLines      *int               `json:"maximumLines"`
	IntervalUnit      *settings.TimeUnit `json:"intervalUnit"`
	IntervalUnitCount *int               `json:"intervalUnitCount"`
}

type certificateAutoRenewSettingsDto struct {
	Enabled           *bool              `json:"enabled"`
	IntervalUnit      *settings.TimeUnit `json:"intervalUnit"`
	IntervalUnitCount *int               `json:"intervalUnitCount"`
}

type nginxTimeoutsSettingsDto struct {
	Read       *int `json:"read"`
	Connect    *int `json:"connect"`
	Send       *int `json:"send"`
	Keepalive  *int `json:"keepalive"`
	ClientBody *int `json:"clientBody"`
}

type nginxBuffersSettingsDto struct {
	ClientBodyKb      *int                `json:"clientBodyKb"`
	ClientHeaderKb    *int                `json:"clientHeaderKb"`
	LargeClientHeader *nginxBufferSizeDto `json:"largeClientHeader"`
	Output            *nginxBufferSizeDto `json:"output"`
}

type nginxBufferSizeDto struct {
	SizeKb *int `json:"sizeKb"`
	Amount *int `json:"amount"`
}

type nginxLogsSettingsDto struct {
	ServerLogsEnabled *bool              `json:"serverLogsEnabled"`
	ServerLogsLevel   *settings.LogLevel `json:"serverLogsLevel"`
	AccessLogsEnabled *bool              `json:"accessLogsEnabled"`
	ErrorLogsEnabled  *bool              `json:"errorLogsEnabled"`
	ErrorLogsLevel    *settings.LogLevel `json:"errorLogsLevel"`
}

type bindingDto struct {
	Type          *binding.Type `json:"type"`
	IP            *string       `json:"ip"`
	Port          *int          `json:"port"`
	CertificateID *uuid.UUID    `json:"certificateId"`
}
