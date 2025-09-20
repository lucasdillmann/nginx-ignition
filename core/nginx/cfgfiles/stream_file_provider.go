package cfgfiles

import (
	"fmt"
	"strings"

	"github.com/aws/smithy-go/ptr"

	"dillmann.com.br/nginx-ignition/core/stream"
)

type streamFileProvider struct{}

func newStreamFileProvider() *streamFileProvider {
	return &streamFileProvider{}
}

func (p *streamFileProvider) provide(ctx *providerContext) ([]File, error) {
	files := make([]File, 0, len(ctx.streams))

	for _, s := range ctx.streams {
		contents, err := p.buildConfigFileContents(s)
		if err != nil {
			return nil, err
		}

		files = append(files, File{
			Name:     fmt.Sprintf("stream-%s.conf", s.ID),
			Contents: *contents,
		})
	}

	return files, nil
}

func (p *streamFileProvider) buildConfigFileContents(s *stream.Stream) (*string, error) {
	switch s.Type {
	case stream.SimpleType:
		return p.buildSimpleStream(s)
	case stream.SNIRouterType:
		return p.buildRoutedStream(s)
	default:
		return nil, fmt.Errorf("unknown stream type: %s", s.Type)
	}
}

func (p *streamFileProvider) buildSimpleStream(s *stream.Stream) (*string, error) {
	upstreamId := fmt.Sprintf("stream_%s_default", s.ID)
	upstream, err := p.buildUpstream([]stream.Backend{s.DefaultBackend}, upstreamId)
	if err != nil {
		return nil, err
	}

	return p.buildStream(s, *upstream, fmt.Sprintf("proxy_pass %s;", upstreamId))
}

func (p *streamFileProvider) buildBinding(s *stream.Stream) (*string, error) {
	instruction := strings.Builder{}
	instruction.WriteString("listen ")

	switch s.Binding.Protocol {
	case stream.SocketProtocol:
		instruction.WriteString(fmt.Sprintf("unix:%s", s.Binding.Address))

	case stream.TCPProtocol:
		instruction.WriteString(fmt.Sprintf("%s:%d", s.Binding.Address, *s.Binding.Port))

		if s.FeatureSet.UseProxyProtocol {
			instruction.WriteString(" proxy_protocol")
		}

		if s.FeatureSet.TCPDeferred {
			instruction.WriteString(" deferred")
		}

		if s.FeatureSet.TCPKeepAlive {
			instruction.WriteString(" so_keepalive=on")
		}

	case stream.UDPProtocol:
		instruction.WriteString(fmt.Sprintf("%s:%d udp", s.Binding.Address, *s.Binding.Port))

	default:
		return nil, fmt.Errorf("unknown binding protocol: %s", s.Binding.Protocol)
	}

	instruction.WriteString(" reuseport;")

	return ptr.String(instruction.String()), nil
}

func (p *streamFileProvider) buildUpstream(backends []stream.Backend, name string) (*string, error) {
	instructions := strings.Builder{}
	instructions.WriteString(fmt.Sprintf("upstream %s {\n", name))

	for _, backend := range backends {
		address := backend.Address
		switch address.Protocol {
		case stream.SocketProtocol:
			instructions.WriteString(fmt.Sprintf("server unix:%s", address.Address))

		case stream.TCPProtocol, stream.UDPProtocol:
			instructions.WriteString(fmt.Sprintf("server %s:%d", address.Address, *address.Port))

		default:
			return nil, fmt.Errorf("unknown backend protocol: %s", address.Protocol)
		}

		if backend.Weight != nil {
			instructions.WriteString(fmt.Sprintf(" weight=%d", *backend.Weight))
		}

		if backend.CircuitBreaker != nil {
			instructions.WriteString(fmt.Sprintf(
				" max_fails=%d fail_timeout=%ds",
				backend.CircuitBreaker.MaxFailures,
				backend.CircuitBreaker.OpenSeconds,
			))
		}

		instructions.WriteString(";\n")
	}

	instructions.WriteString("}\n")
	return ptr.String(instructions.String()), nil
}

func (p *streamFileProvider) buildRoutedStream(s *stream.Stream) (*string, error) {
	mapping := strings.Builder{}
	mappingId := fmt.Sprintf("$stream_%s_router", s.ID)
	mapping.WriteString(fmt.Sprintf("map $ssl_preread_server_name %s {\n", mappingId))

	upstreams := strings.Builder{}
	for routeIndex, route := range s.Routes {
		routeId := fmt.Sprintf("stream_%s_route_%d", s.ID, routeIndex)
		upstream, err := p.buildUpstream(route.Backends, routeId)
		if err != nil {
			return nil, err
		}

		upstreams.WriteString(*upstream + "\n")

		for _, domainName := range route.DomainNames {
			mapping.WriteString(fmt.Sprintf("%s %s;\n", domainName, routeId))
		}
	}

	defaultUpstreamId := fmt.Sprintf("stream_%s_default", s.ID)
	defaultUpstream, err := p.buildUpstream([]stream.Backend{s.DefaultBackend}, defaultUpstreamId)
	if err != nil {
		return nil, err
	}

	upstreams.WriteString(*defaultUpstream + "\n")
	mapping.WriteString(fmt.Sprintf("default %s;\n}", defaultUpstreamId))
	instructions := fmt.Sprintf(
		`
			ssl_preread on;
			proxy_pass %s;
		`,
		mappingId,
	)
	return p.buildStream(s, upstreams.String()+mapping.String(), instructions)
}

func (p *streamFileProvider) buildStream(s *stream.Stream, upstreams, instructions string) (*string, error) {
	binding, err := p.buildBinding(s)
	if err != nil {
		return nil, err
	}

	tcpNoDelay := ""
	if s.Binding.Protocol == stream.TCPProtocol && s.FeatureSet.TCPNoDelay {
		tcpNoDelay = "tcp_nodelay on;"
	}

	socketKeepAlive := ""
	if s.FeatureSet.SocketKeepAlive {
		socketKeepAlive = "proxy_socket_keepalive on;"
	}

	contents := fmt.Sprintf(
		`
		%s 

		server {
			%s
			%s
			%s
			%s
		}
		`,
		upstreams,
		*binding,
		tcpNoDelay,
		socketKeepAlive,
		instructions,
	)

	return &contents, nil
}
