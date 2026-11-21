package eventstream

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"
)

func Serve(ctx *gin.Context, run func(context.Context) (any, error)) {
	writer := New(ctx.Writer)
	if err := writer.Prepare(); err != nil {
		return
	}

	requestContext := ctx.Request.Context()
	done := make(chan struct{})
	defer close(done)

	events := make(chan progress.Event, eventBufferSize)
	runContext, cancel := context.WithCancel(requestContext)
	defer cancel()

	runContext = progress.WithReporter(runContext, &streamReporter{
		events: events,
		done:   done,
	})

	results := make(chan streamResult, 1)
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Errorf(
					"Panic detected while running a streaming operation: %v\n%s",
					recovered,
					debug.Stack(),
				)
				results <- streamResult{panic: recovered}
			}
		}()

		payload, err := run(runContext)
		results <- streamResult{payload: payload, err: err}
	}()

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case event := <-events:
			if err := writer.Event("log", event); err != nil {
				return
			}
		case result := <-results:
			if err := drainEvents(writer, events); err != nil {
				return
			}

			writeTerminal(requestContext, writer, result)
			return
		case <-ticker.C:
			if err := writer.Comment("keep-alive"); err != nil {
				return
			}
		case <-requestContext.Done():
			return
		}
	}
}
