package settings

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type nginxModel struct {
	bun.BaseModel `bun:"settings_nginx"`

	ID                  uuid.UUID `bun:"id,pk"`
	WorkerProcesses     int       `bun:"worker_processes"`
	WorkerConnections   int       `bun:"worker_connections"`
	ServerTokensEnabled bool      `bun:"server_tokens_enabled"`
	SendfileEnabled     bool      `bun:"sendfile_enabled"`
	GzipEnabled         bool      `bun:"gzip_enabled"`
	TcpNoDelayEnabled   bool      `bun:"tcp_nodelay_enabled"`
	DefaultContentType  string    `bun:"default_content_type"`
	MaximumBodySizeMb   int       `bun:"maximum_body_size_mb"`
	ReadTimeout         int       `bun:"read_timeout"`
	ConnectTimeout      int       `bun:"connect_timeout"`
	SendTimeout         int       `bun:"send_timeout"`
	KeepaliveTimeout    int       `bun:"keepalive_timeout"`
	ClientBodyTimeout   int       `bun:"client_body_timeout"`
	ServerLogsEnabled   bool      `bun:"server_logs_enabled"`
	ServerLogsLevel     string    `bun:"server_logs_level"`
	AccessLogsEnabled   bool      `bun:"access_logs_enabled"`
	ErrorLogsEnabled    bool      `bun:"error_logs_enabled"`
	ErrorLogsLevel      string    `bun:"error_logs_level"`
	RuntimeUser         string    `bun:"runtime_user"`
	Custom              *string   `bun:"custom"`
}

type logRotationModel struct {
	bun.BaseModel `bun:"settings_log_rotation"`

	ID                uuid.UUID `bun:"id,pk"`
	Enabled           bool      `bun:"enabled"`
	MaximumLines      int       `bun:"maximum_lines"`
	IntervalUnit      string    `bun:"interval_unit"`
	IntervalUnitCount int       `bun:"interval_unit_count"`
}

type certificateModel struct {
	bun.BaseModel `bun:"settings_certificate_auto_renew"`

	ID                uuid.UUID `bun:"id,pk"`
	Enabled           bool      `bun:"enabled"`
	IntervalUnit      string    `bun:"interval_unit"`
	IntervalUnitCount int       `bun:"interval_unit_count"`
}

type bindingModel struct {
	bun.BaseModel `bun:"settings_global_binding"`

	ID                  uuid.UUID  `bun:"id,pk"`
	Type                string     `bun:"type"`
	IP                  string     `bun:"ip"`
	Port                int        `bun:"port"`
	ServerCertificateID *uuid.UUID `bun:"server_certificate_id"`
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
