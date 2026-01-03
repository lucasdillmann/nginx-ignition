package settings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
)

func validSettings() *Settings {
	return &Settings{
		Nginx: &NginxSettings{
			DefaultContentType: "text/html",
			RuntimeUser:        "nginx",
			Timeouts: &NginxTimeoutsSettings{
				Read:      60,
				Send:      60,
				Connect:   60,
				Keepalive: 65,
			},
			WorkerProcesses:   1,
			WorkerConnections: 1024,
			MaximumBodySizeMb: 1,
		},
		LogRotation: &LogRotationSettings{
			IntervalUnitCount: 1,
			MaximumLines:      1000,
		},
		CertificateAutoRenew: &CertificateAutoRenewSettings{
			IntervalUnitCount: 1,
		},
	}
}

func Test_Validator(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bindingCommands := binding.NewMockedCommands(ctrl)

	t.Run("validate", func(t *testing.T) {
		t.Run("valid settings pass", func(t *testing.T) {
			settings := validSettings()
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.NoError(t, err)
		})

		t.Run("empty default content type fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.DefaultContentType = ""
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("default content type exceeds maximum length fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.DefaultContentType = string(make([]byte, 129))
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("empty runtime user fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.RuntimeUser = ""
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("whitespace-only runtime user fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.RuntimeUser = "   "
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("runtime user exceeds maximum length fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.RuntimeUser = string(make([]byte, 33))
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("runtime user at maximum length passes", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.RuntimeUser = string(make([]byte, 32))
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.NoError(t, err)
		})

		t.Run("default content type at maximum length passes", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.DefaultContentType = string(make([]byte, 128))
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.NoError(t, err)
		})

		t.Run("timeout read below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.Timeouts.Read = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("timeout send below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.Timeouts.Send = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("timeout connect below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.Timeouts.Connect = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("timeout keepalive below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.Timeouts.Keepalive = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("worker processes below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.WorkerProcesses = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("worker processes above range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.WorkerProcesses = 101
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("worker connections below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.WorkerConnections = 31
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("worker connections above range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.WorkerConnections = 4097
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("maximum body size below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.Nginx.MaximumBodySizeMb = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("log rotation interval unit count below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.LogRotation.IntervalUnitCount = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("log rotation maximum lines below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.LogRotation.MaximumLines = -1
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("log rotation maximum lines above range fails", func(t *testing.T) {
			settings := validSettings()
			settings.LogRotation.MaximumLines = 10001
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})

		t.Run("certificate auto renew interval unit count below range fails", func(t *testing.T) {
			settings := validSettings()
			settings.CertificateAutoRenew.IntervalUnitCount = 0
			val := newValidator(bindingCommands)

			err := val.validate(ctx, settings)

			assert.Error(t, err)
		})
	})
}
