package broadcast

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type shutdown struct{}

func registerShutdown(lc *lifecycle.Lifecycle) {
	lc.RegisterShutdown(shutdown{})
}

func (s shutdown) Priority() int {
	return shutdownPriority
}

func (s shutdown) Run(_ context.Context) {
	for _, ch := range channels {
		close(ch)
	}
}
