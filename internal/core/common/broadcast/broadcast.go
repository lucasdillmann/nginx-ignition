package broadcast

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/log"
)

var channels = make(map[string]chan context.Context)

func SendSignal(ctx context.Context, qualifier string) {
	if channels[qualifier] != nil {
		channels[qualifier] <- ctx
	} else {
		log.Warnf("Signal ignored: qualifier %s has no listeners yet", qualifier)
	}
}

func Listen(qualifier string) chan context.Context {
	if channels[qualifier] == nil {
		channels[qualifier] = make(chan context.Context)
	}

	return channels[qualifier]
}
