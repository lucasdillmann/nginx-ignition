package netbird

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"strings"
	"time"

	netbird "github.com/netbirdio/netbird/client/embed"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type netbirdEndpoint struct {
	client        *netbird.Client
	endpoint      vpn.Endpoint
	managementURL string
	setupKey      string
	configDir     string
	listeners     []net.Listener
	servers       []*http.Server
}

func (e *netbirdEndpoint) stop(ctx context.Context) {
	log.Infof("Stopping NetBird endpoint %s...", e.endpoint.SourceName())

	for _, server := range e.servers {
		_ = server.Shutdown(ctx)
	}

	for _, listener := range e.listeners {
		_ = listener.Close()
	}

	_ = e.client.Stop(ctx)
}

func (e *netbirdEndpoint) start(ctx context.Context) error {
	log.Infof("Starting NetBird %s endpoint...", e.endpoint.SourceName())
	basePath := filepath.Join(e.configDir, "netbird", e.endpoint.SourceName())

	opts := netbird.Options{
		DeviceName:    e.endpoint.SourceName(),
		SetupKey:      e.setupKey,
		ManagementURL: e.managementURL,
		StatePath:     filepath.Join(basePath, "state.json"),
		ConfigPath:    filepath.Join(basePath, "config.json"),
		LogOutput:     &noOpLogger{},
	}

	var err error
	if e.client, err = netbird.New(opts); err != nil {
		return err
	}

	if err = e.client.Start(ctx); err != nil {
		return err
	}

	for _, target := range e.endpoint.Targets() {
		if err = e.startListener(target); err != nil {
			e.stop(ctx)
			return err
		}
	}

	status, err := e.client.Status()
	if err == nil {
		log.Infof(
			"NetBird endpoint %s started (IP %v)",
			e.endpoint.SourceName(),
			status.LocalPeerState.IP,
		)
	} else {
		log.Infof("NetBird endpoint %s started", e.endpoint.SourceName())
	}

	return nil
}

func (e *netbirdEndpoint) startListener(target vpn.EndpointTarget) error {
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

	listener, err := e.client.ListenTCP(fmt.Sprintf(":%d", target.Port))
	if err != nil {
		return err
	}

	if target.HTTPS.Enabled {
		tlsCerts, err := buildTLSCertificate(target.HTTPS)
		if err != nil {
			_ = listener.Close()
			return err
		}

		listener = tls.NewListener(listener, &tls.Config{
			Certificates: tlsCerts,
			MinVersion:   tls.VersionTLS12,
		})
	}

	svr := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           http.HandlerFunc(proxy.ServeHTTP),
	}

	e.listeners = append(e.listeners, listener)
	e.servers = append(e.servers, svr)

	go func() {
		_ = svr.Serve(listener)
	}()

	return nil
}

func buildTLSCertificate(cert vpn.EndpointHTTPS) ([]tls.Certificate, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(cert.PublicKey)
	if err != nil {
		return nil, err
	}

	fullChainPem := convertToPemEncodedCertificateString(publicKeyBytes)
	for _, chain := range cert.CertificationChain {
		//nolint:govet
		decodedChain, err := base64.StdEncoding.DecodeString(chain)
		if err != nil {
			return nil, err
		}

		fullChainPem += "\n" + convertToPemEncodedCertificateString(decodedChain)
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(cert.PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKeyPem := convertToPemEncodedPrivateKeyString(privateKeyBytes)
	keyPair, err := tls.X509KeyPair([]byte(fullChainPem), []byte(privateKeyPem))
	if err != nil {
		return nil, err
	}

	return []tls.Certificate{keyPair}, nil
}

func convertToPemEncodedCertificateString(bytes []byte) string {
	if strings.Contains(string(bytes), "CERTIFICATE") {
		return string(bytes)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: bytes,
	})
	return string(certPEM)
}

func convertToPemEncodedPrivateKeyString(bytes []byte) string {
	if strings.Contains(string(bytes), "PRIVATE KEY") {
		return string(bytes)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bytes,
	})
	return string(keyPEM)
}
