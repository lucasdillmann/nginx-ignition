package host

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type hostModel struct {
	bun.BaseModel `bun:"host"`

	AccessListID        *uuid.UUID         `bun:"access_list_id"`
	CacheID             *uuid.UUID         `bun:"cache_id"`
	VPNs                []hostVpnModel     `bun:"rel:has-many,join:id=host_id"`
	DomainNames         []string           `bun:"domain_names,array"`
	Routes              []hostRouteModel   `bun:"rel:has-many,join:id=host_id"`
	Bindings            []hostBindingModel `bun:"rel:has-many,join:id=host_id"`
	ID                  uuid.UUID          `bun:"id,pk"`
	DefaultServer       bool               `bun:"default_server,notnull"`
	UseGlobalBindings   bool               `bun:"use_global_bindings,notnull"`
	RedirectHTTPToHTTPS bool               `bun:"redirect_http_to_https,notnull"`
	HTTP2Support        bool               `bun:"http2_support,notnull"`
	WebsocketSupport    bool               `bun:"websocket_support,notnull"`
	Enabled             bool               `bun:"enabled,notnull"`
}

type hostBindingModel struct {
	bun.BaseModel `bun:"host_binding"`

	CertificateID *uuid.UUID `bun:"certificate_id"`
	Type          string     `bun:"type,notnull"`
	IP            string     `bun:"ip,notnull"`
	Port          int        `bun:"port,notnull"`
	ID            uuid.UUID  `bun:"id,pk"`
	HostID        uuid.UUID  `bun:"host_id,notnull"`
}

type hostVpnModel struct {
	bun.BaseModel `bun:"host_vpn"`

	Host   *string   `bun:"host"`
	Name   string    `bun:"name,notnull"`
	HostID uuid.UUID `bun:"host_id,notnull"`
	VPNID  uuid.UUID `bun:"vpn_id,notnull"`
}

type hostRouteModel struct {
	bun.BaseModel `bun:"host_route"`

	IntegrationID           *uuid.UUID `bun:"integration_id"`
	StaticResponsePayload   *string    `bun:"static_response_payload"`
	CodeMainFunction        *string    `bun:"code_main_function"`
	CodeContents            *string    `bun:"code_contents"`
	CodeLanguage            *string    `bun:"code_language"`
	TargetURI               *string    `bun:"target_uri"`
	CustomSettings          *string    `bun:"custom_settings"`
	IntegrationOptionID     *string    `bun:"integration_option_id"`
	CacheID                 *uuid.UUID `bun:"cache_id"`
	StaticResponseHeaders   *string    `bun:"static_response_headers"`
	AccessListID            *uuid.UUID `bun:"access_list_id"`
	RedirectCode            *int       `bun:"redirect_code"`
	StaticResponseCode      *int       `bun:"static_response_code"`
	SourcePath              string     `bun:"source_path,notnull"`
	Type                    string     `bun:"type,notnull"`
	Priority                int        `bun:"priority,notnull"`
	HostID                  uuid.UUID  `bun:"host_id,notnull"`
	ID                      uuid.UUID  `bun:"id,pk"`
	IncludeForwardHeaders   bool       `bun:"include_forward_headers,notnull"`
	ProxySSLServerName      bool       `bun:"proxy_ssl_server_name,notnull"`
	KeepOriginalDomainName  bool       `bun:"keep_original_domain_name,notnull"`
	DirectoryListingEnabled bool       `bun:"directory_listing_enabled,notnull"`
	Enabled                 bool       `bun:"enabled,notnull"`
}
