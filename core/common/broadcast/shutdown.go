package broadcast

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type shutdown struct {
}

func registerShutdown(lifecycle *lifecycle.Lifecycle) {
	lifecycle.RegisterShutdown(shutdown{})
}

func (s shutdown) Priority() int {
	return shutdownPriority
}

func (s shutdown) Run() {
	for _, ch := range channels {
		close(ch)
	}
}
