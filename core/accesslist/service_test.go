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

func Test_Service(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("valid access list saves successfully", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			accessList := &AccessList{
				Name:           "test",
				DefaultOutcome: AllowOutcome,
				Entries: []Entry{
					{
						Outcome:       AllowOutcome,
						SourceAddress: []string{"192.168.1.1"},
						Priority:      1,
					},
				},
			}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(ctx, accessList).Return(nil)

			svc := newCommands(repo)
			err := svc.Save(ctx, accessList)

			assert.NoError(t, err)
		})

		t.Run("invalid access list returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			accessList := &AccessList{
				Name: "",
			}

			repo := NewMockedRepository(ctrl)
			svc := newCommands(repo)
			err := svc.Save(ctx, accessList)

			assert.Error(t, err)
		})

		t.Run("repository error is returned", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			accessList := &AccessList{
				Name:           "test",
				DefaultOutcome: AllowOutcome,
				Entries: []Entry{
					{
						Outcome:       AllowOutcome,
						SourceAddress: []string{"192.168.1.1"},
						Priority:      1,
					},
				},
			}

			expectedErr := errors.New("repository error")
			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Save(ctx, accessList).Return(expectedErr)

			svc := newCommands(repo)
			err := svc.Save(ctx, accessList)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, id).Return(false, nil)
			repo.EXPECT().DeleteByID(ctx, id).Return(nil)

			svc := newCommands(repo)
			err := svc.Delete(ctx, id)

			assert.NoError(t, err)
		})

		t.Run("returns error when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, id).Return(true, nil)

			svc := newCommands(repo)
			err := svc.Delete(ctx, id)

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

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, id).Return(false, expectedErr)

			svc := newCommands(repo)
			err := svc.Delete(ctx, id)

			assert.Equal(t, expectedErr, err)
		})

		t.Run("returns error when DeleteByID fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("delete failed")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().InUseByID(ctx, id).Return(false, nil)
			repo.EXPECT().DeleteByID(ctx, id).Return(expectedErr)

			svc := newCommands(repo)
			err := svc.Delete(ctx, id)

			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("returns access list when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expected := &AccessList{
				ID:   id,
				Name: "test",
			}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

			svc := newCommands(repo)
			result, err := svc.Get(ctx, id)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("not found")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(ctx, id).Return(nil, expectedErr)

			svc := newCommands(repo)
			result, err := svc.Get(ctx, id)

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
			expectedPage := pagination.New(1, 10, 1, []AccessList{
				{
					Name: "test",
				},
			})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(ctx, 1, 10, &searchTerms).Return(expectedPage, nil)

			svc := newCommands(repo)
			result, err := svc.List(ctx, 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedErr := errors.New("list failed")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(ctx, 1, 10, (*string)(nil)).Return(nil, expectedErr)

			svc := newCommands(repo)
			result, err := svc.List(ctx, 10, 1, nil)

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

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().ExistsByID(ctx, id).Return(true, nil)

			svc := newCommands(repo)
			exists, err := svc.Exists(ctx, id)

			assert.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().ExistsByID(ctx, id).Return(false, nil)

			svc := newCommands(repo)
			exists, err := svc.Exists(ctx, id)

			assert.NoError(t, err)
			assert.False(t, exists)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expectedErr := errors.New("exists check failed")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().ExistsByID(ctx, id).Return(false, expectedErr)

			svc := newCommands(repo)
			exists, err := svc.Exists(ctx, id)

			assert.Error(t, err)
			assert.False(t, exists)
			assert.Equal(t, expectedErr, err)
		})
	})
}
