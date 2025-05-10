package cfgfiles

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type providerContext struct {
	context  context.Context
	basePath string
	hosts    []*host.Host
	streams  []*stream.Stream
}

type fileProvider interface {
	provide(ctx *providerContext) ([]output, error)
}

type output struct {
	name     string
	contents string
}
