package settings

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

type service struct {
	repository      Repository
	bindingCommands *binding.Commands
	scheduler       *scheduler.Scheduler
}

func newService(
	repository Repository,
	bindingCommands *binding.Commands,
	sched *scheduler.Scheduler,
) *service {
	return &service{
		repository:      repository,
		bindingCommands: bindingCommands,
		scheduler:       sched,
	}
}

func (s *service) get(ctx context.Context) (*Settings, error) {
	return s.repository.Get(ctx)
}

func (s *service) save(ctx context.Context, settings *Settings) error {
	if err := newValidator(s.bindingCommands).validate(ctx, settings); err != nil {
		return err
	}

	if err := s.repository.Save(ctx, settings); err != nil {
		return err
	}

	return s.scheduler.Reload(ctx)
}
