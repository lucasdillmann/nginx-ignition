package integration

import (
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

			expected := newIntegration()

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().FindByID(t.Context(), expected.ID).Return(expected, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			result, err := integrationService.Get(t.Context(), expected.ID)

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

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), id).Return(&inUse, nil)
			repository.EXPECT().DeleteByID(t.Context(), id).Return(nil)

			integrationService := newService(repository, func() []Driver { return nil })
			err := integrationService.Delete(t.Context(), id)

			assert.NoError(t, err)
		})

		t.Run("returns error when in use", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			inUse := true

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), id).Return(&inUse, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			err := integrationService.Delete(t.Context(), id)

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

			expectedPage := pagination.Of([]Integration{*newIntegration()})
			searchTerms := "test"

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				FindPage(t.Context(), 10, 1, &searchTerms, false).
				Return(expectedPage, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			result, err := integrationService.List(t.Context(), 10, 1, &searchTerms, false)

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

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().ExistsByID(t.Context(), id).Return(&exists, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			result, err := integrationService.Exists(t.Context(), id)

			assert.NoError(t, err)
			assert.True(t, *result)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("invalid integration returns validation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			integration := newIntegration()
			integration.Name = ""
			inUse := false

			repository := NewMockedRepository(ctrl)
			repository.EXPECT().InUseByID(t.Context(), integration.ID).Return(&inUse, nil)

			integrationService := newService(repository, func() []Driver { return nil })
			err := integrationService.Save(t.Context(), integration)

			assert.Error(t, err)
		})
	})
}
