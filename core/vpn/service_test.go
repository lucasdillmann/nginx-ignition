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
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), id).Return(&inUse, nil)
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
			inUse := true

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), id).Return(&inUse, nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			err := vpnService.Delete(t.Context(), id)

			require.Error(t, err)
			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.VpnErrorInUse, coreErr.Message.Key)
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
			exists := true

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().ExistsByID(t.Context(), id).Return(&exists, nil)

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
			inUse := false

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(t.Context(), vpn.ID).Return(&inUse, nil)

			cfg := &configuration.Configuration{}
			vpnService := newService(cfg, repo, func() []Driver { return nil })
			err := vpnService.Save(t.Context(), vpn)

			assert.Error(t, err)
		})
	})
}
