package user

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type userModel struct {
	bun.BaseModel           `bun:"user"`
	StreamsAccessLevel      string    `bun:"streams_access_level,notnull"`
	IntegrationsAccessLevel string    `bun:"integrations_access_level,notnull"`
	Name                    string    `bun:"name,notnull"`
	Username                string    `bun:"username,notnull"`
	PasswordHash            string    `bun:"password_hash,notnull"`
	PasswordSalt            string    `bun:"password_salt,notnull"`
	HostsAccessLevel        string    `bun:"hosts_access_level,notnull"`
	CachesAccessLevel       string    `bun:"caches_access_level,notnull"`
	VPNsAccessLevel         string    `bun:"vpns_access_level,notnull"`
	CertificatesAccessLevel string    `bun:"certificates_access_level,notnull"`
	LogsAccessLevel         string    `bun:"logs_access_level,notnull"`
	AccessListsAccessLevel  string    `bun:"access_lists_access_level,notnull"`
	SettingsAccessLevel     string    `bun:"settings_access_level,notnull"`
	UsersAccessLevel        string    `bun:"users_access_level,notnull"`
	NginxServerAccessLevel  string    `bun:"nginx_server_access_level,notnull"`
	ExportDataAccessLevel   string    `bun:"export_data_access_level,notnull"`
	ID                      uuid.UUID `bun:"id,pk"`
	Enabled                 bool      `bun:"enabled,notnull"`
}
