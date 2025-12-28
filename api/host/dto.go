package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/host"
)

type hostRequestDto struct {
	Enabled           *bool          `json:"enabled"`
	DefaultServer     *bool          `json:"defaultServer"`
	UseGlobalBindings *bool          `json:"useGlobalBindings"`
	FeatureSet        *featureSetDto `json:"featureSet"`
	AccessListID      *uuid.UUID     `json:"accessListId"`
	CacheID           *uuid.UUID     `json:"cacheId"`
	DomainNames       []string       `json:"domainNames"`
	Routes            []routeDto     `json:"routes"`
	Bindings          []bindingDto   `json:"bindings"`
	VPNs              []vpnDto       `json:"vpns"`
}

type routeDto struct {
	Priority     *int                  `json:"priority"`
	Enabled      *bool                 `json:"enabled"`
	Type         *host.RouteType       `json:"type"`
	SourcePath   *string               `json:"sourcePath"`
	Settings     *routeSettingsDto     `json:"settings"`
	TargetURI    *string               `json:"targetUri"`
	RedirectCode *int                  `json:"redirectCode"`
	Response     *staticResponseDto    `json:"response"`
	Integration  *integrationConfigDto `json:"integration"`
	AccessListID *uuid.UUID            `json:"accessListId"`
	CacheID      *uuid.UUID            `json:"cacheId"`
	SourceCode   *routeSourceCodeDto   `json:"sourceCode"`
}

type routeSourceCodeDto struct {
	Language     *host.CodeLanguage `json:"language"`
	Code         *string            `json:"code"`
	MainFunction *string            `json:"mainFunction"`
}

type routeSettingsDto struct {
	IncludeForwardHeaders   *bool   `json:"includeForwardHeaders"`
	ProxySslServerName      *bool   `json:"proxySslServerName"`
	KeepOriginalDomainName  *bool   `json:"keepOriginalDomainName"`
	DirectoryListingEnabled *bool   `json:"directoryListingEnabled"`
	Custom                  *string `json:"custom"`
}

type integrationConfigDto struct {
	IntegrationID *uuid.UUID `json:"integrationId"`
	OptionID      *string    `json:"optionId"`
}

type staticResponseDto struct {
	StatusCode *int               `json:"statusCode"`
	Payload    *string            `json:"payload"`
	Headers    *map[string]string `json:"headers"`
}

type featureSetDto struct {
	WebsocketsSupport   *bool `json:"websocketsSupport"`
	HTTP2Support        *bool `json:"http2Support"`
	RedirectHTTPToHTTPS *bool `json:"redirectHttpToHttps"`
}

type bindingDto struct {
	Type          *binding.Type `json:"type"`
	IP            *string       `json:"ip"`
	Port          *int          `json:"port"`
	CertificateID *uuid.UUID    `json:"certificateId"`
}

type vpnDto struct {
	VPNID *uuid.UUID `json:"vpnId"`
	Name  *string    `json:"name"`
	Host  *string    `json:"host"`
}

type hostResponseDto struct {
	ID                *uuid.UUID     `json:"id"`
	Enabled           *bool          `json:"enabled"`
	DefaultServer     *bool          `json:"defaultServer"`
	UseGlobalBindings *bool          `json:"useGlobalBindings"`
	FeatureSet        *featureSetDto `json:"featureSet"`
	AccessListID      *uuid.UUID     `json:"accessListId"`
	CacheID           *uuid.UUID     `json:"cacheId"`
	DomainNames       []string       `json:"domainNames"`
	Routes            []routeDto     `json:"routes"`
	Bindings          []bindingDto   `json:"bindings,omitempty"`
	GlobalBindings    []bindingDto   `json:"globalBindings,omitempty"`
	VPNs              []vpnDto       `json:"vpns"`
}
