package broadcast

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/lifecycle"
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
