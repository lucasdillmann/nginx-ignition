package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"nginx-ignition/internal/core/common/configuration"
	"nginx-ignition/internal/core/common/lifecycle"
	"nginx-ignition/internal/core/common/log"
)

type startup struct {
	configuration *configuration.Configuration
	state         *state
}

func registerStartup(
	lc *lifecycle.Lifecycle,
	cfg *configuration.Configuration,
	state *state,
) {
	lc.RegisterStartup(startup{cfg, state})
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
