package host

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/host"
)

type hostRequestDto struct {
	Enabled           *bool          `json:"enabled"`
	DefaultServer     *bool          `json:"defaultServer"`
	UseGlobalBindings *bool          `json:"useGlobalBindings"`
	DomainNames       []*string      `json:"domainNames"`
	Routes            []*routeDto    `json:"routes"`
	Bindings          []*bindingDto  `json:"bindings"`
	VPNs              []*vpnDto      `json:"vpns"`
	FeatureSet        *featureSetDto `json:"featureSet"`
	AccessListID      *uuid.UUID     `json:"accessListId"`
}

type routeDto struct {
	Priority     *int                  `json:"priority"`
	Enabled      *bool                 `json:"enabled"`
	Type         *host.RouteType       `json:"type"`
	SourcePath   *string               `json:"sourcePath"`
	Settings     *routeSettingsDto     `json:"settings"`
	TargetUri    *string               `json:"targetUri"`
	RedirectCode *int                  `json:"redirectCode"`
	Response     *staticResponseDto    `json:"response"`
	Integration  *integrationConfigDto `json:"integration"`
	AccessListID *uuid.UUID            `json:"accessListId"`
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
	Http2Support        *bool `json:"http2Support"`
	RedirectHttpToHttps *bool `json:"redirectHttpToHttps"`
}

type bindingDto struct {
	Type                *host.BindingType `json:"type"`
	IP                  *string           `json:"ip"`
	Port                *int              `json:"port"`
	ServerCertificateID *uuid.UUID        `json:"serverCertificateId"`
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
	DomainNames       []*string      `json:"domainNames"`
	Routes            []*routeDto    `json:"routes"`
	Bindings          []*bindingDto  `json:"bindings,omitempty"`
	GlobalBindings    *[]*bindingDto `json:"globalBindings,omitempty"`
	VPNs              []*vpnDto      `json:"vpns"`
	FeatureSet        *featureSetDto `json:"featureSet"`
	AccessListID      *uuid.UUID     `json:"accessListId"`
}
