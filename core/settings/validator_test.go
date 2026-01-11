package settings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
)

func Test_validator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bindingCommands := binding.NewMockedCommands(ctrl)

	t.Run("validate", func(t *testing.T) {
		t.Run("valid settings pass", func(t *testing.T) {
			s := newSettings()
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.NoError(t, err)
		})

		t.Run("empty default content type fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.DefaultContentType = ""
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("default content type exceeds maximum length fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.DefaultContentType = strings.Repeat("a", 129)
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("empty runtime user fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.RuntimeUser = ""
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("whitespace-only runtime user fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.RuntimeUser = "   "
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("runtime user exceeds maximum length fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.RuntimeUser = strings.Repeat("a", 33)
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("runtime user at maximum length passes", func(t *testing.T) {
			s := newSettings()
			s.Nginx.RuntimeUser = strings.Repeat("a", 32)
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.NoError(t, err)
		})

		t.Run("default content type at maximum length passes", func(t *testing.T) {
			s := newSettings()
			s.Nginx.DefaultContentType = strings.Repeat("a", 128)
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.NoError(t, err)
		})

		t.Run("timeout read below range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.Timeouts.Read = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("timeout send below range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.Timeouts.Send = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("timeout connect below range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.Timeouts.Connect = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("timeout keepalive below range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.Timeouts.Keepalive = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("worker processes below range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.WorkerProcesses = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("worker processes above range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.WorkerProcesses = 101
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("worker connections below range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.WorkerConnections = 31
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("worker connections above range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.WorkerConnections = 4097
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("maximum body size below range fails", func(t *testing.T) {
			s := newSettings()
			s.Nginx.MaximumBodySizeMb = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("log rotation interval unit count below range fails", func(t *testing.T) {
			s := newSettings()
			s.LogRotation.IntervalUnitCount = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("log rotation maximum lines below range fails", func(t *testing.T) {
			s := newSettings()
			s.LogRotation.MaximumLines = -1
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("log rotation maximum lines above range fails", func(t *testing.T) {
			s := newSettings()
			s.LogRotation.MaximumLines = 10001
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("certificate auto renew interval unit count below range fails", func(t *testing.T) {
			s := newSettings()
			s.CertificateAutoRenew.IntervalUnitCount = 0
			settingsValidator := newValidator(bindingCommands)

			err := settingsValidator.validate(t.Context(), s)

			assert.Error(t, err)
		})
	})
}
