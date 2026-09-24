package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/lifecycle"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/log"
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
	serverCfg, err := loadServerConfig(s.configuration)
	if err != nil {
		return err
	}

	log.Infof("Starting HTTP server on port %s", serverCfg.port)
	s.state.server = &http.Server{
		Handler:           s.state.engine.Handler(),
		IdleTimeout:       serverCfg.idleTimeout,
		WriteTimeout:      serverCfg.writeTimeout,
		ReadTimeout:       serverCfg.readTimeout,
		ReadHeaderTimeout: serverCfg.readHeaderTimeout,
		MaxHeaderBytes:    serverCfg.maxHeaderBytes,
		ErrorLog:          log.Std(),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", serverCfg.address, serverCfg.port))
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
