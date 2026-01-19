package cfgfiles

import (
	"fmt"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/common/runtime"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type streamFileProvider struct{}

func newStreamFileProvider() *streamFileProvider {
	return &streamFileProvider{}
}

func (p *streamFileProvider) provide(ctx *providerContext) ([]File, error) {
	if len(ctx.streams) > 0 && ctx.supportedFeatures.StreamType == NoneSupportType {
		return nil, coreerror.New(
			i18n.M(ctx.context, i18n.K.CoreNginxCfgfilesStreamNotEnabled),
			false,
		)
	}

	files := make([]File, 0, len(ctx.streams))

	for _, s := range ctx.streams {
		contents, err := p.buildConfigFileContents(ctx, &s)
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

func (p *streamFileProvider) buildConfigFileContents(
	ctx *providerContext,
	s *stream.Stream,
) (*string, error) {
	switch s.Type {
	case stream.SimpleType:
		return p.buildSimpleStream(s)
	case stream.SNIRouterType:
		return p.buildRoutedStream(ctx, s)
	default:
		return nil, fmt.Errorf("unknown stream type: %s", s.Type)
	}
}

func (p *streamFileProvider) buildSimpleStream(s *stream.Stream) (*string, error) {
	upstreamID := fmt.Sprintf("stream_%s_default", nginxID(s))
	upstream, err := p.buildUpstream([]stream.Backend{s.DefaultBackend}, upstreamID)
	if err != nil {
		return nil, err
	}

	return p.buildStream(s, *upstream, fmt.Sprintf("proxy_pass %s;", upstreamID))
}

func (p *streamFileProvider) buildBinding(s *stream.Stream) (*string, error) {
	instruction := strings.Builder{}
	_, _ = instruction.WriteString("listen ")

	switch s.Binding.Protocol {
	case stream.SocketProtocol:
		_, _ = fmt.Fprintf(&instruction, "unix:\"%s\"", s.Binding.Address)

	case stream.TCPProtocol:
		_, _ = fmt.Fprintf(&instruction, "%s:%d", s.Binding.Address, *s.Binding.Port)

		if s.FeatureSet.UseProxyProtocol {
			_, _ = instruction.WriteString(" proxy_protocol")
		}

		if s.FeatureSet.TCPDeferred {
			_, _ = instruction.WriteString(" deferred")
		}

		if s.FeatureSet.TCPKeepAlive {
			_, _ = instruction.WriteString(" so_keepalive=on")
		}

	case stream.UDPProtocol:
		_, _ = fmt.Fprintf(&instruction, "%s:%d udp", s.Binding.Address, *s.Binding.Port)

	default:
		return nil, fmt.Errorf("unknown binding protocol: %s", s.Binding.Protocol)
	}

	if runtime.IsWindows() {
		_, _ = instruction.WriteString(";")
	} else {
		_, _ = instruction.WriteString(" reuseport;")
	}

	return ptr.Of(instruction.String()), nil
}

func (p *streamFileProvider) buildUpstream(
	backends []stream.Backend,
	name string,
) (*string, error) {
	instructions := strings.Builder{}
	_, _ = fmt.Fprintf(&instructions, "upstream %s {\n", name)

	for _, backend := range backends {
		address := backend.Address
		switch address.Protocol {
		case stream.SocketProtocol:
			_, _ = fmt.Fprintf(&instructions, "server unix:\"%s\"", address.Address)

		case stream.TCPProtocol, stream.UDPProtocol:
			_, _ = fmt.Fprintf(&instructions, "server %s:%d", address.Address, *address.Port)

		default:
			return nil, fmt.Errorf("unknown backend protocol: %s", address.Protocol)
		}

		if backend.Weight != nil {
			_, _ = fmt.Fprintf(&instructions, " weight=%d", *backend.Weight)
		}

		if backend.CircuitBreaker != nil {
			_, _ = fmt.Fprintf(
				&instructions,
				" max_fails=%d fail_timeout=%ds",
				backend.CircuitBreaker.MaxFailures,
				backend.CircuitBreaker.OpenSeconds,
			)
		}

		_, _ = instructions.WriteString(";\n")
	}

	_, _ = instructions.WriteString("}\n")
	return ptr.Of(instructions.String()), nil
}

func (p *streamFileProvider) buildRoutedStream(
	ctx *providerContext,
	s *stream.Stream,
) (*string, error) {
	if ctx.supportedFeatures.TLSSNI == NoneSupportType {
		return nil, coreerror.New(
			i18n.M(ctx.context, i18n.K.CoreNginxCfgfilesStreamSniNotEnabled),
			false,
		)
	}

	mapping := strings.Builder{}
	mappingID := fmt.Sprintf("$stream_%s_router", nginxID(s))
	_, _ = fmt.Fprintf(&mapping, "map $ssl_preread_server_name %s {\n", mappingID)

	upstreams := strings.Builder{}
	for routeIndex, route := range s.Routes {
		routeID := fmt.Sprintf("stream_%s_route_%d", nginxID(s), routeIndex)
		upstream, err := p.buildUpstream(route.Backends, routeID)
		if err != nil {
			return nil, err
		}

		_, _ = upstreams.WriteString(*upstream + "\n")

		for _, domainName := range route.DomainNames {
			_, _ = fmt.Fprintf(&mapping, "%s %s;\n", domainName, routeID)
		}
	}

	defaultUpstreamID := fmt.Sprintf("stream_%s_default", nginxID(s))
	defaultUpstream, err := p.buildUpstream([]stream.Backend{s.DefaultBackend}, defaultUpstreamID)
	if err != nil {
		return nil, err
	}

	_, _ = upstreams.WriteString(*defaultUpstream + "\n")
	_, _ = fmt.Fprintf(&mapping, "default %s;\n}", defaultUpstreamID)
	instructions := fmt.Sprintf(
		`
			ssl_preread on;
			proxy_pass %s;
		`,
		mappingID,
	)
	return p.buildStream(s, upstreams.String()+mapping.String(), instructions)
}

func (p *streamFileProvider) buildStream(
	s *stream.Stream,
	upstreams, instructions string,
) (*string, error) {
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

func nginxID(s *stream.Stream) string {
	return strings.ReplaceAll(s.ID.String(), "-", "")
}
