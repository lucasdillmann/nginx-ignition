package vpn

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_toDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		subject := newVPN()
		result := toDTO(subject)

		assert.NotNil(t, result)
		assert.Equal(t, subject.ID, result.ID)
		assert.Equal(t, subject.Name, result.Name)
		assert.Equal(t, subject.Driver, result.Driver)
		assert.True(t, result.Enabled)
		assert.Equal(t, subject.Parameters, result.Parameters)
	})
}

func Test_doDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		id := uuid.New()
		payload := newVPNRequest()
		result := toDomain(&payload, id)

		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, payload.Name, result.Name)
		assert.Equal(t, payload.Driver, result.Driver)
		assert.True(t, result.Enabled)
		assert.Equal(t, payload.Parameters, result.Parameters)
	})
}
