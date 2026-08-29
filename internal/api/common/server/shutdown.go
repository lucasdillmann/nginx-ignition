package server

import (
	"context"

	"nginx-ignition/internal/core/common/lifecycle"
	"nginx-ignition/internal/core/common/log"
)

type shutdown struct {
	state *state
}

func registerShutdown(lc *lifecycle.Lifecycle, state *state) {
	lc.RegisterShutdown(shutdown{state})
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
