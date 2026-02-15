package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

func Test_validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid integration with driver passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().
				ConfigurationFields(t.Context()).
				Return([]dynamicfields.DynamicField{})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(t.Context(), integration)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Name = ""
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(t.Context(), integration)

			assert.Error(t, err)
		})

		t.Run("whitespace-only name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Name = "   "
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(t.Context(), integration)

			assert.Error(t, err)
		})

		t.Run("empty driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Driver = ""
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(t.Context(), integration)

			assert.Error(t, err)
		})

		t.Run("whitespace-only driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Driver = "   "
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(t.Context(), integration)

			assert.Error(t, err)
		})

		t.Run("invalid driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Driver = "nonexistent"
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(t.Context(), integration)

			assert.Error(t, err)
		})

		t.Run("cannot disable when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Enabled = false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(true), nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().
				ConfigurationFields(t.Context()).
				Return([]dynamicfields.DynamicField{})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(t.Context(), integration)

			assert.Error(t, err)
		})

		t.Run("can disable when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Enabled = false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().
				ConfigurationFields(t.Context()).
				Return([]dynamicfields.DynamicField{})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(t.Context(), integration)

			assert.NoError(t, err)
		})

		t.Run("dynamicfields validation errors are propagated", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Parameters = map[string]any{
				"requiredField": "",
			}
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(new(false), nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().
				ConfigurationFields(t.Context()).
				Return([]dynamicfields.DynamicField{
					{
						ID:       "requiredField",
						Type:     dynamicfields.SingleLineTextType,
						Required: true,
					},
				})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(t.Context(), integration)

			assert.Error(t, err)
		})
	})
}
