package vpn

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns VPN when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := newVPN()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), expected.ID).Return(expected, nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			result, err := vpnService.Get(t.Context(), expected.ID)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), id).Return(new(false), nil)
			repo.EXPECT().DeleteByID(t.Context(), id).Return(nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			err := vpnService.Delete(t.Context(), id)

			assert.NoError(t, err)
		})

		t.Run("returns error when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), id).Return(new(true), nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			err := vpnService.Delete(t.Context(), id)

			require.Error(t, err)
			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CoreVpnInUse, coreErr.Message.Key)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedPage := pagination.Of([]VPN{*newVPN()})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().
				FindPage(t.Context(), 10, 1, &searchTerms, false).
				Return(expectedPage, nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			result, err := vpnService.List(t.Context(), 10, 1, &searchTerms, false)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().ExistsByID(t.Context(), id).Return(new(true), nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			result, err := vpnService.Exists(t.Context(), id)

			assert.NoError(t, err)
			assert.True(t, *result)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("invalid VPN returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			vpn := newVPN()
			vpn.Name = ""
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(new(false), nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			err := vpnService.Save(t.Context(), vpn)

			assert.Error(t, err)
		})
	})

	t.Run("GetAvailableDrivers", func(t *testing.T) {
		t.Run("returns all drivers sorted by name", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			driver1 := NewMockedDriver(ctrl)
			driver1.EXPECT().ID().Return("b_driver").AnyTimes()
			driver1.EXPECT().Name(gomock.Any()).Return(i18n.Static("B Driver")).AnyTimes()
			driver1.EXPECT().ImportantInstructions(gomock.Any()).Return(nil).AnyTimes()
			driver1.EXPECT().ConfigurationFields(gomock.Any()).Return(nil).AnyTimes()
			driver1.EXPECT().
				EndpointSSLSupport(gomock.Any()).
				Return(DriverManagedEndpointSSLSupport).
				AnyTimes()

			driver2 := NewMockedDriver(ctrl)
			driver2.EXPECT().ID().Return("a_driver").AnyTimes()
			driver2.EXPECT().Name(gomock.Any()).Return(i18n.Static("A Driver")).AnyTimes()
			driver2.EXPECT().ImportantInstructions(gomock.Any()).Return(nil).AnyTimes()
			driver2.EXPECT().ConfigurationFields(gomock.Any()).Return(nil).AnyTimes()
			driver2.EXPECT().
				EndpointSSLSupport(gomock.Any()).
				Return(ProviderManagedEndpointSSLSupport).
				AnyTimes()

			cfg := configuration.New()
			vpnService := newService(
				cfg,
				nil,
				func() []Driver { return []Driver{driver1, driver2} },
			)
			result, err := vpnService.GetAvailableDrivers(t.Context())

			assert.NoError(t, err)
			assert.Len(t, result, 2)

			assert.Equal(t, "a_driver", result[0].ID)
			assert.Equal(t, ProviderManagedEndpointSSLSupport, result[0].EndpointSSLSupport)

			assert.Equal(t, "b_driver", result[1].ID)
			assert.Equal(t, DriverManagedEndpointSSLSupport, result[1].EndpointSSLSupport)
		})
	})

	t.Run("GetAvailableDriverByID", func(t *testing.T) {
		t.Run("returns driver when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			driver1 := NewMockedDriver(ctrl)
			driver1.EXPECT().ID().Return("a_driver").AnyTimes()
			driver1.EXPECT().Name(gomock.Any()).Return(i18n.Static("A Driver")).AnyTimes()
			driver1.EXPECT().ImportantInstructions(gomock.Any()).Return(nil).AnyTimes()
			driver1.EXPECT().ConfigurationFields(gomock.Any()).Return(nil).AnyTimes()
			driver1.EXPECT().
				EndpointSSLSupport(gomock.Any()).
				Return(ProviderManagedEndpointSSLSupport).
				AnyTimes()

			cfg := configuration.New()
			vpnService := newService(cfg, nil, func() []Driver { return []Driver{driver1} })
			result, err := vpnService.GetAvailableDriverByID(t.Context(), "a_driver")

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "a_driver", result.ID)
		})

		t.Run("returns error when driver not found", func(t *testing.T) {
			cfg := configuration.New()
			vpnService := newService(cfg, nil, func() []Driver { return nil })
			result, err := vpnService.GetAvailableDriverByID(t.Context(), "non_existent")

			assert.Nil(t, result)
			require.Error(t, err)

			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CoreVpnDriverNotFound, coreErr.Message.Key)
		})
	})
}
