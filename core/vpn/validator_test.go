package vpn

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

func validVPN() *VPN {
	return &VPN{
		ID:         uuid.New(),
		Name:       "test",
		Driver:     "tailscale",
		Enabled:    true,
		Parameters: map[string]any{},
	}
}

func Test_Validator(t *testing.T) {
	ctx := context.Background()

	t.Run("Validate", func(t *testing.T) {
		t.Run("valid VPN with driver passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
			val := newValidator(repo, driverMock)

			err := val.validate(ctx, vpn)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Name = ""
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			val := newValidator(repo, nil)

			err := val.validate(ctx, vpn)

			assert.Error(t, err)
		})

		t.Run("whitespace-only name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Name = "   "
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			val := newValidator(repo, nil)

			err := val.validate(ctx, vpn)

			assert.Error(t, err)
		})

		t.Run("empty driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Driver = ""
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			val := newValidator(repo, nil)

			err := val.validate(ctx, vpn)

			assert.Error(t, err)
		})

		t.Run("whitespace-only driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Driver = "   "
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			val := newValidator(repo, nil)

			err := val.validate(ctx, vpn)

			assert.Error(t, err)
		})

		t.Run("invalid driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Driver = "nonexistent"
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			val := newValidator(repo, nil)

			err := val.validate(ctx, vpn)

			assert.Error(t, err)
		})

		t.Run("cannot disable when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Enabled = false
			inUse := true
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
			val := newValidator(repo, driverMock)

			err := val.validate(ctx, vpn)

			assert.Error(t, err)
		})

		t.Run("can disable when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Enabled = false
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})
			val := newValidator(repo, driverMock)

			err := val.validate(ctx, vpn)

			assert.NoError(t, err)
		})

		t.Run("dynamicfields validation errors are propagated", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := validVPN()
			vpn.Parameters = map[string]any{
				"requiredField": "",
			}
			inUse := false
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)
			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{
				{
					ID:       "requiredField",
					Type:     dynamicfields.SingleLineTextType,
					Required: true,
				},
			})
			val := newValidator(repo, driverMock)

			err := val.validate(ctx, vpn)

			assert.Error(t, err)
		})
	})
}
