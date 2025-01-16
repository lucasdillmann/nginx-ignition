package settings

import "dillmann.com.br/nginx-ignition/core/host"

type Settings struct {
	Nginx                *NginxSettings
	LogRotation          *LogRotationSettings
	CertificateAutoRenew *CertificateAutoRenewSettings
	GlobalBindings       []*host.Binding
}

type NginxSettings struct {
	Logs                *NginxLogsSettings
	Timeouts            *NginxTimeoutsSettings
	WorkerProcesses     int
	WorkerConnections   int
	DefaultContentType  string
	ServerTokensEnabled bool
	MaximumBodySizeMb   int
	SendfileEnabled     bool
	GzipEnabled         bool
}

type LogRotationSettings struct {
	Enabled           bool
	MaximumLines      int
	IntervalUnit      TimeUnit
	IntervalUnitCount int
}

type CertificateAutoRenewSettings struct {
	Enabled           bool
	IntervalUnit      TimeUnit
	IntervalUnitCount int
}

type NginxTimeoutsSettings struct {
	Read      int
	Connect   int
	Send      int
	Keepalive int
}

type NginxLogsSettings struct {
	ServerLogsEnabled bool
	ServerLogsLevel   LogLevel
	AccessLogsEnabled bool
	ErrorLogsEnabled  bool
	ErrorLogsLevel    LogLevel
}

type LogLevel string

const (
	WarnLogLevel  = LogLevel("WARN")
	ErrorLogLevel = LogLevel("ERROR")
	CritLogLevel  = LogLevel("CRIT")
	AlertLogLevel = LogLevel("ALERT")
	EmergLogLevel = LogLevel("EMERG")
)

type TimeUnit string

const (
	MinutesTimeUnit = TimeUnit("MINUTES")
	HoursTimeUnit   = TimeUnit("HOURS")
	DaysTimeUnit    = TimeUnit("DAYS")
)
