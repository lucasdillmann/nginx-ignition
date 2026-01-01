package vpn

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func TestService_GetByID(t *testing.T) {
	t.Run("returns VPN when found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		expected := &VPN{
			ID:     id,
			Name:   "test",
			Driver: "tailscale",
		}

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

		cfg := &configuration.Configuration{}
		svc := newService(cfg, repo, func() []Driver { return nil })
		result, err := svc.getByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestService_DeleteByID(t *testing.T) {
	t.Run("deletes successfully when not in use", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		inUse := false

		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, id).Return(&inUse, nil)
		repo.EXPECT().DeleteByID(ctx, id).Return(nil)

		cfg := &configuration.Configuration{}
		svc := newService(cfg, repo, func() []Driver { return nil })
		err := svc.deleteByID(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("returns error when in use", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		inUse := true

		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, id).Return(&inUse, nil)

		cfg := &configuration.Configuration{}
		svc := newService(cfg, repo, func() []Driver { return nil })
		err := svc.deleteByID(ctx, id)

		require.Error(t, err)
		var coreErr *coreerror.CoreError
		require.ErrorAs(t, err, &coreErr)
		assert.Contains(t, coreErr.Message, "in use")
	})
}

func TestService_List(t *testing.T) {
	t.Run("returns paginated results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expectedPage := pagination.New(1, 10, 1, []VPN{
			{Name: "test"},
		})
		searchTerms := "test"

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindPage(ctx, 10, 1, &searchTerms, false).Return(expectedPage, nil)

		cfg := &configuration.Configuration{}
		svc := newService(cfg, repo, func() []Driver { return nil })
		result, err := svc.list(ctx, 10, 1, &searchTerms, false)

		assert.NoError(t, err)
		assert.Equal(t, expectedPage, result)
	})
}

func TestService_ExistsByID(t *testing.T) {
	t.Run("returns true when exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		exists := true

		repo := NewMockRepository(ctrl)
		repo.EXPECT().ExistsByID(ctx, id).Return(&exists, nil)

		cfg := &configuration.Configuration{}
		svc := newService(cfg, repo, func() []Driver { return nil })
		result, err := svc.existsByID(ctx, id)

		assert.NoError(t, err)
		assert.True(t, *result)
	})
}

func TestService_Save(t *testing.T) {
	t.Run("invalid VPN returns validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		vpn := &VPN{
			Name: "",
		}
		inUse := false

		repo := NewMockRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, vpn.ID).Return(&inUse, nil)

		cfg := &configuration.Configuration{}
		svc := newService(cfg, repo, func() []Driver { return nil })
		err := svc.save(ctx, vpn)

		assert.Error(t, err)
	})
}
