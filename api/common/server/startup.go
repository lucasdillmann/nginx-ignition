package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
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
	lifecycle.RegisterStartup(startup{configuration, state})
}

func (s startup) Run(_ context.Context) error {
	port, err := s.configuration.Get("nginx-ignition.server.binding-port")
	if err != nil {
		return err
	}

	address, err := s.configuration.Get("nginx-ignition.server.binding-address")
	if err != nil {
		return err
	}

	log.Infof("Starting HTTP server on port %s", port)
	s.state.server = &http.Server{Handler: s.state.engine.Handler()}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return err
	}

	s.state.listener = &listener
	go func() {
		_ = s.state.server.Serve(listener)
	}()

	return nil
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return false
}
