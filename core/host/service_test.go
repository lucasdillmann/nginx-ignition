package host

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_service(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("validates input before saving", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			integrationCmds := integration.NewMockedCommands(ctrl)
			vpnCmds := vpn.NewMockedCommands(ctrl)
			aclCmds := accesslist.NewMockedCommands(ctrl)
			cacheCmds := cache.NewMockedCommands(ctrl)
			bindingCmds := binding.NewMockedCommands(ctrl)
			hostService := newCommands(
				repo,
				integrationCmds,
				vpnCmds,
				aclCmds,
				cacheCmds,
				bindingCmds,
			)
			ctx := context.Background()

			input := newHost()
			input.Routes = nil

			repo.EXPECT().FindDefault(ctx).Return(nil, nil).AnyTimes()
			bindingCmds.EXPECT().
				Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)

			err := hostService.Save(ctx, input)
			assert.Error(t, err)
		})

		t.Run("saves valid host", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			integrationCmds := integration.NewMockedCommands(ctrl)
			vpnCmds := vpn.NewMockedCommands(ctrl)
			aclCmds := accesslist.NewMockedCommands(ctrl)
			cacheCmds := cache.NewMockedCommands(ctrl)
			bindingCmds := binding.NewMockedCommands(ctrl)
			hostService := newCommands(
				repo,
				integrationCmds,
				vpnCmds,
				aclCmds,
				cacheCmds,
				bindingCmds,
			)
			ctx := context.Background()

			input := newHost()

			// Mocks for validation
			bindingCmds.EXPECT().
				Validate(ctx, "bindings", 0, &input.Bindings[0], gomock.Any()).
				Return(nil)

			repo.EXPECT().Save(ctx, input).Return(nil)

			err := hostService.Save(ctx, input)
			assert.NoError(t, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("delegates to repository", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			hostService := newCommands(repo, nil, nil, nil, nil, nil)
			ctx := context.Background()
			id := uuid.New()

			repo.EXPECT().DeleteByID(ctx, id).Return(nil)

			err := hostService.Delete(ctx, id)
			assert.NoError(t, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("delegates to repository", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			hostService := newCommands(repo, nil, nil, nil, nil, nil)
			ctx := context.Background()
			pageSize := 10
			pageNumber := 1
			search := ptr.Of("term")

			expectedPage := &pagination.Page[Host]{}
			repo.EXPECT().FindPage(ctx, pageSize, pageNumber, search).Return(expectedPage, nil)

			page, err := hostService.List(ctx, pageSize, pageNumber, search)
			assert.NoError(t, err)
			assert.Equal(t, expectedPage, page)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("delegates to repository", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			hostService := newCommands(repo, nil, nil, nil, nil, nil)
			ctx := context.Background()
			id := uuid.New()

			expectedHost := newHost()
			repo.EXPECT().FindByID(ctx, id).Return(expectedHost, nil)

			result, err := hostService.Get(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, expectedHost, result)
		})
	})

	t.Run("GetAllEnabled", func(t *testing.T) {
		t.Run("delegates to repository", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			hostService := newCommands(repo, nil, nil, nil, nil, nil)
			ctx := context.Background()

			expectedHosts := []Host{*newHost()}
			repo.EXPECT().FindAllEnabled(ctx).Return(expectedHosts, nil)

			result, err := hostService.GetAllEnabled(ctx)
			assert.NoError(t, err)
			assert.Equal(t, expectedHosts, result)
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("delegates to repository", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			hostService := newCommands(repo, nil, nil, nil, nil, nil)
			ctx := context.Background()
			id := uuid.New()

			repo.EXPECT().ExistsByID(ctx, id).Return(true, nil)

			exists, err := hostService.Exists(ctx, id)
			assert.NoError(t, err)
			assert.True(t, exists)
		})
	})
}
