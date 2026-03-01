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
	tsServer  *tsnet.Server
	endpoint  vpn.Endpoint
	serverURL string
	authKey   string
	configDir string
	listeners []net.Listener
	hServers  []*http.Server
}

func (e *tailnetEndpoint) stop(ctx context.Context) {
	log.Infof("Stopping Tailscale endpoint %s...", e.endpoint.SourceName())

	for _, server := range e.hServers {
		_ = server.Shutdown(ctx)
	}

	for _, listener := range e.listeners {
		_ = listener.Close()
	}

	_ = e.client.Logout(ctx)
	_ = e.tsServer.Close()
}

func (e *tailnetEndpoint) start(ctx context.Context) error {
	log.Infof("Starting tailscale %s endpoint...", e.endpoint.SourceName())

	e.tsServer = new(tsnet.Server)
	e.tsServer.AuthKey = e.authKey
	e.tsServer.ControlURL = e.serverURL
	e.tsServer.Hostname = e.endpoint.SourceName()
	e.tsServer.Ephemeral = true
	e.tsServer.UserLogf = noOpLogger
	e.tsServer.Logf = noOpLogger
	e.tsServer.Dir = filepath.Join(e.configDir, "tsnet", e.endpoint.SourceName())

	if _, err := e.tsServer.Up(ctx); err != nil {
		return err
	}

	var err error
	if e.client, err = e.tsServer.LocalClient(); err != nil {
		return err
	}

	for _, target := range e.endpoint.Targets() {
		if err := e.startListener(target); err != nil {
			e.stop(ctx)
			return err
		}
	}

	ipv4, ipv6 := e.tsServer.TailscaleIPs()
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
	if target.HTTPS.Enabled {
		scheme = "https"

		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{
			ServerName: target.Host,
			MinVersion: tls.VersionTLS12,
		}

		proxy.Transport = transport
	}

	proxy.Rewrite = func(pr *httputil.ProxyRequest) {
		ipAddr := target.IP
		if ipAddr == "0.0.0.0" {
			ipAddr = "127.0.0.1"
		}

		pr.SetXForwarded()
		pr.Out.URL.Host = fmt.Sprintf("%s:%d", ipAddr, target.Port)
		pr.Out.URL.Scheme = scheme
		pr.Out.Header.Del("Host")
		pr.Out.Header.Set("Host", target.Host)
		pr.Out.Host = target.Host
	}

	startListener := e.tsServer.Listen
	if target.HTTPS.Enabled {
		startListener = e.tsServer.ListenTLS
	}

	var listener net.Listener
	var err error

	if listener, err = startListener("tcp", fmt.Sprintf(":%d", target.Port)); err != nil {
		return err
	}

	svr := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           http.HandlerFunc(proxy.ServeHTTP),
	}

	e.listeners = append(e.listeners, listener)
	e.hServers = append(e.hServers, svr)

	go func() {
		_ = svr.Serve(listener)
	}()

	return nil
}

func noOpLogger(_ string, _ ...any) {
	// NO-OP
}
