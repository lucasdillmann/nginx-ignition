package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_Converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("converts domain object to DTO", func(t *testing.T) {
			id := uuid.New()
			input := &user.User{
				ID:       id,
				Name:     "John Doe",
				Username: "johndoe",
				Enabled:  true,
				Permissions: user.Permissions{
					Hosts: user.ReadWriteAccessLevel,
				},
			}

			result := toDTO(input)

			assert.NotNil(t, result)
			assert.Equal(t, id, result.ID)
			assert.Equal(t, input.Name, result.Name)
			assert.Equal(t, input.Username, result.Username)
			assert.True(t, result.Enabled)
			assert.Equal(t, string(user.ReadWriteAccessLevel), result.Permissions.Hosts)
		})

		t.Run("returns nil when input is nil", func(t *testing.T) {
			result := toDTO(nil)
			assert.Nil(t, result)
		})
	})

	t.Run("toDomain", func(t *testing.T) {
		t.Run("converts DTO to domain object", func(t *testing.T) {
			input := &userRequestDTO{
				Name:     ptr.Of("John Doe"),
				Username: ptr.Of("johndoe"),
				Enabled:  ptr.Of(true),
				Permissions: userPermissionsDTO{
					Hosts: string(user.ReadWriteAccessLevel),
				},
			}

			result := toDomain(input)

			assert.NotNil(t, result)
			assert.Equal(t, *input.Name, result.Name)
			assert.Equal(t, *input.Username, result.Username)
			assert.True(t, result.Enabled)
			assert.Equal(t, user.ReadWriteAccessLevel, result.Permissions.Hosts)
		})
	})
}
