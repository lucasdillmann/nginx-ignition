package host

import (
	"github.com/google/uuid"
)

type BindingType string

const (
	HttpBindingType  BindingType = "HTTP"
	HttpsBindingType BindingType = "HTTPS"
)

type CodeLanguage string

const (
	JavascriptCodeLanguage CodeLanguage = "JAVASCRIPT"
	LuaCodeLanguage        CodeLanguage = "LUA"
)

type RouteType string

const (
	ProxyRouteType          RouteType = "PROXY"
	RedirectRouteType       RouteType = "REDIRECT"
	StaticResponseRouteType RouteType = "STATIC_RESPONSE"
	IntegrationRouteType    RouteType = "INTEGRATION"
	ExecuteCodeRouteType    RouteType = "EXECUTE_CODE"
	StaticFilesRouteType    RouteType = "STATIC_FILES"
)

type Host struct {
	AccessListID      *uuid.UUID
	CacheID           *uuid.UUID
	DomainNames       []string
	Routes            []Route
	Bindings          []Binding
	VPNs              []VPN
	ID                uuid.UUID
	FeatureSet        FeatureSet
	Enabled           bool
	DefaultServer     bool
	UseGlobalBindings bool
}

type FeatureSet struct {
	WebsocketSupport    bool
	HTTP2Support        bool
	RedirectHTTPToHTTPS bool
}

type Route struct {
	Settings     RouteSettings
	RedirectCode *int
	TargetURI    *string
	AccessListID *uuid.UUID
	CacheID      *uuid.UUID
	Response     *RouteStaticResponse
	Integration  *RouteIntegrationConfig
	SourceCode   *RouteSourceCode
	Type         RouteType
	SourcePath   string
	Priority     int
	ID           uuid.UUID
	Enabled      bool
}

type RouteSourceCode struct {
	MainFunction *string
	Language     CodeLanguage
	Contents     string
}

type RouteSettings struct {
	Custom                  *string
	IncludeForwardHeaders   bool
	ProxySSLServerName      bool
	KeepOriginalDomainName  bool
	DirectoryListingEnabled bool
}

type RouteStaticResponse struct {
	Headers    map[string]string
	Payload    *string
	StatusCode int
}

type RouteIntegrationConfig struct {
	OptionID      string
	IntegrationID uuid.UUID
}

type Binding struct {
	CertificateID *uuid.UUID
	Type          BindingType
	IP            string
	Port          int
	ID            uuid.UUID
}

type VPN struct {
	Host  *string
	Name  string
	VPNID uuid.UUID
}
