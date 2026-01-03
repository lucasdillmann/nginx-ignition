package cfgfiles

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_StreamFileProvider_Provide(t *testing.T) {
	p := &streamFileProvider{}
	id := uuid.New()
	ctx := &providerContext{
		context: context.Background(),
		supportedFeatures: &SupportedFeatures{
			StreamType: StaticSupportType,
		},
		streams: []stream.Stream{
			{
				ID: id,
				Binding: stream.Address{
					Protocol: stream.TCPProtocol,
					Address:  "0.0.0.0",
					Port:     ptr.Of(80),
				},
				Type: stream.SimpleType,
				DefaultBackend: stream.Backend{
					Address: stream.Address{
						Protocol: stream.TCPProtocol,
						Address:  "127.0.0.1",
						Port:     ptr.Of(8080),
					},
				},
			},
		},
	}

	files, err := p.provide(ctx)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, fmt.Sprintf("stream-%s.conf", id), files[0].Name)

	t.Run("returns error when streams present but not supported", func(t *testing.T) {
		ctx.supportedFeatures.StreamType = NoneSupportType
		_, err := p.provide(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Support for streams is not enabled")
	})

	t.Run("returns error for unknown stream type", func(t *testing.T) {
		ctx.supportedFeatures.StreamType = StaticSupportType
		ctx.streams[0].Type = "UNKNOWN"
		_, err := p.provide(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown stream type")
	})
}

func Test_StreamFileProvider_BuildBinding(t *testing.T) {
	p := &streamFileProvider{}

	t.Run("TCP binding with all flags", func(t *testing.T) {
		s := &stream.Stream{
			Binding: stream.Address{
				Protocol: stream.TCPProtocol,
				Address:  "0.0.0.0",
				Port:     ptr.Of(80),
			},
			FeatureSet: stream.FeatureSet{
				UseProxyProtocol: true,
				TCPDeferred:      true,
				TCPKeepAlive:     true,
				TCPNoDelay:       true,
			},
		}

		result, err := p.buildBinding(s)
		assert.NoError(t, err)
		assert.Equal(
			t,
			"listen 0.0.0.0:80 proxy_protocol deferred so_keepalive=on reuseport;",
			*result,
		)
	})

	t.Run("UDP binding", func(t *testing.T) {
		s := &stream.Stream{
			Binding: stream.Address{
				Protocol: stream.UDPProtocol,
				Address:  "127.0.0.1",
				Port:     ptr.Of(53),
			},
		}

		result, err := p.buildBinding(s)
		assert.NoError(t, err)
		assert.Equal(t, "listen 127.0.0.1:53 udp reuseport;", *result)
	})

	t.Run("Unix socket binding", func(t *testing.T) {
		s := &stream.Stream{
			Binding: stream.Address{
				Protocol: stream.SocketProtocol,
				Address:  "/tmp/nginx.sock",
			},
		}

		result, err := p.buildBinding(s)
		assert.NoError(t, err)
		assert.Equal(t, "listen unix:/tmp/nginx.sock reuseport;", *result)
	})

	t.Run("returns error for unknown protocol", func(t *testing.T) {
		s := &stream.Stream{
			Binding: stream.Address{
				Protocol: "GOPHER",
			},
		}
		_, err := p.buildBinding(s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown binding protocol")
	})
}

func Test_StreamFileProvider_BuildUpstream(t *testing.T) {
	p := &streamFileProvider{}

	t.Run("generates upstream with circuit breaker and weight", func(t *testing.T) {
		backends := []stream.Backend{
			{
				Address: stream.Address{
					Protocol: stream.TCPProtocol,
					Address:  "10.0.0.1",
					Port:     ptr.Of(8080),
				},
				Weight: ptr.Of(5),
				CircuitBreaker: &stream.CircuitBreaker{
					MaxFailures: 3,
					OpenSeconds: 30,
				},
			},
			{
				Address: stream.Address{
					Protocol: stream.SocketProtocol,
					Address:  "/var/run/backend.sock",
				},
			},
		}

		result, err := p.buildUpstream(backends, "test_upstream")
		assert.NoError(t, err)
		assert.Contains(t, *result, "upstream test_upstream {")
		assert.Contains(t, *result, "server 10.0.0.1:8080 weight=5 max_fails=3 fail_timeout=30s;")
		assert.Contains(t, *result, "server unix:/var/run/backend.sock;")
	})

	t.Run("returns error for unknown backend protocol", func(t *testing.T) {
		backends := []stream.Backend{
			{
				Address: stream.Address{Protocol: "GOPHER"},
			},
		}
		_, err := p.buildUpstream(backends, "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown backend protocol")
	})
}

func Test_StreamFileProvider_BuildRoutedStream(t *testing.T) {
	p := &streamFileProvider{}
	id := uuid.New()
	idStr := nginxID(&stream.Stream{ID: id})

	t.Run("generates SNI routing configuration", func(t *testing.T) {
		ctx := &providerContext{
			supportedFeatures: &SupportedFeatures{
				TLSSNI: StaticSupportType,
			},
		}
		s := &stream.Stream{
			ID:   id,
			Type: stream.SNIRouterType,
			Binding: stream.Address{
				Protocol: stream.TCPProtocol,
				Address:  "0.0.0.0",
				Port:     ptr.Of(443),
			},
			DefaultBackend: stream.Backend{
				Address: stream.Address{
					Protocol: stream.TCPProtocol,
					Address:  "127.0.0.1",
					Port:     ptr.Of(8443),
				},
			},
			Routes: []stream.Route{
				{
					DomainNames: []string{"example.com"},
					Backends: []stream.Backend{
						{
							Address: stream.Address{
								Protocol: stream.TCPProtocol,
								Address:  "10.0.0.1",
								Port:     ptr.Of(443),
							},
						},
					},
				},
			},
		}

		result, err := p.buildRoutedStream(ctx, s)
		assert.NoError(t, err)
		assert.Contains(
			t,
			*result,
			fmt.Sprintf("map $ssl_preread_server_name $stream_%s_router {", idStr),
		)
		assert.Contains(t, *result, fmt.Sprintf("example.com stream_%s_route_0;", idStr))
		assert.Contains(t, *result, fmt.Sprintf("default stream_%s_default;", idStr))
		assert.Contains(t, *result, "ssl_preread on;")
		assert.Contains(t, *result, fmt.Sprintf("proxy_pass $stream_%s_router;", idStr))
	})

	t.Run("returns error when TLSSNI not supported", func(t *testing.T) {
		ctx := &providerContext{
			supportedFeatures: &SupportedFeatures{
				TLSSNI: NoneSupportType,
			},
		}
		s := &stream.Stream{Type: stream.SNIRouterType}
		_, err := p.buildRoutedStream(ctx, s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Support for TLS SNI is not enabled")
	})
}
