package server

import (
	"time"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
)

type serverConfig struct {
	port              string
	address           string
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	readHeaderTimeout time.Duration
	maxHeaderBytes    int
}

func loadServerConfig(cfg *configuration.Configuration) (*serverConfig, error) {
	serverCfg := cfg.WithPrefix("nginx-ignition.server")

	port, err := serverCfg.Get("port")
	if err != nil {
		return nil, err
	}

	address, err := serverCfg.Get("address")
	if err != nil {
		return nil, err
	}

	readTimeoutSec, err := serverCfg.GetInt("read-timeout-seconds")
	if err != nil {
		return nil, err
	}

	writeTimeoutSec, err := serverCfg.GetInt("write-timeout-seconds")
	if err != nil {
		return nil, err
	}

	idleTimeoutSec, err := serverCfg.GetInt("idle-timeout-seconds")
	if err != nil {
		return nil, err
	}

	readHeaderTimeoutSec, err := serverCfg.GetInt("read-header-timeout-seconds")
	if err != nil {
		return nil, err
	}

	maxHeaderBytes, err := serverCfg.GetInt("max-header-bytes")
	if err != nil {
		return nil, err
	}

	return &serverConfig{
		port:              port,
		address:           address,
		maxHeaderBytes:    maxHeaderBytes,
		readTimeout:       time.Duration(readTimeoutSec) * time.Second,
		writeTimeout:      time.Duration(writeTimeoutSec) * time.Second,
		idleTimeout:       time.Duration(idleTimeoutSec) * time.Second,
		readHeaderTimeout: time.Duration(readHeaderTimeoutSec) * time.Second,
	}, nil
}
