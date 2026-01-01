package binding

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

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

func certCommandsExists(certID uuid.UUID) *certificate.Commands {
	return &certificate.Commands{
		Exists: func(_ context.Context, id uuid.UUID) (bool, error) {
			return id == certID, nil
		},
	}
}

func certCommandsNotExists() *certificate.Commands {
	return &certificate.Commands{
		Exists: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return false, nil
		},
	}
}

func TestValidator_Validate(t *testing.T) {
	ctx := context.Background()

	t.Run("valid HTTP binding passes", func(t *testing.T) {
		binding := validHTTPBinding()
		certCommands := &certificate.Commands{}
		val := newValidator(validation.NewValidator(), certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
	})

	t.Run("valid HTTPS binding with certificate passes", func(t *testing.T) {
		binding, certID := validHTTPSBinding()
		certCommands := certCommandsExists(certID)
		val := newValidator(validation.NewValidator(), certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
	})

	t.Run("valid IPv6 address passes", func(t *testing.T) {
		binding := validHTTPBinding()
		binding.IP = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		certCommands := &certificate.Commands{}
		val := newValidator(validation.NewValidator(), certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
	})

	t.Run("invalid IP fails", func(t *testing.T) {
		binding := validHTTPBinding()
		binding.IP = "invalid.ip"
		certCommands := &certificate.Commands{}
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("port below range fails", func(t *testing.T) {
		binding := validHTTPBinding()
		binding.Port = 0
		certCommands := &certificate.Commands{}
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("port above range fails", func(t *testing.T) {
		binding := validHTTPBinding()
		binding.Port = 65536
		certCommands := &certificate.Commands{}
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("HTTP binding with certificate fails", func(t *testing.T) {
		binding := validHTTPBinding()
		certID := uuid.New()
		binding.CertificateID = &certID
		certCommands := &certificate.Commands{}
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("HTTPS binding without certificate fails", func(t *testing.T) {
		binding, _ := validHTTPSBinding()
		binding.CertificateID = nil
		certCommands := &certificate.Commands{}
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("HTTPS binding with non-existent certificate fails", func(t *testing.T) {
		binding, _ := validHTTPSBinding()
		certCommands := certCommandsNotExists()
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})

	t.Run("invalid binding type fails", func(t *testing.T) {
		binding := validHTTPBinding()
		binding.Type = Type("INVALID")
		certCommands := &certificate.Commands{}
		delegate := validation.NewValidator()
		val := newValidator(delegate, certCommands)

		err := val.validate(ctx, "bindings", binding, 0)

		assert.NoError(t, err)
		assert.Error(t, delegate.Result())
	})
}
