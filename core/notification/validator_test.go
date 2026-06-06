package notification

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func Test_validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid configuration passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			err := newValidator(
				repository,
				testProvider{},
			).validate(t.Context(), userID, configuration)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			configuration := newConfiguration()
			configuration.Name = ""

			err := newValidator(NewMockedRepository(ctrl), testProvider{}).
				validate(t.Context(), uuid.New(), configuration)

			assert.Error(t, err)
		})

		t.Run("duplicate name fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			configuration := newConfiguration()

			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, configuration.Name, &configuration.ID).
				Return(true, nil)

			err := newValidator(
				repository,
				testProvider{},
			).validate(t.Context(), userID, configuration)

			assert.Error(t, err)
		})

		t.Run("empty provider fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			configuration.Provider = ""

			err := newValidator(
				repository,
				testProvider{},
			).validate(t.Context(), userID, configuration)

			assert.Error(t, err)
		})

		t.Run("invalid provider fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			configuration.Provider = "UNKNOWN"

			err := newValidator(repository, nil).validate(t.Context(), userID, configuration)

			assert.Error(t, err)
		})

		t.Run("nil categories passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			configuration.Categories = nil

			err := newValidator(
				repository,
				testProvider{},
			).validate(t.Context(), userID, configuration)

			assert.NoError(t, err)
		})

		t.Run("valid categories passes", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			configuration.Categories = new([]Category{
				CategoryCertificateRenewed,
				CategoryNginxReloadFailed,
			})

			err := newValidator(
				repository,
				testProvider{},
			).validate(t.Context(), userID, configuration)

			assert.NoError(t, err)
		})

		t.Run("invalid category fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			configuration.Categories = new([]Category{Category("INVALID")})

			err := newValidator(
				repository,
				testProvider{},
			).validate(t.Context(), userID, configuration)

			assert.Error(t, err)
		})

		t.Run("invalid category reports field violation", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			configuration.Categories = new([]Category{Category("INVALID")})

			err := newValidator(
				repository,
				testProvider{},
			).validate(t.Context(), userID, configuration)

			var consistencyErr *validation.ConsistencyError
			if assert.ErrorAs(t, err, &consistencyErr) {
				require.Len(t, consistencyErr.Violations, 1)
				assert.Equal(t, "categories[0]", consistencyErr.Violations[0].Path)
				assert.Contains(
					t,
					consistencyErr.Violations[0].Message.Key,
					i18n.K.CommonInvalidValue,
				)
			}
		})

		t.Run("missing required parameter fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "test", gomock.Any()).
				Return(false, nil)

			configuration := newConfiguration()
			configuration.Parameters = map[string]any{}

			err := newValidator(repository, requiredHostProvider{}).
				validate(t.Context(), userID, configuration)

			var consistencyErr *validation.ConsistencyError
			if assert.ErrorAs(t, err, &consistencyErr) {
				require.Len(t, consistencyErr.Violations, 1)
				assert.Equal(t, "parameters.host", consistencyErr.Violations[0].Path)
			}
		})
	})
}
