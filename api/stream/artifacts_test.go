package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
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
			Address:  "0.0.0.0",
			Port:     ptr.Of(80),
			Protocol: stream.TCPProtocol,
		},
		DefaultBackend: stream.Backend{
			Address: stream.Address{
				Address:  "127.0.0.1",
				Port:     ptr.Of(8080),
				Protocol: stream.TCPProtocol,
			},
		},
	}
}

func newStreamRequest() streamRequestDTO {
	return streamRequestDTO{
		Name:    ptr.Of("Test Stream"),
		Type:    ptr.Of(string(stream.SimpleType)),
		Enabled: ptr.Of(true),
		Binding: &addressDTO{
			Address:  ptr.Of("0.0.0.0"),
			Port:     ptr.Of(80),
			Protocol: stream.TCPProtocol,
		},
		DefaultBackend: &backendDTO{
			Target: &addressDTO{
				Address:  ptr.Of("127.0.0.1"),
				Port:     ptr.Of(8080),
				Protocol: stream.TCPProtocol,
			},
		},
	}
}

func newStreamPage() *pagination.Page[stream.Stream] {
	return pagination.Of([]stream.Stream{
		*newStream(),
	})
}
