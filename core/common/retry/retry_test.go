package retry

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry_Start(t *testing.T) {
	t.Run("returns error when action is nil", func(t *testing.T) {
		retry := &Retry{
			Action:               nil,
			Attempts:             3,
			DelayBetweenAttempts: time.Millisecond,
		}

		err := retry.Start()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "action must be set")
	})

	t.Run("returns error when attempts is less than 1", func(t *testing.T) {
		retry := &Retry{
			Action:               func() error { return nil },
			Attempts:             0,
			DelayBetweenAttempts: time.Millisecond,
		}

		err := retry.Start()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "attempts must be greater than 0")
	})

	t.Run("returns error when delay is negative", func(t *testing.T) {
		retry := &Retry{
			Action:               func() error { return nil },
			Attempts:             3,
			DelayBetweenAttempts: -1,
		}

		err := retry.Start()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delay between attempts must be greater than or equal to 0")
	})

	t.Run("starts retry successfully", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		retry := &Retry{
			Action: func() error {
				wg.Done()
				return nil
			},
			Attempts:             1,
			DelayBetweenAttempts: time.Millisecond,
		}

		err := retry.Start()
		assert.NoError(t, err)

		wg.Wait()
	})
}

func TestRetry_executeAttempts(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		var callbackCalled bool
		var callbackAttempt int
		var callbackCompleted bool

		retry := &Retry{
			Action: func() error {
				return nil
			},
			Callback: func(_ error, attempt int, completed bool) {
				callbackCalled = true
				callbackAttempt = attempt
				callbackCompleted = completed
			},
			Attempts:             3,
			DelayBetweenAttempts: time.Millisecond,
		}

		retry.executeAttempts()

		assert.True(t, callbackCalled)
		assert.Equal(t, 0, callbackAttempt)
		assert.True(t, callbackCompleted)
	})

	t.Run("retries on failure and succeeds", func(t *testing.T) {
		attemptCount := 0
		var callbackAttempts []int

		retry := &Retry{
			Action: func() error {
				attemptCount++
				if attemptCount < 2 {
					return errors.New("temporary error")
				}
				return nil
			},
			Callback: func(_ error, attempt int, _ bool) {
				callbackAttempts = append(callbackAttempts, attempt)
			},
			Attempts:             3,
			DelayBetweenAttempts: time.Millisecond,
		}

		retry.executeAttempts()

		assert.Len(t, callbackAttempts, 2)
		assert.Equal(t, 0, callbackAttempts[0])
		assert.Equal(t, 1, callbackAttempts[1])
	})

	t.Run("exhausts all attempts on persistent failure", func(t *testing.T) {
		var callbackAttempts []int
		var callbackCompleted []bool

		retry := &Retry{
			Action: func() error {
				return errors.New("persistent error")
			},
			Callback: func(_ error, attempt int, completed bool) {
				callbackAttempts = append(callbackAttempts, attempt)
				callbackCompleted = append(callbackCompleted, completed)
			},
			Attempts:             3,
			DelayBetweenAttempts: time.Millisecond,
		}

		retry.executeAttempts()

		assert.Len(t, callbackAttempts, 3)
		assert.Equal(t, []int{0, 1, 2}, callbackAttempts)
		assert.Equal(t, []bool{false, false, true}, callbackCompleted)
	})

	t.Run("handles panic in action", func(t *testing.T) {
		var callbackCalled bool
		var callbackErr error

		retry := &Retry{
			Action: func() error {
				panic("test panic")
			},
			Callback: func(err error, _ int, _ bool) {
				callbackCalled = true
				callbackErr = err
			},
			Attempts:             1,
			DelayBetweenAttempts: time.Millisecond,
		}

		retry.executeAttempts()

		assert.True(t, callbackCalled)
		assert.Error(t, callbackErr)
		assert.Contains(t, callbackErr.Error(), "panic")
	})
}

func TestRetry_sendCallback(t *testing.T) {
	t.Run("calls callback when set", func(t *testing.T) {
		var callbackCalled bool
		var callbackErr error
		var callbackAttempt int
		var callbackCompleted bool

		retry := &Retry{
			Callback: func(err error, attempt int, completed bool) {
				callbackCalled = true
				callbackErr = err
				callbackAttempt = attempt
				callbackCompleted = completed
			},
		}

		testErr := errors.New("test error")
		retry.sendCallback(testErr, 5, true)

		assert.True(t, callbackCalled)
		assert.Equal(t, testErr, callbackErr)
		assert.Equal(t, 5, callbackAttempt)
		assert.True(t, callbackCompleted)
	})

	t.Run("does not panic when callback is nil", func(t *testing.T) {
		retry := &Retry{
			Callback: nil,
		}

		assert.NotPanics(t, func() {
			retry.sendCallback(errors.New("test"), 0, false)
		})
	})

	t.Run("calls callback with error details", func(t *testing.T) {
		var receivedErr error
		var receivedAttempt int
		var receivedCompleted bool

		retry := &Retry{
			Callback: func(err error, attempt int, completed bool) {
				receivedErr = err
				receivedAttempt = attempt
				receivedCompleted = completed
			},
		}

		testErr := errors.New("test error")
		retry.sendCallback(testErr, 5, true)

		assert.Equal(t, testErr, receivedErr)
		assert.Equal(t, 5, receivedAttempt)
		assert.True(t, receivedCompleted)
	})
}
