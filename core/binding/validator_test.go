package binding

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func Test_validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid HTTP binding passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPBinding()
			certificateCommands := certificate.NewMockedCommands(ctrl)
			bindingValidator := newValidator(validation.NewValidator(), certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
		})

		t.Run("valid HTTPS binding with certificate passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPSBinding()
			certificateCommands := certCommandsExists(ctrl, *binding.CertificateID)
			bindingValidator := newValidator(validation.NewValidator(), certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
		})

		t.Run("valid IPv6 address passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPBinding()
			binding.IP = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
			certificateCommands := certificate.NewMockedCommands(ctrl)
			bindingValidator := newValidator(validation.NewValidator(), certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
		})

		t.Run("invalid IP fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPBinding()
			binding.IP = "invalid.ip"
			certificateCommands := certificate.NewMockedCommands(ctrl)
			delegate := validation.NewValidator()
			bindingValidator := newValidator(delegate, certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
			assert.Error(t, delegate.Result())
		})

		t.Run("port below range fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPBinding()
			binding.Port = 0
			certificateCommands := certificate.NewMockedCommands(ctrl)
			delegate := validation.NewValidator()
			bindingValidator := newValidator(delegate, certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
			assert.Error(t, delegate.Result())
		})

		t.Run("port above range fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPBinding()
			binding.Port = 65536
			certificateCommands := certificate.NewMockedCommands(ctrl)
			delegate := validation.NewValidator()
			bindingValidator := newValidator(delegate, certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
			assert.Error(t, delegate.Result())
		})

		t.Run("HTTP binding with certificate fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPBinding()
			binding.CertificateID = new(uuid.New())
			certificateCommands := certificate.NewMockedCommands(ctrl)
			delegate := validation.NewValidator()
			bindingValidator := newValidator(delegate, certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
			assert.Error(t, delegate.Result())
		})

		t.Run("HTTPS binding without certificate fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPSBinding()
			binding.CertificateID = nil
			certificateCommands := certificate.NewMockedCommands(ctrl)
			delegate := validation.NewValidator()
			bindingValidator := newValidator(delegate, certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
			assert.Error(t, delegate.Result())
		})

		t.Run("HTTPS binding with non-existent certificate fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPSBinding()
			certificateCommands := certCommandsNotExists(ctrl)
			delegate := validation.NewValidator()
			bindingValidator := newValidator(delegate, certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
			assert.Error(t, delegate.Result())
		})

		t.Run("invalid binding type fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			binding := newHTTPBinding()
			binding.Type = "INVALID"
			certificateCommands := certificate.NewMockedCommands(ctrl)
			delegate := validation.NewValidator()
			bindingValidator := newValidator(delegate, certificateCommands)

			err := bindingValidator.validate(t.Context(), "bindings", binding, 0)

			assert.NoError(t, err)
			assert.Error(t, delegate.Result())
		})
	})
}
