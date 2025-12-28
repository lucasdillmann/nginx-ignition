package settings

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type settingsDto struct {
	GlobalBindings       []bindingDto                    `json:"globalBindings"`
	CertificateAutoRenew certificateAutoRenewSettingsDto `json:"certificateAutoRenew"`
	LogRotation          logRotationSettingsDto          `json:"logRotation"`
	Nginx                nginxSettingsDto                `json:"nginx"`
}

type nginxSettingsDto struct {
	Custom              *string                  `json:"custom"`
	Logs                nginxLogsSettingsDto     `json:"logs"`
	DefaultContentType  string                   `json:"defaultContentType"`
	RuntimeUser         string                   `json:"runtimeUser"`
	API                 nginxApiDto              `json:"api"`
	Buffers             nginxBuffersSettingsDto  `json:"buffers"`
	Timeouts            nginxTimeoutsSettingsDto `json:"timeouts"`
	WorkerConnections   int                      `json:"workerConnections"`
	WorkerProcesses     int                      `json:"workerProcesses"`
	MaximumBodySizeMb   int                      `json:"maximumBodySizeMb"`
	ServerTokensEnabled bool                     `json:"serverTokensEnabled"`
	SendfileEnabled     bool                     `json:"sendfileEnabled"`
	GzipEnabled         bool                     `json:"gzipEnabled"`
	TcpNoDelayEnabled   bool                     `json:"tcpNoDelayEnabled"`
}

type logRotationSettingsDto struct {
	IntervalUnit      settings.TimeUnit `json:"intervalUnit"`
	MaximumLines      int               `json:"maximumLines"`
	IntervalUnitCount int               `json:"intervalUnitCount"`
	Enabled           bool              `json:"enabled"`
}

type certificateAutoRenewSettingsDto struct {
	IntervalUnit      settings.TimeUnit `json:"intervalUnit"`
	IntervalUnitCount int               `json:"intervalUnitCount"`
	Enabled           bool              `json:"enabled"`
}

type nginxTimeoutsSettingsDto struct {
	Read       int `json:"read"`
	Connect    int `json:"connect"`
	Send       int `json:"send"`
	Keepalive  int `json:"keepalive"`
	ClientBody int `json:"clientBody"`
}

type nginxBuffersSettingsDto struct {
	ClientBodyKb      int                `json:"clientBodyKb"`
	ClientHeaderKb    int                `json:"clientHeaderKb"`
	LargeClientHeader nginxBufferSizeDto `json:"largeClientHeader"`
	Output            nginxBufferSizeDto `json:"output"`
}

type nginxBufferSizeDto struct {
	SizeKb int `json:"sizeKb"`
	Amount int `json:"amount"`
}

type nginxLogsSettingsDto struct {
	ServerLogsLevel   settings.LogLevel `json:"serverLogsLevel"`
	ErrorLogsLevel    settings.LogLevel `json:"errorLogsLevel"`
	ServerLogsEnabled bool              `json:"serverLogsEnabled"`
	AccessLogsEnabled bool              `json:"accessLogsEnabled"`
	ErrorLogsEnabled  bool              `json:"errorLogsEnabled"`
}

type nginxApiDto struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	Enabled      bool   `json:"enabled"`
	WriteEnabled bool   `json:"writeEnabled"`
}

type bindingDto struct {
	CertificateID *uuid.UUID   `json:"certificateId"`
	Type          binding.Type `json:"type"`
	IP            string       `json:"ip"`
	Port          int          `json:"port"`
}
