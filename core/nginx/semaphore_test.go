package nginx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Semaphore_ChangeState(t *testing.T) {
	t.Run("updates state when action succeeds", func(t *testing.T) {
		s := newSemaphore()
		err := s.changeState(runningState, func() error {
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, runningState, s.currentState())
	})

	t.Run("does not update state when action fails", func(t *testing.T) {
		s := newSemaphore()
		expectedErr := errors.New("failed")
		err := s.changeState(runningState, func() error {
			return expectedErr
		})

		assert.Equal(t, expectedErr, err)
		assert.Equal(t, stoppedState, s.currentState())
	})
}

func Test_Semaphore_CurrentState(t *testing.T) {
	t.Run("returns initial state", func(t *testing.T) {
		s := newSemaphore()
		assert.Equal(t, stoppedState, s.currentState())
	})
}
