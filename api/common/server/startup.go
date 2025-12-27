package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

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
	port, err := s.configuration.Get("nginx-ignition.server.port")
	if err != nil {
		return err
	}

	address, err := s.configuration.Get("nginx-ignition.server.address")
	if err != nil {
		return err
	}

	log.Infof("Starting HTTP server on port %s", port)
	s.state.server = &http.Server{
		Handler:           s.state.engine.Handler(),
		IdleTimeout:       120 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       15 * time.Second,
		ErrorLog:          log.Std(),
	}

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
