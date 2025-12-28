package settings

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type nginxModel struct {
	bun.BaseModel       `bun:"settings_nginx"`
	Custom              *string   `bun:"custom"`
	RuntimeUser         string    `bun:"runtime_user"`
	DefaultContentType  string    `bun:"default_content_type"`
	ErrorLogsLevel      string    `bun:"error_logs_level"`
	ServerLogsLevel     string    `bun:"server_logs_level"`
	APIAddress          string    `bun:"api_address"`
	KeepaliveTimeout    int       `bun:"keepalive_timeout"`
	WorkerConnections   int       `bun:"worker_connections"`
	ReadTimeout         int       `bun:"read_timeout"`
	ConnectTimeout      int       `bun:"connect_timeout"`
	SendTimeout         int       `bun:"send_timeout"`
	WorkerProcesses     int       `bun:"worker_processes"`
	ClientBodyTimeout   int       `bun:"client_body_timeout"`
	APIPort             int       `bun:"api_port"`
	MaximumBodySizeMb   int       `bun:"maximum_body_size_mb"`
	ID                  uuid.UUID `bun:"id,pk"`
	ServerLogsEnabled   bool      `bun:"server_logs_enabled"`
	ErrorLogsEnabled    bool      `bun:"error_logs_enabled"`
	TCPNoDelayEnabled   bool      `bun:"tcp_nodelay_enabled"`
	GzipEnabled         bool      `bun:"gzip_enabled"`
	ServerTokensEnabled bool      `bun:"server_tokens_enabled"`
	APIEnabled          bool      `bun:"api_enabled"`
	APIWriteEnabled     bool      `bun:"api_write_enabled"`
	AccessLogsEnabled   bool      `bun:"access_logs_enabled"`
	SendfileEnabled     bool      `bun:"sendfile_enabled"`
}

type logRotationModel struct {
	bun.BaseModel `bun:"settings_log_rotation"`

	IntervalUnit      string    `bun:"interval_unit"`
	MaximumLines      int       `bun:"maximum_lines"`
	IntervalUnitCount int       `bun:"interval_unit_count"`
	ID                uuid.UUID `bun:"id,pk"`
	Enabled           bool      `bun:"enabled"`
}

type certificateModel struct {
	bun.BaseModel `bun:"settings_certificate_auto_renew"`

	IntervalUnit      string    `bun:"interval_unit"`
	IntervalUnitCount int       `bun:"interval_unit_count"`
	ID                uuid.UUID `bun:"id,pk"`
	Enabled           bool      `bun:"enabled"`
}

type bindingModel struct {
	bun.BaseModel `bun:"settings_global_binding"`

	CertificateID *uuid.UUID `bun:"certificate_id"`
	Type          string     `bun:"type"`
	IP            string     `bun:"ip"`
	Port          int        `bun:"port"`
	ID            uuid.UUID  `bun:"id,pk"`
}

type buffersModel struct {
	bun.BaseModel `bun:"nginx_settings_buffers"`

	ID                      uuid.UUID `bun:"id,pk"`
	ClientBodyKb            int       `bun:"client_body_kb"`
	ClientHeaderKb          int       `bun:"client_header_kb"`
	LargeClientHeaderSizeKb int       `bun:"large_client_header_size_kb"`
	LargeClientHeaderAmount int       `bun:"large_client_header_amount"`
	OutputSizeKb            int       `bun:"output_size_kb"`
	OutputAmount            int       `bun:"output_amount"`
}
