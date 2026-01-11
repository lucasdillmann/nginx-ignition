package vpn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

func Test_validator(t *testing.T) {
	t.Run("Validate", func(t *testing.T) {
		t.Run("valid VPN with driver passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})

			vpnValidator := newValidator(repo, driverMock)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Name = ""
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			vpnValidator := newValidator(repo, nil)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.Error(t, err)
		})

		t.Run("whitespace-only name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Name = "   "
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			vpnValidator := newValidator(repo, nil)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.Error(t, err)
		})

		t.Run("empty driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Driver = ""
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			vpnValidator := newValidator(repo, nil)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.Error(t, err)
		})

		t.Run("whitespace-only driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Driver = "   "
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			vpnValidator := newValidator(repo, nil)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.Error(t, err)
		})

		t.Run("invalid driver fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Driver = "nonexistent"
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			vpnValidator := newValidator(repo, nil)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.Error(t, err)
		})

		t.Run("cannot disable when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Enabled = false
			inUse := true

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})

			vpnValidator := newValidator(repo, driverMock)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.Error(t, err)
		})

		t.Run("can disable when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Enabled = false
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{})

			vpnValidator := newValidator(repo, driverMock)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.NoError(t, err)
		})

		t.Run("dynamicfields validation errors are propagated", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Parameters = map[string]any{
				"requiredField": "",
			}
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			driverMock := NewMockedDriver(ctrl)
			driverMock.EXPECT().ConfigurationFields().Return([]dynamicfields.DynamicField{
				{
					ID:       "requiredField",
					Type:     dynamicfields.SingleLineTextType,
					Required: true,
				},
			})

			vpnValidator := newValidator(repo, driverMock)
			err := vpnValidator.validate(t.Context(), vpn)

			assert.Error(t, err)
		})
	})
}
