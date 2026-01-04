package accesslist

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

type hostModel struct {
	bun.BaseModel `bun:"host"`

	DomainNames         []string  `bun:"domain_names,array"`
	ID                  uuid.UUID `bun:"id,pk"`
	AccessListID        uuid.UUID `bun:"access_list_id"`
	Enabled             bool      `bun:"enabled,notnull"`
	DefaultServer       bool      `bun:"default_server,notnull"`
	WebsocketSupport    bool      `bun:"websocket_support,notnull"`
	HTTP2Support        bool      `bun:"http2_support,notnull"`
	RedirectHTTPToHTTPS bool      `bun:"redirect_http_to_https,notnull"`
	UseGlobalBindings   bool      `bun:"use_global_bindings,notnull"`
}

func newAccessList() *accesslist.AccessList {
	return &accesslist.AccessList{
		ID:             uuid.New(),
		Name:           uuid.NewString(),
		Realm:          "Restricted Area",
		SatisfyAll:     false,
		DefaultOutcome: accesslist.AllowOutcome,
		Entries: []accesslist.Entry{
			{
				Outcome:       accesslist.DenyOutcome,
				SourceAddress: []string{"192.168.1.1"},
				Priority:      1,
			},
		},
		Credentials: []accesslist.Credentials{
			{
				Username: "user",
				Password: "password",
			},
		},
	}
}
