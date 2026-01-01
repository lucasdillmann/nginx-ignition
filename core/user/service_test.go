package user

import (
	"context"
	"errors"
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
	t.Run("returns user when found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		expected := &User{
			ID:       id,
			Username: "testuser",
		}

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		result, err := svc.getByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()
		expectedErr := errors.New("not found")

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindByID(ctx, id).Return(nil, expectedErr)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		result, err := svc.getByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestService_DeleteByID(t *testing.T) {
	t.Run("deletes successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()

		repo := NewMockRepository(ctrl)
		repo.EXPECT().DeleteByID(ctx, id).Return(nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		err := svc.deleteByID(ctx, id)

		assert.NoError(t, err)
	})
}

func TestService_List(t *testing.T) {
	t.Run("returns paginated results", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expectedPage := pagination.New(1, 10, 1, []User{
			{Username: "user1"},
		})
		searchTerms := "test"

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindPage(ctx, 10, 1, &searchTerms).Return(expectedPage, nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		result, err := svc.list(ctx, 10, 1, &searchTerms)

		assert.NoError(t, err)
		assert.Equal(t, expectedPage, result)
	})
}

func TestService_Count(t *testing.T) {
	t.Run("returns count", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		expectedCount := 5

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Count(ctx).Return(expectedCount, nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		count, err := svc.count(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
	})
}

func TestService_IsOnboardingCompleted(t *testing.T) {
	t.Run("returns true when users exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Count(ctx).Return(1, nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		completed, err := svc.isOnboardingCompleted(ctx)

		assert.NoError(t, err)
		assert.True(t, completed)
	})

	t.Run("returns false when no users exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repo := NewMockRepository(ctrl)
		repo.EXPECT().Count(ctx).Return(0, nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		completed, err := svc.isOnboardingCompleted(ctx)

		assert.NoError(t, err)
		assert.False(t, completed)
	})
}

func TestService_Authenticate(t *testing.T) {
	t.Run("returns error when user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repo := NewMockRepository(ctrl)
		repo.EXPECT().FindByUsername(ctx, "nonexistent").Return(nil, nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		result, err := svc.authenticate(ctx, "nonexistent", "password")

		require.Error(t, err)
		assert.Nil(t, result)
		var coreErr *coreerror.CoreError
		require.ErrorAs(t, err, &coreErr)
		assert.Contains(t, coreErr.Message, "Invalid username or password")
	})
}

func TestService_IsEnabled(t *testing.T) {
	t.Run("returns true when enabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		id := uuid.New()

		repo := NewMockRepository(ctrl)
		repo.EXPECT().IsEnabledByID(ctx, id).Return(true, nil)

		cfg := &configuration.Configuration{}
		svc := &service{
			repository:    repo,
			configuration: cfg,
		}
		enabled, err := svc.isEnabled(ctx, id)

		assert.NoError(t, err)
		assert.True(t, enabled)
	})
}
