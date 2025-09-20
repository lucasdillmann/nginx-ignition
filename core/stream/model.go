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
	ID             uuid.UUID
	Enabled        bool
	Name           string
	Type           Type
	Routes         []Route
	Binding        Address
	DefaultBackend Backend
	FeatureSet     FeatureSet
}

type Route struct {
	DomainName string
	Backends   []Backend
}

type Backend struct {
	Weight         *int
	Address        Address
	CircuitBreaker *CircuitBreaker
}

type CircuitBreaker struct {
	MaxFailures int
	OpenSeconds int
}

type Address struct {
	Protocol Protocol
	Address  string
	Port     *int
}

type FeatureSet struct {
	UseProxyProtocol bool
	SocketKeepAlive  bool
	TCPKeepAlive     bool
	TCPNoDelay       bool
	TCPDeferred      bool
}
