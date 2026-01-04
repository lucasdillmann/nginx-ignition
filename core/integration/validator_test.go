package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

func Test_validator(t *testing.T) {
	ctx := context.Background()

	t.Run("validate", func(t *testing.T) {
		t.Run("valid integration with driver passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(ctx, integration)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Name = ""
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(ctx, integration)

			assert.Error(t, err)
		})

		t.Run("whitespace-only name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Name = "   "
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(ctx, integration)

			assert.Error(t, err)
		})

		t.Run("empty driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Driver = ""
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(ctx, integration)

			assert.Error(t, err)
		})

		t.Run("whitespace-only driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Driver = "   "
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(ctx, integration)

			assert.Error(t, err)
		})

		t.Run("invalid driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Driver = "nonexistent"
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			integrationValidator := newValidator(repository, nil)

			err := integrationValidator.validate(ctx, integration)

			assert.Error(t, err)
		})

		t.Run("cannot disable when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Enabled = false
			inUse := true
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(ctx, integration)

			assert.Error(t, err)
		})

		t.Run("can disable when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Enabled = false
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(ctx, integration)

			assert.NoError(t, err)
		})

		t.Run("dynamicfields validation errors are propagated", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Parameters = map[string]any{
				"requiredField": "",
			}
			inUse := false
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{
				{
					ID:       "requiredField",
					Type:     dynamicfields.SingleLineTextType,
					Required: true,
				},
			})
			integrationValidator := newValidator(repository, driverMock)

			err := integrationValidator.validate(ctx, integration)

			assert.Error(t, err)
		})
	})
}
