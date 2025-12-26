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
	ID                uuid.UUID
	Enabled           bool
	DefaultServer     bool
	UseGlobalBindings bool
	DomainNames       []string
	Routes            []Route
	Bindings          []Binding
	VPNs              []VPN
	FeatureSet        FeatureSet
	AccessListID      *uuid.UUID
	CacheID           *uuid.UUID
}

type FeatureSet struct {
	WebsocketSupport    bool
	HTTP2Support        bool
	RedirectHTTPToHTTPS bool
}

type Route struct {
	ID           uuid.UUID
	Priority     int
	Enabled      bool
	Type         RouteType
	SourcePath   string
	TargetURI    *string
	RedirectCode *int
	AccessListID *uuid.UUID
	CacheID      *uuid.UUID
	Settings     RouteSettings
	Response     *RouteStaticResponse
	Integration  *RouteIntegrationConfig
	SourceCode   *RouteSourceCode
}

type RouteSourceCode struct {
	Language     CodeLanguage
	Contents     string
	MainFunction *string
}

type RouteSettings struct {
	IncludeForwardHeaders   bool
	ProxySSLServerName      bool
	KeepOriginalDomainName  bool
	DirectoryListingEnabled bool
	Custom                  *string
}

type RouteStaticResponse struct {
	StatusCode int
	Headers    map[string]string
	Payload    *string
}

type RouteIntegrationConfig struct {
	IntegrationID uuid.UUID
	OptionID      string
}

type Binding struct {
	ID            uuid.UUID
	Type          BindingType
	IP            string
	Port          int
	CertificateID *uuid.UUID
}

type VPN struct {
	VPNID uuid.UUID
	Name  string
	Host  *string
}
