package server

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type shutdown struct {
	state *state
}

func registerShutdown(lifecycle *lifecycle.Lifecycle, state *state) {
	lifecycle.RegisterShutdown(shutdown{state})
}

func (s shutdown) Run(_ context.Context) {
	log.Infof("Stopping the HTTP server")

	if err := s.state.server.Close(); err != nil {
		log.Warnf("Failed to stop HTTP server: %v", err)
	}
}

func (s shutdown) Priority() int {
	return shutdownPriority
}
