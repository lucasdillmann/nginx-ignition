package authorization

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
)

const (
	testJwtSecret        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	testAnotherJwtSecret = "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
)

func newServiceWithMockedRepository(
	t *testing.T,
	overrides map[string]string,
) (*service, *MockedRepository) {
	t.Helper()

	controller := gomock.NewController(t)
	repository := NewMockedRepository(controller)
	return newService(repository, configuration.NewWithOverrides(overrides)), repository
}

func Test_service_JwtSecret(t *testing.T) {
	t.Run("returns an empty value before the initialization", func(t *testing.T) {
		svc, _ := newServiceWithMockedRepository(t, map[string]string{})

		assert.Empty(t, svc.JwtSecret())
	})
}

func Test_service_initialize(t *testing.T) {
	t.Run("uses the secret stored in the database", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{})
		stored := testJwtSecret
		repository.EXPECT().FindJwtSecret(t.Context()).Return(&stored, nil)

		err := svc.initialize(t.Context())

		require.NoError(t, err)
		assert.Equal(t, testJwtSecret, svc.JwtSecret())
	})

	t.Run("ignores a configured secret that differs from the stored one", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{
			"nginx-ignition.security.jwt.secret": testAnotherJwtSecret,
		})
		stored := testJwtSecret
		repository.EXPECT().FindJwtSecret(t.Context()).Return(&stored, nil)

		err := svc.initialize(t.Context())

		require.NoError(t, err)
		assert.Equal(t, testJwtSecret, svc.JwtSecret())
	})

	t.Run("returns an error for an invalid stored secret length", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{})
		stored := "too-short"
		repository.EXPECT().FindJwtSecret(t.Context()).Return(&stored, nil)

		err := svc.initialize(t.Context())

		assert.EqualError(
			t,
			err,
			"JWT secret stored in the database should be 64 characters long but is 9 characters long",
		)
		assert.Empty(t, svc.JwtSecret())
	})

	t.Run("stores the configured secret when the table is empty", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{
			"nginx-ignition.security.jwt.secret": testJwtSecret,
		})
		repository.EXPECT().FindJwtSecret(t.Context()).Return(nil, nil)

		var stored string
		repository.EXPECT().
			SaveJwtSecret(t.Context(), gomock.Any()).
			DoAndReturn(func(_ context.Context, secret *string) error {
				stored = *secret
				return nil
			})

		err := svc.initialize(t.Context())

		require.NoError(t, err)
		assert.Equal(t, testJwtSecret, stored)
		assert.Equal(t, testJwtSecret, svc.JwtSecret())
	})

	t.Run("returns an error for an invalid configured secret length", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{
			"nginx-ignition.security.jwt.secret": "too-short",
		})
		repository.EXPECT().FindJwtSecret(t.Context()).Return(nil, nil)

		err := svc.initialize(t.Context())

		assert.EqualError(
			t,
			err,
			"a custom JWT secret should be 64 characters long, but the provided value is 9 characters long",
		)
		assert.Empty(t, svc.JwtSecret())
	})

	t.Run("generates and stores a random secret when none is configured", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{})
		repository.EXPECT().FindJwtSecret(t.Context()).Return(nil, nil)

		var stored string
		repository.EXPECT().
			SaveJwtSecret(t.Context(), gomock.Any()).
			DoAndReturn(func(_ context.Context, secret *string) error {
				stored = *secret
				return nil
			})

		err := svc.initialize(t.Context())

		require.NoError(t, err)
		assert.Len(t, stored, 64)
		assert.Equal(t, stored, svc.JwtSecret())
	})

	t.Run("returns an error when the database lookup fails", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{})
		repository.EXPECT().FindJwtSecret(t.Context()).Return(nil, assert.AnError)

		err := svc.initialize(t.Context())

		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, svc.JwtSecret())
	})

	t.Run("returns an error when the secret can't be stored", func(t *testing.T) {
		svc, repository := newServiceWithMockedRepository(t, map[string]string{})
		repository.EXPECT().FindJwtSecret(t.Context()).Return(nil, nil)
		repository.EXPECT().
			SaveJwtSecret(t.Context(), gomock.Any()).
			Return(assert.AnError)

		err := svc.initialize(t.Context())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
