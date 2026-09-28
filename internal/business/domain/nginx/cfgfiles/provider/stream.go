package provider

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/runtime"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
)

type streamProvider struct{}

func newStreamProvider() *streamProvider {
	return &streamProvider{}
}

func (p *streamProvider) Provide(ctx *Context) ([]File, error) {
	if len(ctx.Streams) > 0 && ctx.SupportedFeatures.StreamType == NoneSupportType {
		return nil, coreerror.New(
			i18n.M(ctx.Context, i18n.K.CoreNginxCfgfilesStreamNotEnabled),
			false,
		)
	}

	files := make([]File, 0, len(ctx.Streams))

	for _, s := range ctx.Streams {
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

func (p *streamProvider) buildConfigFileContents(
	ctx *Context,
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

func (p *streamProvider) buildSimpleStream(s *stream.Stream) (*string, error) {
	upstreamID := fmt.Sprintf("stream_%s_default", nginxID(s))
	upstream, err := p.buildUpstream([]stream.Backend{s.DefaultBackend}, upstreamID)
	if err != nil {
		return nil, err
	}

	return p.buildStream(
		s,
		*upstream,
		nginxSprintf("proxy_pass %s;", upstreamID),
	)
}

func (p *streamProvider) buildBinding(s *stream.Stream) (*string, error) {
	instruction := strings.Builder{}
	_, _ = instruction.WriteString("listen ")

	switch s.Binding.Protocol {
	case stream.SocketProtocol:
		nginxFprintf(&instruction, "%s", "unix:"+s.Binding.Address)

	case stream.TCPProtocol:
		nginxFprintf(
			&instruction,
			"%s",
			net.JoinHostPort(s.Binding.Address, strconv.Itoa(*s.Binding.Port)),
		)

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
		nginxFprintf(
			&instruction,
			"%s udp",
			net.JoinHostPort(s.Binding.Address, strconv.Itoa(*s.Binding.Port)),
		)

	default:
		return nil, fmt.Errorf("unknown binding protocol: %s", s.Binding.Protocol)
	}

	if runtime.IsWindows() {
		_, _ = instruction.WriteString(";")
	} else {
		_, _ = instruction.WriteString(" reuseport;")
	}

	return new(instruction.String()), nil
}

func (p *streamProvider) buildUpstream(
	backends []stream.Backend,
	name string,
) (*string, error) {
	instructions := strings.Builder{}
	nginxFprintf(&instructions, "upstream %s {\n", name)

	for _, backend := range backends {
		address := backend.Address
		switch address.Protocol {
		case stream.SocketProtocol:
			nginxFprintf(&instructions, "server %s", "unix:"+address.Address)

		case stream.TCPProtocol, stream.UDPProtocol:
			nginxFprintf(
				&instructions,
				"server %s",
				net.JoinHostPort(address.Address, strconv.Itoa(*address.Port)),
			)

		default:
			return nil, fmt.Errorf("unknown backend protocol: %s", address.Protocol)
		}

		if backend.Weight != nil {
			nginxFprintf(&instructions, " weight=%d", *backend.Weight)
		}

		if backend.CircuitBreaker != nil {
			nginxFprintf(
				&instructions,
				" max_fails=%d fail_timeout=%ds",
				backend.CircuitBreaker.MaxFailures,
				backend.CircuitBreaker.OpenSeconds,
			)
		}

		_, _ = instructions.WriteString(";\n")
	}

	_, _ = instructions.WriteString("}\n")
	return new(instructions.String()), nil
}

func (p *streamProvider) buildRoutedStream(
	ctx *Context,
	s *stream.Stream,
) (*string, error) {
	if ctx.SupportedFeatures.TLSSNI == NoneSupportType {
		return nil, coreerror.New(
			i18n.M(ctx.Context, i18n.K.CoreNginxCfgfilesStreamSniNotEnabled),
			false,
		)
	}

	mapping := strings.Builder{}
	mappingID := fmt.Sprintf("$stream_%s_router", nginxID(s))
	nginxFprintf(&mapping, "map $ssl_preread_server_name %s {\n", mappingID)

	upstreams := strings.Builder{}
	for routeIndex, route := range s.Routes {
		routeID := fmt.Sprintf("stream_%s_route_%d", nginxID(s), routeIndex)
		upstream, err := p.buildUpstream(route.Backends, routeID)
		if err != nil {
			return nil, err
		}

		_, _ = upstreams.WriteString(*upstream + "\n")

		for _, domainName := range route.DomainNames {
			nginxFprintf(&mapping, "%s %s;\n", domainName, routeID)
		}
	}

	defaultUpstreamID := fmt.Sprintf("stream_%s_default", nginxID(s))
	defaultUpstream, err := p.buildUpstream([]stream.Backend{s.DefaultBackend}, defaultUpstreamID)
	if err != nil {
		return nil, err
	}

	_, _ = upstreams.WriteString(*defaultUpstream + "\n")
	nginxFprintf(&mapping, "default %s;\n}", defaultUpstreamID)
	instructions := nginxSprintf(
		`
			ssl_preread on;
			proxy_pass %s;
		`,
		mappingID,
	)
	return p.buildStream(s, upstreams.String()+mapping.String(), instructions)
}

func (p *streamProvider) buildStream(
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

	return new(nginxSprintf(
		`
		%s 

		server {
			%s
			%s
			%s
			%s
		}
		`,
		directiveFragment(upstreams),
		directiveFragment(*binding),
		directiveFragment(tcpNoDelay),
		directiveFragment(socketKeepAlive),
		directiveFragment(instructions),
	)), nil
}

func nginxID(s *stream.Stream) string {
	return strings.ReplaceAll(s.ID.String(), "-", "")
}
