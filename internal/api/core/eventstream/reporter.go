package eventstream

import "github.com/lucasdillmann/nginx-ignition/internal/business/core/progress"

type streamReporter struct {
	events chan progress.Event
	done   <-chan struct{}
}

func (r *streamReporter) Report(event progress.Event) {
	select {
	case r.events <- event:
	case <-r.done:
	}
}
