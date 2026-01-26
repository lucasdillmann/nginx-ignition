package user

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_service(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns user when found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := newUser()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), expected.ID).Return(expected, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			result, err := svc.Get(t.Context(), expected.ID)

			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("returns error when repository fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := errors.New("not found")

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByID(t.Context(), id).Return(nil, expectedErr)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			result, err := svc.Get(t.Context(), id)

			assert.Error(t, err)
			assert.Nil(t, result)
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

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			err := svc.Delete(t.Context(), id)

			assert.NoError(t, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("returns paginated results", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedPage := pagination.Of([]User{*newUser()})
			searchTerms := "test"

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindPage(t.Context(), 10, 1, &searchTerms).Return(expectedPage, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			result, err := svc.List(t.Context(), 10, 1, &searchTerms)

			assert.NoError(t, err)
			assert.Equal(t, expectedPage, result)
		})
	})

	t.Run("GetCount", func(t *testing.T) {
		t.Run("returns count", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedCount := 5

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(t.Context()).Return(expectedCount, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			count, err := svc.GetCount(t.Context())

			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	})

	t.Run("OnboardingCompleted", func(t *testing.T) {
		t.Run("returns true when users exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(t.Context()).Return(1, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			completed, err := svc.OnboardingCompleted(t.Context())

			assert.NoError(t, err)
			assert.True(t, completed)
		})

		t.Run("returns false when no users exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().Count(t.Context()).Return(0, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			completed, err := svc.OnboardingCompleted(t.Context())

			assert.NoError(t, err)
			assert.False(t, completed)
		})
	})

	t.Run("Authenticate", func(t *testing.T) {
		t.Run("returns error when user not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().FindByUsername(t.Context(), "nonexistent").Return(nil, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			result, err := svc.Authenticate(t.Context(), "nonexistent", "password")

			require.Error(t, err)
			assert.Nil(t, result)
			var coreErr *coreerror.CoreError
			require.ErrorAs(t, err, &coreErr)
			assert.Equal(t, i18n.K.CoreUserInvalidCredentials, coreErr.Message.Key)
		})
	})

	t.Run("GetStatus", func(t *testing.T) {
		t.Run("returns true when enabled", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()

			repo := NewMockedRepository(ctrl)
			repo.EXPECT().IsEnabledByID(t.Context(), id).Return(true, nil)

			cfg := &configuration.Configuration{}
			svc, _ := newCommands(repo, cfg)
			enabled, err := svc.GetStatus(t.Context(), id)

			assert.NoError(t, err)
			assert.True(t, enabled)
		})
	})
}
