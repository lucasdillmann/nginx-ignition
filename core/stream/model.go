package stream

import "github.com/google/uuid"

type Protocol string

const (
	UDPProtocol    Protocol = "UDP"
	TCPProtocol    Protocol = "TCP"
	SocketProtocol Protocol = "SOCKET"
)

type Stream struct {
	ID         uuid.UUID
	Enabled    bool
	Name       string
	Binding    Address
	Backend    Address
	FeatureSet FeatureSet
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
