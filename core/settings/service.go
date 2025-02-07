package settings

import (
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

func (s *service) get() (*Settings, error) {
	return (*s.repository).Get()
}

func (s *service) save(settings *Settings) error {
	if err := newValidator(s.validateBindingCommand).validate(settings); err != nil {
		return err
	}

	if err := (*s.repository).Save(settings); err != nil {
		return err
	}

	return (*s.scheduler).Reload()
}
