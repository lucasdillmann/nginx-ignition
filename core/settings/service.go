package settings

import "dillmann.com.br/nginx-ignition/core/host"

type service struct {
	repository             *Repository
	validateBindingCommand *host.ValidateBindingCommand
}

func newService(repository *Repository, validateBindingCommand *host.ValidateBindingCommand) *service {
	return &service{
		repository:             repository,
		validateBindingCommand: validateBindingCommand,
	}
}

func (s *service) get() (*Settings, error) {
	return (*s.repository).Get()
}

func (s *service) save(settings *Settings) error {
	if err := newValidator(s.validateBindingCommand).validate(settings); err != nil {
		return err
	}

	return (*s.repository).Save(settings)
}
