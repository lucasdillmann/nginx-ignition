package stream

import "github.com/google/uuid"

type Protocol string

const (
	UDPProtocol    Protocol = "UDP"
	TCPProtocol    Protocol = "TCP"
	SocketProtocol Protocol = "SOCKET"
)

type Type string

const (
	SimpleType    Type = "SIMPLE"
	SNIRouterType Type = "SNI_ROUTER"
)

type Stream struct {
	DefaultBackend Backend
	Binding        Address
	Name           string
	Type           Type
	Routes         []Route
	ID             uuid.UUID
	FeatureSet     FeatureSet
	Enabled        bool
}

type Route struct {
	DomainNames []string
	Backends    []Backend
}

type Backend struct {
	Weight         *int
	CircuitBreaker *CircuitBreaker
	Address        Address
}

type CircuitBreaker struct {
	MaxFailures int
	OpenSeconds int
}

type Address struct {
	Port     *int
	Protocol Protocol
	Address  string
}

type FeatureSet struct {
	UseProxyProtocol bool
	SocketKeepAlive  bool
	TCPKeepAlive     bool
	TCPNoDelay       bool
	TCPDeferred      bool
}
