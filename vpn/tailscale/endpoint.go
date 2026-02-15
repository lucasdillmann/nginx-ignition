package tailscale

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"time"

	"tailscale.com/client/local"
	"tailscale.com/tsnet"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type tailnetEndpoint struct {
	client    *local.Client
	server    *tsnet.Server
	endpoint  vpn.Endpoint
	serverURL string
	authKey   string
	configDir string
	listeners []net.Listener
}

func (e *tailnetEndpoint) stop(ctx context.Context) {
	log.Infof("Stopping Tailscale endpoint %s...", e.endpoint.SourceName())

	for _, listener := range e.listeners {
		_ = listener.Close()
	}

	_ = e.client.Logout(ctx)
	_ = e.server.Close()
}

func (e *tailnetEndpoint) start(ctx context.Context) error {
	log.Infof("Starting tailscale %s endpoint...", e.endpoint.SourceName())

	e.server = new(tsnet.Server)
	e.server.AuthKey = e.authKey
	e.server.ControlURL = e.serverURL
	e.server.Hostname = e.endpoint.SourceName()
	e.server.Ephemeral = true
	e.server.UserLogf = noOpLogger
	e.server.Logf = noOpLogger
	e.server.Dir = filepath.Join(e.configDir, "tsnet", e.endpoint.SourceName())

	if _, err := e.server.Up(ctx); err != nil {
		return err
	}

	var err error
	if e.client, err = e.server.LocalClient(); err != nil {
		return err
	}

	for _, target := range e.endpoint.Targets() {
		if err := e.startListener(target); err != nil {
			return err
		}
	}

	ipv4, ipv6 := e.server.TailscaleIPs()
	log.Infof(
		"Tailscale endpoint %s started (IPv4 %v; IPv6 %v)",
		e.endpoint.SourceName(),
		ipv4,
		ipv6,
	)

	return nil
}

func (e *tailnetEndpoint) startListener(target vpn.EndpointTarget) error {
	proxy := new(httputil.ReverseProxy)
	proxy.ErrorLog = log.Std()

	scheme := "http"
	if target.HTTPS {
		scheme = "https"
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: target.Host,
				MinVersion: tls.VersionTLS12,
			},
		}
	}

	proxy.Rewrite = func(pr *httputil.ProxyRequest) {
		ipAddr := target.IP
		if ipAddr == "0.0.0.0" {
			ipAddr = "127.0.0.1"
		}

		pr.Out.URL.Host = fmt.Sprintf("%s:%d", ipAddr, target.Port)
		pr.Out.URL.Scheme = scheme
		pr.Out.Header.Del("Host")
		pr.Out.Header.Set("Host", target.Host)
		pr.Out.Host = target.Host
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
