package cache

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/i18n"
)

func Test_service(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("valid cache saves successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cache := newCache()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().Save(t.Context(), cache).Return(nil)

			cacheService := newCommands(repository)
			err := cacheService.Save(t.Context(), cache)

			assert.NoError(t, err)
		})

		t.Run("invalid cache returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cache := newCache()
			cache.Name = ""

			repository := NewMockedRepository(ctrl)
			cacheService := newCommands(repository)
			err := cacheService.Save(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("repository error is returned", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cache := newCache()
			expectedErr := errors.New("repository error")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().Save(t.Context(), cache).Return(expectedErr)

			cacheService := newCommands(repository)
			err := cacheService.Save(t.Context(), cache)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), id).Return(false, nil)
			repository.EXPECT().DeleteByID(t.Context(), id).Return(nil)

			cacheService := newCommands(repository)
			err := cacheService.Delete(t.Context(), id)

			assert.NoError(t, err)
		})

		t.Run("returns error when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), id).Return(true, nil)

			cacheService := newCommands(repository)
			err := cacheService.Delete(t.Context(), id)

			require.Error(t, err)
			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CacheErrorInUse, coreErr.Message.Key)
		})

		t.Run("returns error when InUseByID fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := errors.New("check failed")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), id).Return(false, expectedErr)

			cacheService := newCommands(repository)
			err := cacheService.Delete(t.Context(), id)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("returns cache when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := newCache()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindByID(t.Context(), expected.ID).Return(expected, nil)

			cacheService := newCommands(repository)
			result, err := cacheService.Get(t.Context(), expected.ID)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := errors.New("not found")

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindByID(t.Context(), id).Return(nil, expectedErr)

			cacheService := newCommands(repository)
			result, err := cacheService.Get(t.Context(), id)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedPage := pagination.Of([]Cache{*newCache()})
			searchTerms := "test"

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindPage(t.Context(), 1, 10, &searchTerms).Return(expectedPage, nil)

			cacheService := newCommands(repository)
			result, err := cacheService.List(t.Context(), 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().ExistsByID(t.Context(), id).Return(true, nil)

			cacheService := newCommands(repository)
			exists, err := cacheService.Exists(t.Context(), id)

			assert.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().ExistsByID(t.Context(), id).Return(false, nil)

			cacheService := newCommands(repository)
			exists, err := cacheService.Exists(t.Context(), id)

			assert.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("GetAllInUse", func(t *testing.T) {
		t.Run("returns all in use caches", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := []Cache{*newCache(), *newCache()}

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindAllInUse(t.Context()).Return(expected, nil)

			cacheService := newCommands(repository)
			result, err := cacheService.GetAllInUse(t.Context())

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})
}
