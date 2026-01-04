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

func Test_Service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns user when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()
			expected := &User{
				ID:       id,
				Username: "testuser",
			}

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(ctx, id).Return(expected, nil)

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
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

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			result, err := svc.Get(ctx, id)

			assert.Error(t, err)
			assert.Nil(t, result)
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

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			err := svc.Delete(ctx, id)

			assert.NoError(t, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedPage := pagination.New(1, 10, 1, []User{
				{
					Username: "user1",
				},
			})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(ctx, 10, 1, &searchTerms).Return(expectedPage, nil)

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			result, err := svc.List(ctx, 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})
	})

	t.Run("GetCount", func(t *testing.T) {
		t.Run("returns count", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			expectedCount := 5

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(ctx).Return(expectedCount, nil)

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			count, err := svc.GetCount(ctx)

			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	})

	t.Run("OnboardingCompleted", func(t *testing.T) {
		t.Run("returns true when users exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(ctx).Return(1, nil)

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			completed, err := svc.OnboardingCompleted(ctx)

			assert.NoError(t, err)
			assert.True(t, completed)
		})

		t.Run("returns false when no users exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(ctx).Return(0, nil)

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			completed, err := svc.OnboardingCompleted(ctx)

			assert.NoError(t, err)
			assert.False(t, completed)
		})
	})

	t.Run("Authenticate", func(t *testing.T) {
		t.Run("returns error when user not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(ctx, "nonexistent").Return(nil, nil)

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			result, err := svc.Authenticate(ctx, "nonexistent", "password")

			require.Error(t, err)
			assert.Nil(t, result)
			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Contains(t, coreErr.Message, "Invalid username or password")
		})
	})

	t.Run("GetStatus", func(t *testing.T) {
		t.Run("returns true when enabled", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().IsEnabledByID(ctx, id).Return(true, nil)

			cfg := &configuration.Configuration{}
			svc := &service{
				repository:    repo,
				configuration: cfg,
			}
			enabled, err := svc.GetStatus(ctx, id)

			assert.NoError(t, err)
			assert.True(t, enabled)
		})
	})
}
