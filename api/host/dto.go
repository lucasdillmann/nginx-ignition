package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/google/uuid"
)

type hostRequestDto struct {
	Enabled           *bool          `json:"enabled" validate:"required"`
	DefaultServer     *bool          `json:"defaultServer" validate:"required"`
	UseGlobalBindings *bool          `json:"useGlobalBindings" validate:"required"`
	DomainNames       []*string      `json:"domainNames"`
	Routes            []*routeDto    `json:"routes"`
	Bindings          []*bindingDto  `json:"bindings"`
	FeatureSet        *featureSetDto `json:"featureSet" validate:"required"`
	AccessListId      *uuid.UUID     `json:"accessListId"`
}

type routeDto struct {
	Priority     *int                  `json:"priority" validate:"required"`
	Enabled      *bool                 `json:"enabled" validate:"required"`
	Type         *host.RouteType       `json:"type" validate:"required"`
	SourcePath   *string               `json:"sourcePath" validate:"required"`
	Settings     *routeSettingsDto     `json:"settings" validate:"required"`
	TargetUri    *string               `json:"targetUri"`
	RedirectCode *int                  `json:"redirectCode"`
	Response     *staticResponseDto    `json:"response"`
	Integration  *integrationConfigDto `json:"integration"`
	AccessListId *uuid.UUID            `json:"accessListId"`
	SourceCode   *routeSourceCodeDto   `json:"sourceCode"`
}

type routeSourceCodeDto struct {
	Language     *host.CodeLanguage `json:"language" validate:"required"`
	Code         *string            `json:"code" validate:"required"`
	MainFunction *string            `json:"mainFunction"`
}

type routeSettingsDto struct {
	IncludeForwardHeaders   *bool   `json:"includeForwardHeaders" validate:"required"`
	ProxySslServerName      *bool   `json:"proxySslServerName" validate:"required"`
	KeepOriginalDomainName  *bool   `json:"keepOriginalDomainName" validate:"required"`
	DirectoryListingEnabled *bool   `json:"directoryListingEnabled" validate:"required"`
	Custom                  *string `json:"custom"`
}

type integrationConfigDto struct {
	IntegrationId *string `json:"integrationId" validate:"required"`
	OptionId      *string `json:"optionId" validate:"required"`
}

type staticResponseDto struct {
	StatusCode *int               `json:"statusCode" validate:"required"`
	Payload    *string            `json:"payload"`
	Headers    *map[string]string `json:"headers"`
}

type featureSetDto struct {
	WebsocketsSupport   *bool `json:"websocketsSupport" validate:"required"`
	Http2Support        *bool `json:"http2Support" validate:"required"`
	RedirectHttpToHttps *bool `json:"redirectHttpToHttps" validate:"required"`
}

type bindingDto struct {
	Type          *host.BindingType `json:"type" validate:"required"`
	Ip            *string           `json:"ip" validate:"required"`
	Port          *int              `json:"port" validate:"required"`
	CertificateId *uuid.UUID        `json:"certificateId"`
}

type hostResponseDto struct {
	ID                *uuid.UUID     `json:"id" validate:"required"`
	Enabled           *bool          `json:"enabled" validate:"required"`
	DefaultServer     *bool          `json:"defaultServer" validate:"required"`
	UseGlobalBindings *bool          `json:"useGlobalBindings" validate:"required"`
	DomainNames       []*string      `json:"domainNames"`
	Routes            []*routeDto    `json:"routes" validate:"required"`
	Bindings          []*bindingDto  `json:"bindings" validate:"required"`
	FeatureSet        *featureSetDto `json:"featureSet" validate:"required"`
	AccessListId      *uuid.UUID     `json:"accessListId"`
}
