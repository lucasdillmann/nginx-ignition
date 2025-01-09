package server

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"log"
	"net"
	"net/http"
)

type startup struct {
	configuration *configuration.Configuration
	state         *state
}

func registerStartup(
	lifecycle *lifecycle.Lifecycle,
	configuration *configuration.Configuration,
	state *state,
) {
	command := &startup{configuration, state}
	lifecycle.RegisterStartup(command)
}

func (s *startup) Run() error {
	port, err := s.configuration.Get("nginx-ignition.server.port")
	if err != nil {
		return err
	}

	log.Printf("Starting HTTP server on port %s", port)
	s.state.server = &http.Server{Handler: s.state.engine.Handler()}

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}

	s.state.listener = &listener
	go func() {
		_ = s.state.server.Serve(listener)
	}()

	return nil
}

func (s *startup) Priority() int {
	return 500
}

func (s *startup) Async() bool {
	return false
}
