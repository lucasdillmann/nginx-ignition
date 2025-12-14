package host

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type hostModel struct {
	bun.BaseModel `bun:"host"`

	ID                  uuid.UUID           `bun:"id,pk"`
	Enabled             bool                `bun:"enabled,notnull"`
	DefaultServer       bool                `bun:"default_server,notnull"`
	DomainNames         []string            `bun:"domain_names,array"`
	WebsocketSupport    bool                `bun:"websocket_support,notnull"`
	HTTP2Support        bool                `bun:"http2_support,notnull"`
	RedirectHTTPToHTTPS bool                `bun:"redirect_http_to_https,notnull"`
	UseGlobalBindings   bool                `bun:"use_global_bindings,notnull"`
	AccessListID        *uuid.UUID          `bun:"access_list_id"`
	ClientCertificateID *uuid.UUID          `bun:"client_certificate_id"`
	Bindings            []*hostBindingModel `bun:"rel:has-many,join:id=host_id"`
	Routes              []*hostRouteModel   `bun:"rel:has-many,join:id=host_id"`
	VPNs                []*hostVpnModel     `bun:"rel:has-many,join:id=host_id"`
}

type hostBindingModel struct {
	bun.BaseModel `bun:"host_binding"`

	ID                  uuid.UUID  `bun:"id,pk"`
	HostID              uuid.UUID  `bun:"host_id,notnull"`
	Type                string     `bun:"type,notnull"`
	IP                  string     `bun:"ip,notnull"`
	Port                int        `bun:"port,notnull"`
	ServerCertificateID *uuid.UUID `bun:"server_certificate_id"`
}

type hostVpnModel struct {
	bun.BaseModel `bun:"host_vpn"`

	HostID uuid.UUID `bun:"host_id,notnull"`
	VPNID  uuid.UUID `bun:"vpn_id,notnull"`
	Name   string    `bun:"name,notnull"`
	Host   *string   `bun:"host"`
}

type hostRouteModel struct {
	bun.BaseModel `bun:"host_route"`

	ID                      uuid.UUID  `bun:"id,pk"`
	HostID                  uuid.UUID  `bun:"host_id,notnull"`
	Priority                int        `bun:"priority,notnull"`
	Type                    string     `bun:"type,notnull"`
	SourcePath              string     `bun:"source_path,notnull"`
	TargetURI               *string    `bun:"target_uri"`
	CustomSettings          *string    `bun:"custom_settings"`
	StaticResponseCode      *int       `bun:"static_response_code"`
	StaticResponsePayload   *string    `bun:"static_response_payload"`
	StaticResponseHeaders   *string    `bun:"static_response_headers"`
	RedirectCode            *int       `bun:"redirect_code"`
	IntegrationID           *uuid.UUID `bun:"integration_id"`
	IntegrationOptionID     *string    `bun:"integration_option_id"`
	IncludeForwardHeaders   bool       `bun:"include_forward_headers,notnull"`
	ProxySSLServerName      bool       `bun:"proxy_ssl_server_name,notnull"`
	KeepOriginalDomainName  bool       `bun:"keep_original_domain_name,notnull"`
	DirectoryListingEnabled bool       `bun:"directory_listing_enabled,notnull"`
	AccessListID            *uuid.UUID `bun:"access_list_id"`
	ClientCertificateID     *uuid.UUID `bun:"client_certificate_id"`
	CodeLanguage            *string    `bun:"code_language"`
	CodeContents            *string    `bun:"code_contents"`
	CodeMainFunction        *string    `bun:"code_main_function"`
	Enabled                 bool       `bun:"enabled,notnull"`
}
