package stream

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func newStream() *Stream {
	return &Stream{
		ID:   uuid.New(),
		Name: "test",
		Type: SimpleType,
		Binding: Address{
			Protocol: TCPProtocol,
			Address:  "127.0.0.1",
			Port:     ptr.Of(8080),
		},
		DefaultBackend: Backend{
			Address: Address{
				Protocol: TCPProtocol,
				Address:  "127.0.0.1",
				Port:     ptr.Of(8081),
			},
		},
	}
}
