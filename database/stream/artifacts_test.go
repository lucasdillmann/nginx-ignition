package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func newStream() *stream.Stream {
	return &stream.Stream{
		ID:      uuid.New(),
		Name:    "Test Stream",
		Type:    stream.SimpleType,
		Enabled: true,
		Binding: stream.Address{
			Port:     ptr.Of(8080),
			Protocol: stream.TCPProtocol,
			Address:  "0.0.0.0",
		},
		DefaultBackend: stream.Backend{
			Address: stream.Address{
				Port:     ptr.Of(9090),
				Protocol: stream.TCPProtocol,
				Address:  "127.0.0.1",
			},
			Weight: ptr.Of(1),
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
