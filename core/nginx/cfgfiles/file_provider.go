package cfgfiles

import (
	"dillmann.com.br/nginx-ignition/core/host"
)

type fileProvider interface {
	provide(basePath string, hosts []*host.Host) ([]output, error)
}

type output struct {
	name     string
	contents string
}
