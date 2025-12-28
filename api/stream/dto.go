package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/stream"
)

type streamRequestDTO struct {
	Enabled        *bool          `json:"enabled"`
	Name           *string        `json:"name"`
	Type           *string        `json:"type"`
	FeatureSet     *featureSetDTO `json:"featureSet"`
	DefaultBackend *backendDTO    `json:"defaultBackend"`
	Binding        *addressDTO    `json:"binding"`
	Routes         []routeDTO     `json:"routes"`
}

type featureSetDTO struct {
	UseProxyProtocol *bool `json:"useProxyProtocol"`
	SocketKeepAlive  *bool `json:"socketKeepAlive"`
	TCPKeepAlive     *bool `json:"tcpKeepAlive"`
	TCPNoDelay       *bool `json:"tcpNoDelay"`
	TCPDeferred      *bool `json:"tcpDeferred"`
}

type addressDTO struct {
	Address  *string         `json:"address"`
	Port     *int            `json:"port"`
	Protocol stream.Protocol `json:"protocol"`
}

type backendDTO struct {
	Weight         *int               `json:"weight"`
	Target         *addressDTO        `json:"target"`
	CircuitBreaker *circuitBreakerDTO `json:"circuitBreaker"`
}

type circuitBreakerDTO struct {
	MaxFailures *int `json:"maxFailures"`
	OpenSeconds *int `json:"openSeconds"`
}

type routeDTO struct {
	DomainNames []string     `json:"domainNames"`
	Backends    []backendDTO `json:"backends"`
}

type streamResponseDTO struct {
	ID             *uuid.UUID     `json:"id"`
	Enabled        *bool          `json:"enabled"`
	Name           *string        `json:"name"`
	Type           *string        `json:"type"`
	FeatureSet     *featureSetDTO `json:"featureSet"`
	DefaultBackend *backendDTO    `json:"defaultBackend"`
	Binding        *addressDTO    `json:"binding"`
	Routes         []routeDTO     `json:"routes"`
}
