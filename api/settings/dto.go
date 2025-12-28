package settings

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type settingsDTO struct {
	GlobalBindings       []bindingDTO                    `json:"globalBindings"`
	CertificateAutoRenew certificateAutoRenewSettingsDTO `json:"certificateAutoRenew"`
	LogRotation          logRotationSettingsDTO          `json:"logRotation"`
	Nginx                nginxSettingsDTO                `json:"nginx"`
}

type nginxSettingsDTO struct {
	Custom              *string                  `json:"custom"`
	Logs                nginxLogsSettingsDTO     `json:"logs"`
	DefaultContentType  string                   `json:"defaultContentType"`
	RuntimeUser         string                   `json:"runtimeUser"`
	API                 nginxAPIDTO              `json:"api"`
	Buffers             nginxBuffersSettingsDTO  `json:"buffers"`
	Timeouts            nginxTimeoutsSettingsDTO `json:"timeouts"`
	WorkerConnections   int                      `json:"workerConnections"`
	WorkerProcesses     int                      `json:"workerProcesses"`
	MaximumBodySizeMb   int                      `json:"maximumBodySizeMb"`
	ServerTokensEnabled bool                     `json:"serverTokensEnabled"`
	SendfileEnabled     bool                     `json:"sendfileEnabled"`
	GzipEnabled         bool                     `json:"gzipEnabled"`
	TCPNoDelayEnabled   bool                     `json:"tcpNoDelayEnabled"`
}

type logRotationSettingsDTO struct {
	IntervalUnit      settings.TimeUnit `json:"intervalUnit"`
	MaximumLines      int               `json:"maximumLines"`
	IntervalUnitCount int               `json:"intervalUnitCount"`
	Enabled           bool              `json:"enabled"`
}

type certificateAutoRenewSettingsDTO struct {
	IntervalUnit      settings.TimeUnit `json:"intervalUnit"`
	IntervalUnitCount int               `json:"intervalUnitCount"`
	Enabled           bool              `json:"enabled"`
}

type nginxTimeoutsSettingsDTO struct {
	Read       int `json:"read"`
	Connect    int `json:"connect"`
	Send       int `json:"send"`
	Keepalive  int `json:"keepalive"`
	ClientBody int `json:"clientBody"`
}

type nginxBuffersSettingsDTO struct {
	ClientBodyKb      int                `json:"clientBodyKb"`
	ClientHeaderKb    int                `json:"clientHeaderKb"`
	LargeClientHeader nginxBufferSizeDTO `json:"largeClientHeader"`
	Output            nginxBufferSizeDTO `json:"output"`
}

type nginxBufferSizeDTO struct {
	SizeKb int `json:"sizeKb"`
	Amount int `json:"amount"`
}

type nginxLogsSettingsDTO struct {
	ServerLogsLevel   settings.LogLevel `json:"serverLogsLevel"`
	ErrorLogsLevel    settings.LogLevel `json:"errorLogsLevel"`
	ServerLogsEnabled bool              `json:"serverLogsEnabled"`
	AccessLogsEnabled bool              `json:"accessLogsEnabled"`
	ErrorLogsEnabled  bool              `json:"errorLogsEnabled"`
}

type nginxAPIDTO struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	Enabled      bool   `json:"enabled"`
	WriteEnabled bool   `json:"writeEnabled"`
}

type bindingDTO struct {
	CertificateID *uuid.UUID   `json:"certificateId"`
	Type          binding.Type `json:"type"`
	IP            string       `json:"ip"`
	Port          int          `json:"port"`
}
