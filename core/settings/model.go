package settings

import (
	"dillmann.com.br/nginx-ignition/core/binding"
)

type Settings struct {
	Nginx                *NginxSettings
	LogRotation          *LogRotationSettings
	CertificateAutoRenew *CertificateAutoRenewSettings
	GlobalBindings       []binding.Binding
}

type NginxSettings struct {
	Timeouts            *NginxTimeoutsSettings
	Buffers             *NginxBuffersSettings
	Logs                *NginxLogsSettings
	Custom              *string
	RuntimeUser         string
	DefaultContentType  string
	WorkerProcesses     int
	WorkerConnections   int
	MaximumBodySizeMb   int
	ServerTokensEnabled bool
	TCPNoDelayEnabled   bool
	GzipEnabled         bool
	SendfileEnabled     bool
}

type LogRotationSettings struct {
	IntervalUnit      TimeUnit
	MaximumLines      int
	IntervalUnitCount int
	Enabled           bool
}

type CertificateAutoRenewSettings struct {
	IntervalUnit      TimeUnit
	IntervalUnitCount int
	Enabled           bool
}

type NginxTimeoutsSettings struct {
	Read       int
	Connect    int
	Send       int
	Keepalive  int
	ClientBody int
}

type NginxBuffersSettings struct {
	LargeClientHeader *NginxBufferSize
	Output            *NginxBufferSize
	ClientBodyKb      int
	ClientHeaderKb    int
}

type NginxBufferSize struct {
	SizeKb int
	Amount int
}

type NginxLogsSettings struct {
	ServerLogsLevel   LogLevel
	ErrorLogsLevel    LogLevel
	ServerLogsEnabled bool
	AccessLogsEnabled bool
	ErrorLogsEnabled  bool
}

type LogLevel string

const (
	WarnLogLevel  LogLevel = "WARN"
	ErrorLogLevel LogLevel = "ERROR"
	CritLogLevel  LogLevel = "CRIT"
	AlertLogLevel LogLevel = "ALERT"
	EmergLogLevel LogLevel = "EMERG"
)

type TimeUnit string

const (
	MinutesTimeUnit TimeUnit = "MINUTES"
	HoursTimeUnit   TimeUnit = "HOURS"
	DaysTimeUnit    TimeUnit = "DAYS"
)
