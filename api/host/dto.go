package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/host"
)

type hostRequestDTO struct {
	Enabled           *bool          `json:"enabled"`
	DefaultServer     *bool          `json:"defaultServer"`
	UseGlobalBindings *bool          `json:"useGlobalBindings"`
	FeatureSet        *featureSetDTO `json:"featureSet"`
	AccessListID      *uuid.UUID     `json:"accessListId"`
	CacheID           *uuid.UUID     `json:"cacheId"`
	DomainNames       []string       `json:"domainNames"`
	Routes            []routeDTO     `json:"routes"`
	Bindings          []bindingDTO   `json:"bindings"`
	VPNs              []vpnDTO       `json:"vpns"`
}

type routeDTO struct {
	Priority     *int                  `json:"priority"`
	Enabled      *bool                 `json:"enabled"`
	Type         *host.RouteType       `json:"type"`
	SourcePath   *string               `json:"sourcePath"`
	Settings     *routeSettingsDTO     `json:"settings"`
	TargetURI    *string               `json:"targetUri"`
	RedirectCode *int                  `json:"redirectCode"`
	Response     *staticResponseDTO    `json:"response"`
	Integration  *integrationConfigDTO `json:"integration"`
	AccessListID *uuid.UUID            `json:"accessListId"`
	CacheID      *uuid.UUID            `json:"cacheId"`
	SourceCode   *routeSourceCodeDTO   `json:"sourceCode"`
}

type routeSourceCodeDTO struct {
	Language     *host.CodeLanguage `json:"language"`
	Code         *string            `json:"code"`
	MainFunction *string            `json:"mainFunction"`
}

type routeSettingsDTO struct {
	IncludeForwardHeaders   *bool   `json:"includeForwardHeaders"`
	ProxySslServerName      *bool   `json:"proxySslServerName"`
	KeepOriginalDomainName  *bool   `json:"keepOriginalDomainName"`
	DirectoryListingEnabled *bool   `json:"directoryListingEnabled"`
	Custom                  *string `json:"custom"`
}

type integrationConfigDTO struct {
	IntegrationID *uuid.UUID `json:"integrationId"`
	OptionID      *string    `json:"optionId"`
}

type staticResponseDTO struct {
	StatusCode *int               `json:"statusCode"`
	Payload    *string            `json:"payload"`
	Headers    *map[string]string `json:"headers"`
}

type featureSetDTO struct {
	WebsocketsSupport   *bool `json:"websocketsSupport"`
	HTTP2Support        *bool `json:"http2Support"`
	RedirectHTTPToHTTPS *bool `json:"redirectHttpToHttps"`
}

type bindingDTO struct {
	Type          *binding.Type `json:"type"`
	IP            *string       `json:"ip"`
	Port          *int          `json:"port"`
	CertificateID *uuid.UUID    `json:"certificateId"`
}

type vpnDTO struct {
	VPNID *uuid.UUID `json:"vpnId"`
	Name  *string    `json:"name"`
	Host  *string    `json:"host"`
}

type hostResponseDTO struct {
	ID                *uuid.UUID     `json:"id"`
	Enabled           *bool          `json:"enabled"`
	DefaultServer     *bool          `json:"defaultServer"`
	UseGlobalBindings *bool          `json:"useGlobalBindings"`
	FeatureSet        *featureSetDTO `json:"featureSet"`
	AccessListID      *uuid.UUID     `json:"accessListId"`
	CacheID           *uuid.UUID     `json:"cacheId"`
	DomainNames       []string       `json:"domainNames"`
	Routes            []routeDTO     `json:"routes"`
	Bindings          []bindingDTO   `json:"bindings,omitempty"`
	GlobalBindings    []bindingDTO   `json:"globalBindings,omitempty"`
	VPNs              []vpnDTO       `json:"vpns"`
}
