package tailscale

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"tailscale.com/client/local"
	"tailscale.com/tsnet"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type tailnetEndpoint struct {
	client      *local.Client
	server      *tsnet.Server
	destination vpn.Destination
	serverURL   string
	authKey     string
	configDir   string
	listeners   []net.Listener
}

func (e *tailnetEndpoint) Stop(ctx context.Context) {
	log.Infof("Stopping Tailscale endpoint %s...", e.destination.SourceName())

	for _, listener := range e.listeners {
		_ = listener.Close()
	}

	_ = e.client.Logout(ctx)
	_ = e.server.Close()
}

func (e *tailnetEndpoint) Start(ctx context.Context) error {
	log.Infof("Starting tailscale %s endpoint...", e.destination.SourceName())

	e.server = new(tsnet.Server)
	e.server.AuthKey = e.authKey
	e.server.ControlURL = e.serverURL
	e.server.Hostname = e.destination.SourceName()
	e.server.Ephemeral = true
	e.server.UserLogf = noOpLogger
	e.server.Logf = noOpLogger
	e.server.Dir = fmt.Sprintf("%s/tsnet/%s", e.configDir, e.destination.SourceName())

	if _, err := e.server.Up(ctx); err != nil {
		return err
	}

	var err error
	if e.client, err = e.server.LocalClient(); err != nil {
		return err
	}

	for _, target := range e.destination.Targets() {
		if err := e.startListener(target); err != nil {
			return err
		}
	}

	ipv4, ipv6 := e.server.TailscaleIPs()
	log.Infof(
		"Tailscale endpoint %s started (IPv4 %v; IPv6 %v)",
		e.destination.SourceName(),
		ipv4,
		ipv6,
	)

	return nil
}

func (e *tailnetEndpoint) startListener(target vpn.DestinationTarget) error {
	proxy := new(httputil.ReverseProxy)
	proxy.ErrorLog = log.Std()
	proxy.Director = func(req *http.Request) {
		ipAddr := target.IP
		if ipAddr == "0.0.0.0" {
			ipAddr = "127.0.0.1"
		}

		req.URL.Host = fmt.Sprintf("%s:%d", ipAddr, target.Port)
		req.URL.Scheme = "http"
		if target.HTTPS {
			req.URL.Scheme = "https"
		}

		req.Header.Del("Host")
		req.Header.Set("Host", target.Host)
		req.Host = target.Host
	}

	startListener := e.server.Listen
	if target.HTTPS {
		startListener = e.server.ListenTLS
	}

	var listener net.Listener
	var err error

	if listener, err = startListener("tcp", fmt.Sprintf(":%d", target.Port)); err != nil {
		return err
	}

	e.listeners = append(e.listeners, listener)

	go func() {
		svr := &http.Server{
			ReadHeaderTimeout: 10 * time.Second,
			Handler:           http.HandlerFunc(proxy.ServeHTTP),
		}
		_ = svr.Serve(listener)
	}()

	return nil
}

func noOpLogger(_ string, _ ...any) {
	// NO-OP
}
