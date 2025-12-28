package retry

import (
	"errors"
	"fmt"
	"runtime/debug"
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
		return errors.New("action must be set")
	}

	if r.Attempts < 1 {
		return errors.New("attempts must be greater than 0")
	}

	if r.DelayBetweenAttempts < 0 {
		return errors.New("delay between attempts must be greater than or equal to 0")
	}

	go func() {
		defer func() {
			if errDetails := recover(); errDetails != nil {
				err := fmt.Errorf("panic: %v\n%s", errDetails, debug.Stack())
				r.sendCallback(err, -1, true)
			}
		}()

		r.executeAttempts()
	}()

	return nil
}

func (r *Retry) executeAttempts() {
	for attempt := 0; attempt < r.Attempts; attempt++ {
		var err error

		func() {
			defer func() {
				if errDetails := recover(); errDetails != nil {
					err = fmt.Errorf("panic: %v\n%s", errDetails, debug.Stack())
				}
			}()

			err = r.Action()
		}()

		if err != nil {
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
