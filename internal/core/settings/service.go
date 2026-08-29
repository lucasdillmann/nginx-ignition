package settings

import (
	"context"

	"nginx-ignition/internal/core/binding"
	"nginx-ignition/internal/core/common/scheduler"
)

type service struct {
	repository      Repository
	bindingCommands binding.Commands
	scheduler       *scheduler.Scheduler
}

func newCommands(
	repository Repository,
	bindingCommands binding.Commands,
	sched *scheduler.Scheduler,
) Commands {
	return &service{
		repository:      repository,
		bindingCommands: bindingCommands,
		scheduler:       sched,
	}
}

func (s *service) Get(ctx context.Context) (*Settings, error) {
	return s.repository.Get(ctx)
}

func (s *service) Save(ctx context.Context, settings *Settings) error {
	if err := newValidator(s.bindingCommands).validate(ctx, settings); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, settings); err != nil {
		return err
	}

	return s.scheduler.Reload(ctx)
}
