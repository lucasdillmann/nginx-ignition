package stream

import "github.com/google/uuid"

type Protocol string

const (
	UDPProtocol    Protocol = "UDP"
	TCPProtocol    Protocol = "TCP"
	SocketProtocol Protocol = "SOCKET"
)

type Stream struct {
	ID          uuid.UUID
	Enabled     bool
	Description string
	Binding     Address
	Backend     Address
	FeatureSet  FeatureSet
}

type Address struct {
	Protocol Protocol
	Address  string
	Port     *int
}

type FeatureSet struct {
	UseProxyProtocol bool
	SSL              bool
	TCPKeepAlive     bool
	TCPNoDelay       bool
	TCPDeferred      bool
}
