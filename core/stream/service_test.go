package stream

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_service(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("valid stream saves successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			s := newStream()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(ctx, s).Return(nil)

			streamService := newCommands(repo)
			err := streamService.Save(ctx, s)

			assert.NoError(t, err)
		})

		t.Run("invalid stream returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			s := newStream()
			s.Name = ""

			repo := NewMockedRepository(ctrl)
			streamService := newCommands(repo)
			err := streamService.Save(ctx, s)

			assert.Error(t, err)
		})

		t.Run("repository error is returned", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			s := newStream()

			expectedErr := errors.New("repository error")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(ctx, s).Return(expectedErr)

			streamService := newCommands(repo)
			err := streamService.Save(ctx, s)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().DeleteByID(ctx, id).Return(nil)

			streamService := newCommands(repo)
			err := streamService.Delete(ctx, id)

			assert.NoError(t, err)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("delete failed")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().DeleteByID(ctx, id).Return(expectedErr)

			streamService := newCommands(repo)
			err := streamService.Delete(ctx, id)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("returns stream when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expected := newStream()
			expected.ID = id

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

			streamService := newCommands(repo)
			result, err := streamService.Get(ctx, id)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedPage := pagination.Of([]Stream{*newStream()})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(ctx, 10, 1, &searchTerms).Return(expectedPage, nil)

			streamService := newCommands(repo)
			result, err := streamService.List(ctx, 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().ExistsByID(ctx, id).Return(true, nil)

			streamService := newCommands(repo)
			exists, err := streamService.Exists(ctx, id)

			assert.NoError(t, err)
			assert.True(t, exists)
		})
	})
}
