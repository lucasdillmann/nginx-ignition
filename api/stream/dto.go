package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/stream"
)

type streamRequestDto struct {
	Enabled        *bool          `json:"enabled" validate:"required"`
	Name           *string        `json:"name" validate:"required"`
	Type           *string        `json:"type" validate:"required"`
	FeatureSet     *featureSetDto `json:"featureSet" validate:"required"`
	DefaultBackend *backendDto    `json:"defaultBackend" validate:"required"`
	Binding        *addressDto    `json:"binding" validate:"required"`
	Routes         *[]routeDto    `json:"routes"`
}

type featureSetDto struct {
	UseProxyProtocol *bool `json:"useProxyProtocol"`
	SocketKeepAlive  *bool `json:"socketKeepAlive"`
	TCPKeepAlive     *bool `json:"tcpKeepAlive"`
	TCPNoDelay       *bool `json:"tcpNoDelay"`
	TCPDeferred      *bool `json:"tcpDeferred"`
}

type addressDto struct {
	Protocol stream.Protocol `json:"protocol"`
	Address  *string         `json:"address"`
	Port     *int            `json:"port"`
}

type backendDto struct {
	Weight         *int               `json:"weight"`
	Target         *addressDto        `json:"Target" validate:"required"`
	CircuitBreaker *circuitBreakerDto `json:"circuitBreaker"`
}

type circuitBreakerDto struct {
	MaxFailures *int `json:"maxFailures" validate:"required"`
	OpenSeconds *int `json:"openSeconds" validate:"required"`
}

type routeDto struct {
	DomainName *string       `json:"domainName" validate:"required"`
	Backends   *[]backendDto `json:"backends" validate:"required"`
}

type streamResponseDto struct {
	ID             *uuid.UUID     `json:"id" validate:"required"`
	Enabled        *bool          `json:"enabled" validate:"required"`
	Name           *string        `json:"name" validate:"required"`
	Type           *string        `json:"type" validate:"required"`
	FeatureSet     *featureSetDto `json:"featureSet" validate:"required"`
	DefaultBackend *backendDto    `json:"defaultBackend" validate:"required"`
	Binding        *addressDto    `json:"binding" validate:"required"`
	Routes         *[]routeDto    `json:"routes"`
}
