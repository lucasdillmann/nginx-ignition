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

func Test_service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns integration when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expected := newIntegration()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindByID(ctx, expected.ID).Return(expected, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			result, err := integrationService.Get(ctx, expected.ID)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("deletes successfully when not in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			inUse := false

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, id).Return(&inUse, nil)
			repository.EXPECT().DeleteByID(ctx, id).Return(nil)

			integrationService := newService(repository, func() []Driver { return nil })
			err := integrationService.Delete(ctx, id)

			assert.NoError(t, err)
		})

		t.Run("returns error when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			inUse := true

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, id).Return(&inUse, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			err := integrationService.Delete(ctx, id)

			require.Error(t, err)
			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Contains(t, coreErr.Message, "in use")
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedPage := pagination.Of([]Integration{*newIntegration()})
			searchTerms := "test"

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindPage(ctx, 10, 1, &searchTerms, false).Return(expectedPage, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			result, err := integrationService.List(ctx, 10, 1, &searchTerms, false)

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
			exists := true

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().ExistsByID(ctx, id).Return(&exists, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			result, err := integrationService.Exists(ctx, id)

			assert.NoError(t, err)
			assert.True(t, *result)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("invalid integration returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			integration := newIntegration()
			integration.Name = ""
			inUse := false

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(ctx, integration.ID).Return(&inUse, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			err := integrationService.Save(ctx, integration)

			assert.Error(t, err)
		})
	})
}
