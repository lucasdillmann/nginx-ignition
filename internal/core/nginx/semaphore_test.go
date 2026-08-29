package nginx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_semaphore(t *testing.T) {
	t.Run("changeState", func(t *testing.T) {
		t.Run("updates state when action succeeds", func(t *testing.T) {
			sem := newSemaphore()
			err := sem.changeState(runningState, func() error {
				return nil
			})

			assert.NoError(t, err)
			assert.Equal(t, runningState, sem.currentState())
		})

		t.Run("does not update state when action fails", func(t *testing.T) {
			sem := newSemaphore()
			expectedErr := errors.New("failed")
			err := sem.changeState(runningState, func() error {
				return expectedErr
			})

			assert.Equal(t, expectedErr, err)
			assert.Equal(t, stoppedState, sem.currentState())
		})
	})

	t.Run("currentState", func(t *testing.T) {
		t.Run("returns initial state", func(t *testing.T) {
			sem := newSemaphore()
			assert.Equal(t, stoppedState, sem.currentState())
		})
	})
}
