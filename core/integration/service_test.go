package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_Service_GetByID(t *testing.T) {
	t.Run("returns integration when found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		expected := &Integration{
			ID:     id,
			Name:   "test",
			Driver: "docker",
		}

		repo := NewMockedRepository(ctrl)
		repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

		svc := newService(repo, func() []Driver { return nil })
		result, err := svc.Get(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func Test_Service_DeleteByID(t *testing.T) {
	t.Run("deletes successfully when not in use", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		inUse := false

		repo := NewMockedRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, id).Return(&inUse, nil)
		repo.EXPECT().DeleteByID(ctx, id).Return(nil)

		svc := newService(repo, func() []Driver { return nil })
		err := svc.Delete(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("returns error when in use", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		inUse := true

		repo := NewMockedRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, id).Return(&inUse, nil)

		svc := newService(repo, func() []Driver { return nil })
		err := svc.Delete(ctx, id)

		require.Error(t, err)
		var coreErr *coreerror.CoreError
		require.ErrorAs(t, err, &coreErr)
		assert.Contains(t, coreErr.Message, "in use")
	})
}

func Test_Service_List(t *testing.T) {
	t.Run("returns paginated results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expectedPage := pagination.New(1, 10, 1, []Integration{
			{Name: "test"},
		})
		searchTerms := "test"

		repo := NewMockedRepository(ctrl)
		repo.EXPECT().FindPage(ctx, 10, 1, &searchTerms, false).Return(expectedPage, nil)

		svc := newService(repo, func() []Driver { return nil })
		result, err := svc.List(ctx, 10, 1, &searchTerms, false)

		assert.NoError(t, err)
		assert.Equal(t, expectedPage, result)
	})
}

func Test_Service_ExistsByID(t *testing.T) {
	t.Run("returns true when exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		exists := true

		repo := NewMockedRepository(ctrl)
		repo.EXPECT().ExistsByID(ctx, id).Return(&exists, nil)

		svc := newService(repo, func() []Driver { return nil })
		result, err := svc.Exists(ctx, id)

		assert.NoError(t, err)
		assert.True(t, *result)
	})
}

func Test_Service_Save(t *testing.T) {
	t.Run("invalid integration returns validation error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		integration := &Integration{
			Name: "",
		}
		inUse := false

		repo := NewMockedRepository(ctrl)
		repo.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)

		svc := newService(repo, func() []Driver { return nil })
		err := svc.Save(ctx, integration)

		assert.Error(t, err)
	})
}
