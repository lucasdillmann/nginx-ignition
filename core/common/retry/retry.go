package retry

import (
	"fmt"
	"time"
)

type Retry struct {
	Action               func() error
	Callback             func(err error, attempt int, completed bool)
	Attempts             int
	DelayBetweenAttempts time.Duration
}

func (r *Retry) Start() error {
	if r.Action == nil {
		return fmt.Errorf("action must be set")
	}

	if r.Attempts < 1 {
		return fmt.Errorf("attempts must be greater than 0")
	}

	if r.DelayBetweenAttempts < 0 {
		return fmt.Errorf("delay between attempts must be greater than or equal to 0")
	}

	go func() {
		if errDetails := recover(); errDetails != nil {
			err := fmt.Errorf("%v", errDetails)
			r.sendCallback(err, -1, true)
			return
		}

		r.executeAttempts()
	}()

	return nil
}

func (r *Retry) executeAttempts() {
	for attempt := 0; attempt < r.Attempts; attempt++ {
		if err := r.Action(); err != nil {
			willRetry := attempt+1 < r.Attempts
			r.sendCallback(err, attempt, !willRetry)
			time.Sleep(r.DelayBetweenAttempts)
			continue
		}

		r.sendCallback(nil, attempt, true)
		return
	}
}

func (r *Retry) sendCallback(err error, attempt int, completed bool) {
	if r.Callback != nil {
		r.Callback(err, attempt, completed)
	}
}
