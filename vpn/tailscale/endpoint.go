package tailscale

import (
	"context"
	"crypto/tls"
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
	listener    net.Listener
	configDir   string
}

func (e *tailnetEndpoint) Stop(ctx context.Context) {
	log.Infof("Stopping Tailscale endpoint %s...", e.destination.Name())

	_ = e.listener.Close()
	_ = e.client.Logout(ctx)
	_ = e.server.Close()
}

func (e *tailnetEndpoint) Start(ctx context.Context) error {
	log.Infof(
		"Starting tailscale %s endpoint for domain %s...",
		e.destination.Name(),
		e.destination.DomainName(),
	)

	e.server = new(tsnet.Server)
	e.server.AuthKey = e.authKey
	e.server.ControlURL = e.serverURL
	e.server.Hostname = e.destination.Name()
	e.server.Ephemeral = true
	e.server.UserLogf = noOpLogger
	e.server.Logf = noOpLogger
	e.server.Dir = fmt.Sprintf("%s/tsnet/%s", e.configDir, e.destination.Name())

	if _, err := e.server.Up(ctx); err != nil {
		return err
	}

	var err error
	if e.client, err = e.server.LocalClient(); err != nil {
		return err
	}

	scheme := "http"
	if e.destination.HTTPS() {
		scheme = "https"
	}

	proxy := new(httputil.ReverseProxy)
	proxy.ErrorLog = log.Std()
	proxy.Director = func(req *http.Request) {
		ipAddr := e.destination.IP()
		if ipAddr == "0.0.0.0" {
			ipAddr = "127.0.0.1"
		}

		req.URL.Host = fmt.Sprintf("%s:%d", ipAddr, e.destination.Port())
		req.URL.Scheme = scheme

		req.Header.Del("Host")
		req.Header.Set("Host", e.destination.DomainName())
		req.Host = e.destination.DomainName()
	}

	port := fmt.Sprintf(":%d", e.destination.Port())
	if e.listener, err = e.server.Listen("tcp", port); err != nil {
		return err
	}

	if e.destination.HTTPS() {
		e.listener = tls.NewListener(
			e.listener,
			&tls.Config{
				MinVersion:     tls.VersionTLS12,
				GetCertificate: e.client.GetCertificate,
			},
		)
	}

	go func() {
		svr := &http.Server{
			ReadHeaderTimeout: 10 * time.Second,
			Handler:           http.HandlerFunc(proxy.ServeHTTP),
		}
		_ = svr.Serve(e.listener)
	}()

	ipv4, ipv6 := e.server.TailscaleIPs()
	log.Infof(
		"Tailscale endpoint %s started on hostname %s, IPv4 %v and IPv6 %v",
		e.destination.Name(),
		e.server.Hostname,
		ipv4,
		ipv6,
	)

	return nil
}

func noOpLogger(_ string, _ ...any) {
	// NO-OP
}
