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
	destination *vpn.Destination
	name        string
	serverURL   string
	authKey     string
	listener    net.Listener
	configDir   string
}

func (e *tailnetEndpoint) Stop(ctx context.Context) {
	log.Infof("Stopping Tailscale endpoint %s...", e.name)

	_ = e.listener.Close()
	_ = e.client.Logout(ctx)
	_ = e.server.Close()
}

func (e *tailnetEndpoint) Start(ctx context.Context) error {
	log.Infof(
		"Starting tailscale %s endpoint for domain %s...",
		e.name,
		e.destination.DomainName,
	)

	server := new(tsnet.Server)
	server.AuthKey = e.authKey
	server.ControlURL = e.serverURL
	server.Hostname = e.name
	server.Ephemeral = true
	server.UserLogf = noOpLogger
	server.Logf = noOpLogger
	server.Dir = fmt.Sprintf("%s/tsnet/%s", e.configDir, e.name)

	if _, err := server.Up(ctx); err != nil {
		return err
	}

	ipv4, ipv6 := server.TailscaleIPs()
	log.Infof(
		"Tailscale endpoint %s started on hostname %s, IPv4 %v and IPv6 %v. Starting reverse proxy...",
		e.name,
		server.Hostname,
		ipv4,
		ipv6,
	)

	var err error
	if e.client, err = server.LocalClient(); err != nil {
		return err
	}

	scheme := "http"
	if e.destination.HTTPS {
		scheme = "https"
	}

	proxy := new(httputil.ReverseProxy)
	proxy.Director = func(req *http.Request) {
		req.Host = e.destination.DomainName
		req.URL.Host = e.destination.DomainName
		req.URL.Scheme = scheme
		req.Header.Del("Host")
		req.Header.Set("Host", req.URL.Host)
	}

	port := fmt.Sprintf(":%d", e.destination.Port)
	if e.listener, err = server.Listen("tcp", port); err != nil {
		return err
	}

	if e.destination.HTTPS {
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

	log.Infof("Tailscale endpoint %s started successfully", e.name)

	return nil
}

func noOpLogger(_ string, _ ...any) {
	// NO-OP
}
