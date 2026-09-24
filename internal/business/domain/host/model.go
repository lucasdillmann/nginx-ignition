package host

import (
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/binding"
)

type CodeLanguage string

const (
	JavascriptCodeLanguage CodeLanguage = "JAVASCRIPT"
	LuaCodeLanguage        CodeLanguage = "LUA"
)

type RouteProtocol string

const (
	HTTP10RouteProtocol RouteProtocol = "HTTP_1_0"
	HTTP11RouteProtocol RouteProtocol = "HTTP_1_1"
	GRPCRouteProtocol   RouteProtocol = "GRPC"
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
	Bindings          []binding.Binding
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
	StatsEnabled        bool
}

type Route struct {
	Integration  *RouteIntegrationConfig
	TargetURI    *string
	AccessListID *uuid.UUID
	CacheID      *uuid.UUID
	Response     *RouteStaticResponse
	SourceCode   *RouteSourceCode
	RedirectCode *int
	Settings     RouteSettings
	Protocol     RouteProtocol
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
	IndexFile               *string
	IncludeForwardHeaders   bool
	ProxySSLServerName      bool
	IgnoreSSLErrors         bool
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
	UseHTTPS      bool
}

type VPN struct {
	Host          *string
	CertificateID *uuid.UUID
	Name          string
	VPNID         uuid.UUID
	EnableHTTPS   bool
}
