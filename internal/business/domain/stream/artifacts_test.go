package stream

import (
	"github.com/google/uuid"
)

func newStream() *Stream {
	return &Stream{
		ID:   uuid.New(),
		Name: "test",
		Type: SimpleType,
		Binding: Address{
			Protocol: TCPProtocol,
			Address:  "127.0.0.1",
			Port:     new(8080),
		},
		DefaultBackend: Backend{
			Address: Address{
				Protocol: TCPProtocol,
				Address:  "127.0.0.1",
				Port:     new(8081),
			},
		},
	}
}
