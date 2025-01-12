package server

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type shutdown struct {
	state *state
}

func registerShutdown(lifecycle *lifecycle.Lifecycle, state *state) {
	command := &shutdown{state}
	lifecycle.RegisterShutdown(command)
}

func (s *shutdown) Run() {
	log.Info("Stopping the HTTP server")

	if err := s.state.server.Close(); err != nil {
		log.Warn("Failed to stop HTTP server: %v", err)
	}
}

func (s *shutdown) Priority() int {
	return 500
}
