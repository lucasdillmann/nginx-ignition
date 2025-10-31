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
	"dillmann.com.br/nginx-ignition/core/host"
)

type Endpoint struct {
	client      *local.Client
	server      *tsnet.Server
	destination *host.Host
	bindings    []*host.Binding
	name        string
	configDir   string
	authKey     string
	listeners   []net.Listener
	started     bool
}

func NewEndpoint(
	name, authKey, configDir string,
	destination *host.Host,
	globalBindings []*host.Binding,
) *Endpoint {
	bindings := globalBindings
	if len(destination.Bindings) > 0 {
		bindings = destination.Bindings
	}

	return &Endpoint{
		name:        name,
		authKey:     authKey,
		configDir:   configDir,
		started:     false,
		destination: destination,
		bindings:    bindings,
		listeners:   make([]net.Listener, len(bindings)),
	}
}

func (e *Endpoint) Stop(ctx context.Context) {
	log.Infof("Stopping Tailscale endpoint %s...", e.name)

	if !e.started {
		return
	}

	e.started = false

	for _, listener := range e.listeners {
		_ = listener.Close()
	}

	_ = e.client.Logout(ctx)
	_ = e.server.Close()
}

func (e *Endpoint) Start(ctx context.Context) error {
	domainName := *e.destination.DomainNames[0]
	log.Infof(
		"Starting tailscale %s endpoint forwarding requests to %s...",
		e.name,
		domainName,
	)

	server := new(tsnet.Server)
	server.AuthKey = e.authKey
	server.Hostname = e.name
	server.Ephemeral = true
	server.UserLogf = noOpLogger
	server.Logf = noOpLogger
	server.Dir = fmt.Sprintf("%s/tsnet/%s", e.configDir, e.destination.ID)

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

	proxy := new(httputil.ReverseProxy)
	proxy.Director = func(req *http.Request) {
		req.Host = domainName
		req.URL.Host = domainName
		req.URL.Scheme = domainName
		req.Header.Del("Host")
		req.Header.Set("Host", req.URL.Host)
	}

	for index, binding := range e.bindings {
		port := fmt.Sprintf(":%d", binding.Port)
		if e.listeners[index], err = server.Listen("tcp", port); err != nil {
			return err
		}

		if binding.Type == host.HttpsBindingType {
			e.listeners[index] = tls.NewListener(
				e.listeners[index],
				&tls.Config{
					MinVersion:     tls.VersionTLS12,
					GetCertificate: e.client.GetCertificate,
				},
			)
		}

		go func() {
			svr := &http.Server{
				Handler:      http.HandlerFunc(proxy.ServeHTTP),
				ReadTimeout:  60 * time.Second,
				WriteTimeout: 60 * time.Second,
				IdleTimeout:  60 * time.Second,
			}
			_ = svr.Serve(e.listeners[index])
		}()
	}

	e.started = true
	log.Infof("Tailscale endpoint %s started successfully", e.name)

	return nil
}

func noOpLogger(_ string, _ ...any) {
	// NO-OP
}
