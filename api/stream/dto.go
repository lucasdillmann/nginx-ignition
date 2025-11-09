package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/stream"
)

type streamRequestDto struct {
	Enabled        *bool          `json:"enabled"`
	Name           *string        `json:"name"`
	Type           *string        `json:"type"`
	FeatureSet     *featureSetDto `json:"featureSet"`
	DefaultBackend *backendDto    `json:"defaultBackend"`
	Binding        *addressDto    `json:"binding"`
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
	Target         *addressDto        `json:"target"`
	CircuitBreaker *circuitBreakerDto `json:"circuitBreaker"`
}

type circuitBreakerDto struct {
	MaxFailures *int `json:"maxFailures"`
	OpenSeconds *int `json:"openSeconds"`
}

type routeDto struct {
	DomainNames *[]string     `json:"domainNames"`
	Backends    *[]backendDto `json:"backends"`
}

type streamResponseDto struct {
	ID             *uuid.UUID     `json:"id"`
	Enabled        *bool          `json:"enabled"`
	Name           *string        `json:"name"`
	Type           *string        `json:"type"`
	FeatureSet     *featureSetDto `json:"featureSet"`
	DefaultBackend *backendDto    `json:"defaultBackend"`
	Binding        *addressDto    `json:"binding"`
	Routes         *[]routeDto    `json:"routes"`
}
