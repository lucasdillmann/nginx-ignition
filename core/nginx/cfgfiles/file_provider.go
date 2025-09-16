package cfgfiles

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type providerContext struct {
	context context.Context
	paths   *Paths
	hosts   []*host.Host
	streams []*stream.Stream
}

type Paths struct {
	AbsoluteConfig string
	AbsoluteLogs   string
}

type fileProvider interface {
	provide(ctx *providerContext) ([]File, error)
}

type File struct {
	Name     string
	Contents string
}
