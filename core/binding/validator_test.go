package binding

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func validHTTPBinding() *Binding {
	return &Binding{
		Type: HTTPBindingType,
		IP:   "192.168.1.1",
		Port: 80,
	}
}

func validHTTPSBinding() (*Binding, uuid.UUID) {
	certID := uuid.New()
	return &Binding{
		Type:          HTTPSBindingType,
		IP:            "192.168.1.1",
		Port:          443,
		CertificateID: &certID,
	}, certID
}

func certCommandsExists(ctrl *gomock.Controller, certID uuid.UUID) certificate.Commands {
	m := certificate.NewMockedCommands(ctrl)
	m.EXPECT().Exists(gomock.Any(), certID).AnyTimes().Return(true, nil)
	m.EXPECT().Exists(gomock.Any(), gomock.Not(certID)).AnyTimes().Return(false, nil)
	return m
}

func certCommandsNotExists(ctrl *gomock.Controller) certificate.Commands {
	m := certificate.NewMockedCommands(ctrl)
	m.EXPECT().Exists(gomock.Any(), gomock.Any()).AnyTimes().Return(false, nil)
	return m
}

func Test_Validator_Validate(t *testing.T) {
	ctx := context.Background()

	t.Run("valid HTTP binding passes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding := validHTTPBinding()
		certCommands := certificate.NewMockedCommands(ctrl)
		val := newValidator(validation.NewValidator(), certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
	})

	t.Run("valid HTTPS binding with certificate passes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding, certID := validHTTPSBinding()
		certCommands := certCommandsExists(ctrl, certID)
		val := newValidator(validation.NewValidator(), certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
	})

	t.Run("valid IPv6 address passes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding := validHTTPBinding()
		binding.IP = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		certCommands := certificate.NewMockedCommands(ctrl)
		val := newValidator(validation.NewValidator(), certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
	})

	t.Run("invalid IP fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding := validHTTPBinding()
		binding.IP = "invalid.ip"
		certCommands := certificate.NewMockedCommands(ctrl)
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("port below range fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding := validHTTPBinding()
		binding.Port = 0
		certCommands := certificate.NewMockedCommands(ctrl)
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("port above range fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding := validHTTPBinding()
		binding.Port = 65536
		certCommands := certificate.NewMockedCommands(ctrl)
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("HTTP binding with certificate fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding := validHTTPBinding()
		certID := uuid.New()
		binding.CertificateID = &certID
		certCommands := certificate.NewMockedCommands(ctrl)
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("HTTPS binding without certificate fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding, _ := validHTTPSBinding()
		binding.CertificateID = nil
		certCommands := certificate.NewMockedCommands(ctrl)
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("HTTPS binding with non-existent certificate fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding, _ := validHTTPSBinding()
		certCommands := certCommandsNotExists(ctrl)
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("invalid binding type fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		binding := validHTTPBinding()
		binding.Type = Type("INVALID")
		certCommands := certificate.NewMockedCommands(ctrl)
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})
}
