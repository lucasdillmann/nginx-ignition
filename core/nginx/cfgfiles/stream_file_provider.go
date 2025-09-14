package cfgfiles

import (
	"fmt"
	"strings"

	"dillmann.com.br/nginx-ignition/core/stream"
	"github.com/aws/smithy-go/ptr"
)

type streamFileProvider struct{}

func newStreamFileProvider() *streamFileProvider {
	return &streamFileProvider{}
}

func (p *streamFileProvider) provide(ctx *providerContext) ([]output, error) {
	files := make([]output, 0, len(ctx.streams))

	for _, s := range ctx.streams {
		contents, err := p.buildConfigFileContents(s)
		if err != nil {
			return nil, err
		}

		files = append(files, output{
			name:     fmt.Sprintf("stream-%s.conf", s.ID),
			contents: *contents,
		})
	}

	return files, nil
}

func (p *streamFileProvider) buildConfigFileContents(s *stream.Stream) (*string, error) {
	binding, err := p.buildBinding(s)
	if err != nil {
		return nil, err
	}

	backend, err := p.buildBackend(s)
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
		server {
			%s
			%s
			%s
			%s
		}
		`,
		*binding,
		*backend,
		tcpNoDelay,
		socketKeepAlive,
	)

	return &contents, nil
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

func (p *streamFileProvider) buildBackend(s *stream.Stream) (*string, error) {
	instruction := strings.Builder{}
	instruction.WriteString("proxy_pass ")

	switch s.Backend.Protocol {
	case stream.SocketProtocol:
		instruction.WriteString(fmt.Sprintf("unix:%s", s.Backend.Address))

	case stream.TCPProtocol, stream.UDPProtocol:
		instruction.WriteString(fmt.Sprintf("%s:%d", s.Backend.Address, *s.Backend.Port))

	default:
		return nil, fmt.Errorf("unknown backend protocol: %s", s.Backend.Protocol)
	}

	instruction.WriteString(";")

	return ptr.String(instruction.String()), nil
}
