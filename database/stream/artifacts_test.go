package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/stream"
)

func newStream() *stream.Stream {
	return &stream.Stream{
		ID:      uuid.New(),
		Name:    "Test Stream",
		Type:    stream.SimpleType,
		Enabled: true,
		Binding: stream.Address{
			Port:     new(8080),
			Protocol: stream.TCPProtocol,
			Address:  "0.0.0.0",
		},
		DefaultBackend: stream.Backend{
			Address: stream.Address{
				Port:     new(9090),
				Protocol: stream.TCPProtocol,
				Address:  "127.0.0.1",
			},
			Weight: new(1),
			CircuitBreaker: &stream.CircuitBreaker{
				MaxFailures: 3,
				OpenSeconds: 10,
			},
		},
		FeatureSet: stream.FeatureSet{
			UseProxyProtocol: true,
			SocketKeepAlive:  true,
			TCPKeepAlive:     true,
			TCPNoDelay:       true,
			TCPDeferred:      false,
		},
	}
}
