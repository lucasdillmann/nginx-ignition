package cfgfiles

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/host"
)

type fileProvider interface {
	provide(ctx context.Context, basePath string, hosts []*host.Host) ([]output, error)
}

type output struct {
	name     string
	contents string
}
