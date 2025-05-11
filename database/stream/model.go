package stream

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type streamModel struct {
	bun.BaseModel `bun:"stream"`

	ID               uuid.UUID `bun:"id,pk"`
	Enabled          bool      `bun:"enabled,notnull"`
	Name             string    `bun:"name,notnull"`
	BindingProtocol  string    `bun:"binding_protocol,notnull"`
	BindingAddress   string    `bun:"binding_address,notnull"`
	BindingPort      *int      `bun:"binding_port"`
	BackendProtocol  string    `bun:"backend_protocol,notnull"`
	BackendAddress   string    `bun:"backend_address,notnull"`
	BackendPort      *int      `bun:"backend_port"`
	UseProxyProtocol bool      `bun:"use_proxy_protocol,notnull"`
	SocketKeepAlive  bool      `bun:"socket_keep_alive,notnull"`
	TCPKeepAlive     bool      `bun:"tcp_keep_alive,notnull"`
	TCPNoDelay       bool      `bun:"tcp_no_delay,notnull"`
	TCPDeferred      bool      `bun:"tcp_deferred,notnull"`
}
