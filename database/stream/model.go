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
	Type             string    `bun:"type,notnull"`
	BindingProtocol  string    `bun:"binding_protocol,notnull"`
	BindingAddress   string    `bun:"binding_address,notnull"`
	BindingPort      *int      `bun:"binding_port"`
	UseProxyProtocol bool      `bun:"use_proxy_protocol,notnull"`
	SocketKeepAlive  bool      `bun:"socket_keep_alive,notnull"`
	TCPKeepAlive     bool      `bun:"tcp_keep_alive,notnull"`
	TCPNoDelay       bool      `bun:"tcp_no_delay,notnull"`
	TCPDeferred      bool      `bun:"tcp_deferred,notnull"`
}

type streamRouteModel struct {
	bun.BaseModel `bun:"stream_route"`

	ID         uuid.UUID `bun:"id,pk"`
	StreamID   uuid.UUID `bun:"stream_id,notnull"`
	DomainName string    `bun:"domain_name,notnull"`
}

type streamBackendModel struct {
	bun.BaseModel `bun:"stream_backend"`

	ID            uuid.UUID  `bun:"id,pk"`
	StreamID      *uuid.UUID `bun:"stream_id"`
	StreamRouteID *uuid.UUID `bun:"stream_route_id"`
	Protocol      string     `bun:"protocol,notnull"`
	Address       string     `bun:"address,notnull"`
	Port          *int       `bun:"port"`
	Weight        *int       `bun:"weight"`
	MaxFailures   *int       `bun:"max_failures"`
	OpenSeconds   *int       `bun:"open_seconds"`
}
