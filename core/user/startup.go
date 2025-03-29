package user

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type startup struct {
	service       *service
	configuration *configuration.Configuration
}

func registerStartup(lifecycle *lifecycle.Lifecycle, service *service, configuration *configuration.Configuration) {
	commandInstance := &startup{service, configuration}
	lifecycle.RegisterStartup(commandInstance)
}

func (s startup) Run(ctx context.Context) error {
	username, err := s.configuration.Get("nginx-ignition.password-reset.username")
	if err != nil || username == "" {
		return nil
	}

	newPassword, err := s.service.resetPassword(ctx, username)
	if err != nil {
		log.Errorf("Error resetting password for the user %s: %s", username, err)
		return err
	}

	log.Infof("Password reset completed successfully for the user %s. New password: %s", username, newPassword)
	return core_error.New(
		"Application was started using the password reset procedure. Please disable it in order to continue.",
		true,
	)
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return false
}
