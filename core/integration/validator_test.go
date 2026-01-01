package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

func validIntegration() *Integration {
	return &Integration{
		ID:         uuid.New(),
		Name:       "test",
		Driver:     "docker",
		Enabled:    true,
		Parameters: map[string]any{},
	}
}

func TestValidator_Validate(t *testing.T) {
	ctx := context.Background()

	t.Run("valid integration with driver passes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		driverMock := NewMockDriver(ctrl)
		driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
		val := newValidator(repo, driverMock)

		err := val.validate(ctx, integration)

		assert.NoError(t, err)
	})

	t.Run("empty name fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Name = ""
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		val := newValidator(repo, nil)

		err := val.validate(ctx, integration)

		assert.Error(t, err)
	})

	t.Run("whitespace-only name fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Name = "   "
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		val := newValidator(repo, nil)

		err := val.validate(ctx, integration)

		assert.Error(t, err)
	})

	t.Run("empty driver fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Driver = ""
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		val := newValidator(repo, nil)

		err := val.validate(ctx, integration)

		assert.Error(t, err)
	})

	t.Run("whitespace-only driver fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Driver = "   "
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		val := newValidator(repo, nil)

		err := val.validate(ctx, integration)

		assert.Error(t, err)
	})

	t.Run("invalid driver fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Driver = "nonexistent"
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		val := newValidator(repo, nil)

		err := val.validate(ctx, integration)

		assert.Error(t, err)
	})

	t.Run("cannot disable when in use", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Enabled = false
		inUse := true
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		driverMock := NewMockDriver(ctrl)
		driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
		val := newValidator(repo, driverMock)

		err := val.validate(ctx, integration)

		assert.Error(t, err)
	})

	t.Run("can disable when not in use", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Enabled = false
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		driverMock := NewMockDriver(ctrl)
		driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
		val := newValidator(repo, driverMock)

		err := val.validate(ctx, integration)

		assert.NoError(t, err)
	})

	t.Run("dynamicfields validation errors are propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		integration := validIntegration()
		integration.Parameters = map[string]any{
			"requiredField": "",
		}
		inUse := false
		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)
		driverMock := NewMockDriver(ctrl)
		driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{
			{
				ID:       "requiredField",
				Type:     dynamicfields.SingleLineTextType,
				Required: true,
			},
		})
		val := newValidator(repo, driverMock)

		err := val.validate(ctx, integration)

		assert.Error(t, err)
	})
}
