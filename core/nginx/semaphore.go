package nginx

import (
	"sync"
)

type nginxState int

const (
	stoppedState = nginxState(0)
	runningState = nginxState(1)
)

type semaphore struct {
	lock  sync.Mutex
	state nginxState
}

func newSemaphore() *semaphore {
	return &semaphore{
		state: stoppedState,
	}
}

func (s *semaphore) changeState(newState nginxState, action func() error) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := action()
	if err == nil {
		s.state = newState
	}

	return err
}

func (s *semaphore) currentState() nginxState {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.state
}
