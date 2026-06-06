package notification

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_configurationAcceptsCategory(t *testing.T) {
	category := CategoryCertificateRenewed

	t.Run("returns true when allowed is nil", func(t *testing.T) {
		assert.True(t, configurationAcceptsCategory(nil, category))
	})

	t.Run("returns false when allowed is empty", func(t *testing.T) {
		assert.False(t, configurationAcceptsCategory(new([]Category{}), category))
	})

	t.Run("returns true when category is in allow-list", func(t *testing.T) {
		assert.True(t, configurationAcceptsCategory(new([]Category{
			CategoryCertificateRenewed,
			CategoryNginxReloadFailed,
		}), category))
	})

	t.Run("returns false when category is not in allow-list", func(t *testing.T) {
		assert.False(t, configurationAcceptsCategory(
			new([]Category{CategoryNginxReloadFailed}),
			category,
		))
	})
}

func Test_service_configuration(t *testing.T) {
	t.Run("SaveConfiguration", func(t *testing.T) {
		t.Run("preserves owner on update", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			configuration := newConfiguration()
			configuration.ID = configurationID
			configuration.UserID = uuid.New()
			configuration.Name = "Updated"

			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(&Configuration{ID: configurationID, UserID: userID}, nil)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "Updated", &configurationID).
				Return(false, nil)
			repository.EXPECT().
				SaveConfiguration(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, saved *Configuration) error {
					assert.Equal(t, userID, saved.UserID)
					return nil
				})

			err := serviceInstance.SaveConfiguration(t.Context(), userID, configuration)
			assert.NoError(t, err)
		})

		t.Run("merges missing sensitive parameters on update", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			configuration := newConfiguration()
			configuration.ID = configurationID
			configuration.Parameters = map[string]any{"host": "smtp.example.com"}

			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(&Configuration{
					ID:     configurationID,
					UserID: userID,
					Parameters: map[string]any{
						"host":     "smtp.example.com",
						"password": "stored-secret",
					},
				}, nil)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, configuration.Name, &configurationID).
				Return(false, nil)
			repository.EXPECT().
				SaveConfiguration(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, saved *Configuration) error {
					assert.Equal(t, "stored-secret", saved.Parameters["password"])
					return nil
				})

			err := serviceInstance.SaveConfiguration(t.Context(), userID, configuration)
			assert.NoError(t, err)
		})

		t.Run("returns not found when configuration belongs to another user", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			configuration := newConfiguration()
			configuration.ID = configurationID
			configuration.UserID = userID
			configuration.Name = "Updated"

			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(nil, nil)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "Updated", &configurationID).
				Return(false, nil)
			repository.EXPECT().
				SaveConfiguration(t.Context(), configuration).
				Return(ErrConfigurationNotFound)

			err := serviceInstance.SaveConfiguration(t.Context(), userID, configuration)
			assert.ErrorIs(t, err, ErrConfigurationNotFound)
		})

		t.Run("assigns user on create", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			configuration := newConfiguration()
			configuration.ID = configurationID
			configuration.Name = "Primary SMTP"

			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(nil, nil)
			repository.EXPECT().
				ConfigurationExistsByName(t.Context(), userID, "Primary SMTP", &configurationID).
				Return(false, nil)
			repository.EXPECT().
				SaveConfiguration(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, saved *Configuration) error {
					assert.Equal(t, userID, saved.UserID)
					return nil
				})

			err := serviceInstance.SaveConfiguration(t.Context(), userID, configuration)
			assert.NoError(t, err)
		})

		t.Run("propagates repository errors", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			expectedErr := errors.New("database error")
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			configuration := newConfiguration()
			configuration.ID = configurationID

			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(nil, expectedErr)

			err := serviceInstance.SaveConfiguration(t.Context(), userID, configuration)
			assert.ErrorIs(t, err, expectedErr)
		})
	})

	t.Run("ListConfigurations", func(t *testing.T) {
		t.Run("removes sensitive parameters", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				FindConfigurationsByUserID(t.Context(), userID).
				Return([]Configuration{
					{
						ID:       configurationID,
						UserID:   userID,
						Provider: "SMTP",
						Parameters: map[string]any{
							"host":     "smtp.example.com",
							"password": "secret",
						},
					},
				}, nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			configurations, err := serviceInstance.ListConfigurations(t.Context(), userID)

			require.NoError(t, err)
			require.Len(t, configurations, 1)
			assert.Equal(t, "smtp.example.com", configurations[0].Parameters["host"])
			assert.NotContains(t, configurations[0].Parameters, "password")
		})
	})

	t.Run("GetConfiguration", func(t *testing.T) {
		t.Run("removes sensitive parameters", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(&Configuration{
					ID:       configurationID,
					UserID:   userID,
					Provider: "SMTP",
					Parameters: map[string]any{
						"host":     "smtp.example.com",
						"password": "secret",
					},
				}, nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			configuration, err := serviceInstance.GetConfiguration(
				t.Context(),
				userID,
				configurationID,
			)

			require.NoError(t, err)
			require.NotNil(t, configuration)
			assert.NotContains(t, configuration.Parameters, "password")
		})

		t.Run("returns nil when configuration does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(nil, nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			configuration, err := serviceInstance.GetConfiguration(
				t.Context(),
				userID,
				configurationID,
			)

			require.NoError(t, err)
			assert.Nil(t, configuration)
		})
	})
}
