package stream

import (
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

			s := newStream()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(t.Context(), s).Return(nil)

			streamService := newCommands(repo)
			err := streamService.Save(t.Context(), s)

			assert.NoError(t, err)
		})

		t.Run("invalid stream returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := newStream()
			s.Name = ""

			repo := NewMockedRepository(ctrl)
			streamService := newCommands(repo)
			err := streamService.Save(t.Context(), s)

			assert.Error(t, err)
		})

		t.Run("repository error is returned", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := newStream()

			expectedErr := errors.New("repository error")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(t.Context(), s).Return(expectedErr)

			streamService := newCommands(repo)
			err := streamService.Save(t.Context(), s)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().DeleteByID(t.Context(), id).Return(nil)

			streamService := newCommands(repo)
			err := streamService.Delete(t.Context(), id)

			assert.NoError(t, err)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := errors.New("delete failed")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().DeleteByID(t.Context(), id).Return(expectedErr)

			streamService := newCommands(repo)
			err := streamService.Delete(t.Context(), id)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("returns stream when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expected := newStream()
			expected.ID = id

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), id).Return(expected, nil)

			streamService := newCommands(repo)
			result, err := streamService.Get(t.Context(), id)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedPage := pagination.Of([]Stream{*newStream()})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(t.Context(), 10, 1, &searchTerms).Return(expectedPage, nil)

			streamService := newCommands(repo)
			result, err := streamService.List(t.Context(), 10, 1, &searchTerms)

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
			repo.EXPECT().ExistsByID(t.Context(), id).Return(true, nil)

			streamService := newCommands(repo)
			exists, err := streamService.Exists(t.Context(), id)

			assert.NoError(t, err)
			assert.True(t, exists)
		})
	})
}
