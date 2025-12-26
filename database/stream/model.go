package stream

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type streamModel struct {
	bun.BaseModel    `bun:"stream"`
	BindingPort      *int      `bun:"binding_port"`
	BindingAddress   string    `bun:"binding_address,notnull"`
	Name             string    `bun:"name,notnull"`
	Type             string    `bun:"type,notnull"`
	BindingProtocol  string    `bun:"binding_protocol,notnull"`
	ID               uuid.UUID `bun:"id,pk"`
	Enabled          bool      `bun:"enabled,notnull"`
	UseProxyProtocol bool      `bun:"use_proxy_protocol,notnull"`
	SocketKeepAlive  bool      `bun:"socket_keep_alive,notnull"`
	TCPKeepAlive     bool      `bun:"tcp_keep_alive,notnull"`
	TCPNoDelay       bool      `bun:"tcp_no_delay,notnull"`
	TCPDeferred      bool      `bun:"tcp_deferred,notnull"`
}

type streamRouteModel struct {
	bun.BaseModel `bun:"stream_route"`
	DomainNames   []string  `bun:"domain_names,array,notnull"`
	ID            uuid.UUID `bun:"id,pk"`
	StreamID      uuid.UUID `bun:"stream_id,notnull"`
}

type streamBackendModel struct {
	bun.BaseModel `bun:"stream_backend"`
	StreamID      *uuid.UUID `bun:"stream_id"`
	StreamRouteID *uuid.UUID `bun:"stream_route_id"`
	Port          *int       `bun:"port"`
	Weight        *int       `bun:"weight"`
	MaxFailures   *int       `bun:"max_failures"`
	OpenSeconds   *int       `bun:"open_seconds"`
	Protocol      string     `bun:"protocol,notnull"`
	Address       string     `bun:"address,notnull"`
	ID            uuid.UUID  `bun:"id,pk"`
}
