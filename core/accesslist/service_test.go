package accesslist

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_service(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("valid access list saves successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			accessList := newAccessList()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().Save(ctx, accessList).Return(nil)

			accessListService := newCommands(repository)
			err := accessListService.Save(ctx, accessList)

			assert.NoError(t, err)
		})

		t.Run("invalid access list returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			accessList := newAccessList()
			accessList.Name = ""

			repository := NewMockedRepository(ctrl)
			accessListService := newCommands(repository)
			err := accessListService.Save(ctx, accessList)

			assert.Error(t, err)
		})

		t.Run("repository error is returned", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			accessList := newAccessList()
			expectedErr := errors.New("repository error")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().Save(ctx, accessList).Return(expectedErr)

			accessListService := newCommands(repository)
			err := accessListService.Save(ctx, accessList)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, id).Return(false, nil)
			repository.EXPECT().DeleteByID(ctx, id).Return(nil)

			accessListService := newCommands(repository)
			err := accessListService.Delete(ctx, id)

			assert.NoError(t, err)
		})

		t.Run("returns error when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, id).Return(true, nil)

			accessListService := newCommands(repository)
			err := accessListService.Delete(ctx, id)

			require.Error(t, err)
			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Contains(t, coreErr.Message, "in use")
		})

		t.Run("returns error when InUseByID fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("check failed")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, id).Return(false, expectedErr)

			accessListService := newCommands(repository)
			err := accessListService.Delete(ctx, id)

			assert.Equal(t, expectedErr, err)
		})

		t.Run("returns error when DeleteByID fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("delete failed")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, id).Return(false, nil)
			repository.EXPECT().DeleteByID(ctx, id).Return(expectedErr)

			accessListService := newCommands(repository)
			err := accessListService.Delete(ctx, id)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("returns access list when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expected := newAccessList()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindByID(ctx, expected.ID).Return(expected, nil)

			accessListService := newCommands(repository)
			result, err := accessListService.Get(ctx, expected.ID)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("not found")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindByID(ctx, id).Return(nil, expectedErr)

			accessListService := newCommands(repository)
			result, err := accessListService.Get(ctx, id)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedPage := pagination.Of([]AccessList{*newAccessList()})
			searchTerms := "test"

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindPage(ctx, 1, 10, &searchTerms).Return(expectedPage, nil)

			accessListService := newCommands(repository)
			result, err := accessListService.List(ctx, 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedErr := errors.New("list failed")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindPage(ctx, 1, 10, (*string)(nil)).Return(nil, expectedErr)

			accessListService := newCommands(repository)
			result, err := accessListService.List(ctx, 10, 1, nil)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().ExistsByID(ctx, id).Return(true, nil)

			accessListService := newCommands(repository)
			exists, err := accessListService.Exists(ctx, id)

			assert.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().ExistsByID(ctx, id).Return(false, nil)

			accessListService := newCommands(repository)
			exists, err := accessListService.Exists(ctx, id)

			assert.NoError(t, err)
			assert.False(t, exists)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("exists check failed")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().ExistsByID(ctx, id).Return(false, expectedErr)

			accessListService := newCommands(repository)
			exists, err := accessListService.Exists(ctx, id)

			assert.Error(t, err)
			assert.False(t, exists)
			assert.Equal(t, expectedErr, err)
		})
	})
}
