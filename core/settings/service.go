package settings

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/host"
)

type service struct {
	repository             *Repository
	validateBindingCommand *host.ValidateBindingCommand
	scheduler              *scheduler.Scheduler
}

func newService(
	repository *Repository,
	validateBindingCommand *host.ValidateBindingCommand,
	scheduler *scheduler.Scheduler,
) *service {
	return &service{
		repository:             repository,
		validateBindingCommand: validateBindingCommand,
		scheduler:              scheduler,
	}
}

func (s *service) get(ctx context.Context) (*Settings, error) {
	return (*s.repository).Get(ctx)
}

func (s *service) save(ctx context.Context, settings *Settings) error {
	if err := newValidator(s.validateBindingCommand).validate(ctx, settings); err != nil {
		return err
	}

	if err := (*s.repository).Save(ctx, settings); err != nil {
		return err
	}

	return (*s.scheduler).Reload(ctx)
}
