package notification

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_resolveSendRequest(t *testing.T) {
	t.Run("resolves all deliverable fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := i18n.NewMockedCommands(ctrl)
		lang := language.AmericanEnglish

		request := SendRequest{
			Title: i18n.DetachedMessage{
				Key:       "title-key",
				Variables: map[string]any{"domain": "example.com"},
			},
			Summary: i18n.DetachedMessage{Key: "summary-key"},
			Sections: []ContentSection{
				{
					Title: new(i18n.DetachedMessage{
						Key:       "section-title-key",
						Variables: map[string]any{"name": "host"},
					}),
					Body: i18n.DetachedMessage{Key: "body-key"},
				},
			},
			Actions: []Action{
				{
					Label: i18n.DetachedMessage{Key: "action-label"},
					URL:   "https://example.com",
				},
			},
			OccurredAt: time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC),
			Tags:       map[string]string{"reminder_days": "30"},
			Category:   CategoryCertificateExpiring,
		}

		commands.EXPECT().
			Translate(lang, "title-key", map[string]any{"domain": "example.com"}).
			Return("Title resolved")
		commands.EXPECT().
			Translate(lang, "summary-key", map[string]any{}).
			Return("Summary resolved")
		commands.EXPECT().
			Translate(lang, "section-title-key", map[string]any{"name": "host"}).
			Return("Section title resolved")
		commands.EXPECT().
			Translate(lang, "body-key", map[string]any{}).
			Return("Body resolved")
		commands.EXPECT().
			Translate(lang, "action-label", map[string]any{}).
			Return("Action label resolved")

		deliverable := resolveSendRequest(commands, lang, request)

		assert.Equal(t, "Title resolved", deliverable.Title)
		assert.Equal(t, "Summary resolved", deliverable.Summary)
		require.Len(t, deliverable.Sections, 1)
		require.NotNil(t, deliverable.Sections[0].Title)
		assert.Equal(t, "Section title resolved", *deliverable.Sections[0].Title)
		assert.Equal(t, "Body resolved", deliverable.Sections[0].Body)
		require.Len(t, deliverable.Actions, 1)
		assert.Equal(t, "Action label resolved", deliverable.Actions[0].Label)
		assert.Equal(t, "https://example.com", deliverable.Actions[0].URL)
		assert.Equal(t, request.OccurredAt, deliverable.OccurredAt)
		assert.Equal(t, request.Tags, deliverable.Tags)
		assert.Equal(t, CategoryCertificateExpiring, deliverable.Category)
	})
}

func Test_resolveLanguage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commands := i18n.NewMockedCommands(ctrl)
	defaultLanguage := language.AmericanEnglish
	commands.EXPECT().DefaultLanguage().Return(defaultLanguage).AnyTimes()
	commands.EXPECT().Supports(gomock.Any()).DoAndReturn(func(tag language.Tag) bool {
		return tag != language.Make("xx")
	}).AnyTimes()

	t.Run("uses default language when value is empty", func(t *testing.T) {
		assert.Equal(t, defaultLanguage, resolveLanguage(commands, ""))
	})

	t.Run("uses parsed language when supported", func(t *testing.T) {
		tag, _ := language.Parse("pt")
		assert.Equal(t, tag, resolveLanguage(commands, "pt"))
	})

	t.Run("falls back to default when unsupported", func(t *testing.T) {
		assert.Equal(t, defaultLanguage, resolveLanguage(commands, "xx"))
	})
}

func Test_service_publish(t *testing.T) {
	t.Run("Publish", func(t *testing.T) {
		t.Run("persists notification and marks delivery completed", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			certificateID := uuid.New()
			repository := NewMockedRepository(ctrl)
			userCommands := user.NewMockedCommands(ctrl)
			commands := i18n.NewMockedCommands(ctrl)
			lang := language.AmericanEnglish

			userCommands.EXPECT().
				Get(t.Context(), userID).
				Return(&user.User{NotificationLanguage: "en"}, nil)
			commands.EXPECT().DefaultLanguage().Return(lang).AnyTimes()
			commands.EXPECT().Supports(gomock.Any()).Return(true).AnyTimes()
			commands.EXPECT().Translate(gomock.Any(), "title-key", gomock.Any()).Return("Title")
			commands.EXPECT().Translate(gomock.Any(), "summary-key", gomock.Any()).Return("Summary")

			repository.EXPECT().
				SaveNotification(t.Context(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, notification *Notification, relatedEntities []StoredRelatedEntity) error {
					assert.Equal(t, userID, notification.UserID)
					assert.Equal(t, "Title", notification.Title)
					assert.Equal(t, "Summary", notification.Summary)
					assert.Equal(t, CategoryCertificateExpiring, notification.Category)
					require.Len(t, relatedEntities, 1)
					assert.Equal(t, "certificate", relatedEntities[0].Type)
					assert.Equal(t, certificateID, relatedEntities[0].ID)
					return nil
				})

			repository.EXPECT().
				FindEnabledConfigurationsByUserID(t.Context(), userID).
				Return([]Configuration{}, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), gomock.Any(), true).
				Return(nil)

			serviceInstance := newService(
				repository,
				userCommands,
				commands,
				func() []Provider { return nil },
			)

			notification, err := serviceInstance.Publish(t.Context(), userID, SendRequest{
				Title:   i18n.DetachedMessage{Key: "title-key"},
				Summary: i18n.DetachedMessage{Key: "summary-key"},
				RelatedEntities: []RelatedEntity{
					{Type: "certificate", ID: certificateID, Name: "example.com"},
				},
				OccurredAt: time.Now(),
				Category:   CategoryCertificateExpiring,
			})

			require.NoError(t, err)
			require.NotNil(t, notification)
			assert.True(t, notification.DeliveryCompleted)
		})

		t.Run("returns error for invalid category", func(t *testing.T) {
			serviceInstance := newService(nil, nil, nil, testProviders)

			notification, err := serviceInstance.Publish(t.Context(), uuid.New(), SendRequest{
				Title:      i18n.DetachedMessage{Key: "title-key"},
				Summary:    i18n.DetachedMessage{Key: "summary-key"},
				OccurredAt: time.Now(),
				Category:   Category("INVALID"),
			})

			assert.Error(t, err)
			assert.Nil(t, notification)
		})

		t.Run("returns error when user does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			userCommands := user.NewMockedCommands(ctrl)
			userCommands.EXPECT().Get(t.Context(), userID).Return(nil, nil)

			serviceInstance := newService(nil, userCommands, nil, testProviders)

			notification, err := serviceInstance.Publish(t.Context(), userID, SendRequest{
				Title:      i18n.DetachedMessage{Key: "title-key"},
				Summary:    i18n.DetachedMessage{Key: "summary-key"},
				OccurredAt: time.Now(),
				Category:   CategoryCertificateExpiring,
			})

			assert.Error(t, err)
			assert.Nil(t, notification)
		})

		t.Run("uses user notification language", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			userCommands := user.NewMockedCommands(ctrl)
			commands := i18n.NewMockedCommands(ctrl)
			portuguese, _ := language.Parse("pt")

			userCommands.EXPECT().
				Get(t.Context(), userID).
				Return(&user.User{NotificationLanguage: "pt"}, nil)
			commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish).AnyTimes()
			commands.EXPECT().Supports(gomock.Any()).Return(true).AnyTimes()
			commands.EXPECT().
				Translate(portuguese, "title-key", gomock.Any()).
				Return("Título")
			commands.EXPECT().
				Translate(portuguese, "summary-key", gomock.Any()).
				Return("Resumo")

			repository.EXPECT().SaveNotification(t.Context(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ any, notification *Notification, _ []StoredRelatedEntity) error {
					assert.Equal(t, "Título", notification.Title)
					assert.Equal(t, "Resumo", notification.Summary)
					return nil
				})
			repository.EXPECT().
				FindEnabledConfigurationsByUserID(t.Context(), userID).
				Return(nil, nil)
			repository.EXPECT().SetDeliveryCompleted(t.Context(), gomock.Any(), true).Return(nil)

			serviceInstance := newService(repository, userCommands, commands, testProviders)

			notification, err := serviceInstance.Publish(t.Context(), userID, SendRequest{
				Title:      i18n.DetachedMessage{Key: "title-key"},
				Summary:    i18n.DetachedMessage{Key: "summary-key"},
				OccurredAt: time.Now(),
				Category:   CategoryCertificateExpiring,
			})

			require.NoError(t, err)
			require.NotNil(t, notification)
		})

		t.Run("creates submissions for matching configurations", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			userCommands := user.NewMockedCommands(ctrl)
			commands := i18n.NewMockedCommands(ctrl)
			lang := language.AmericanEnglish
			userCommands.EXPECT().
				Get(t.Context(), userID).
				Return(&user.User{NotificationLanguage: "en"}, nil)
			commands.EXPECT().DefaultLanguage().Return(lang).AnyTimes()
			commands.EXPECT().Supports(gomock.Any()).Return(true).AnyTimes()
			commands.EXPECT().Translate(gomock.Any(), "title-key", gomock.Any()).Return("Title")
			commands.EXPECT().Translate(gomock.Any(), "summary-key", gomock.Any()).Return("Summary")

			repository.EXPECT().
				SaveNotification(t.Context(), gomock.Any(), gomock.Any()).
				Return(nil)
			repository.EXPECT().
				FindEnabledConfigurationsByUserID(t.Context(), userID).
				Return([]Configuration{
					{
						ID:         configurationID,
						Provider:   "SMTP",
						Enabled:    true,
						Categories: new([]Category{CategoryCertificateExpiring}),
					},
					{
						ID:         uuid.New(),
						Provider:   "SMTP",
						Enabled:    true,
						Categories: new([]Category{CategoryNginxReloadFailed}),
					},
				}, nil)
			repository.EXPECT().
				SaveProviderSubmissions(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, submissions []ProviderSubmission) error {
					require.Len(t, submissions, 1)
					assert.Equal(t, configurationID, submissions[0].ConfigurationID)
					assert.Equal(t, SubmissionStatusPending, submissions[0].Status)
					return nil
				})

			serviceInstance := newService(repository, userCommands, commands, testProviders)

			notification, err := serviceInstance.Publish(t.Context(), userID, SendRequest{
				Title:      i18n.DetachedMessage{Key: "title-key"},
				Summary:    i18n.DetachedMessage{Key: "summary-key"},
				OccurredAt: time.Now(),
				Category:   CategoryCertificateExpiring,
			})

			require.NoError(t, err)
			require.NotNil(t, notification)
			assert.False(t, notification.DeliveryCompleted)
		})
	})

	t.Run("Broadcast", func(t *testing.T) {
		t.Run("publishes to all enabled users", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			firstUserID := uuid.New()
			secondUserID := uuid.New()
			repository := NewMockedRepository(ctrl)
			userCommands := user.NewMockedCommands(ctrl)
			commands := i18n.NewMockedCommands(ctrl)
			lang := language.AmericanEnglish

			userCommands.EXPECT().
				ListEnabledIDs(t.Context()).
				Return([]uuid.UUID{firstUserID, secondUserID}, nil)
			userCommands.EXPECT().
				Get(t.Context(), firstUserID).
				Return(&user.User{NotificationLanguage: "en"}, nil)
			userCommands.EXPECT().
				Get(t.Context(), secondUserID).
				Return(&user.User{NotificationLanguage: "pt"}, nil)

			commands.EXPECT().DefaultLanguage().Return(lang).AnyTimes()
			commands.EXPECT().Supports(gomock.Any()).Return(true).AnyTimes()
			commands.EXPECT().
				Translate(gomock.Any(), "title-key", gomock.Any()).
				Return("Title").
				Times(2)
			commands.EXPECT().
				Translate(gomock.Any(), "summary-key", gomock.Any()).
				Return("Summary").
				Times(2)

			repository.EXPECT().
				SaveNotification(t.Context(), gomock.Any(), gomock.Any()).
				Return(nil).
				Times(2)
			repository.EXPECT().
				FindEnabledConfigurationsByUserID(t.Context(), firstUserID).
				Return(nil, nil)
			repository.EXPECT().
				FindEnabledConfigurationsByUserID(t.Context(), secondUserID).
				Return(nil, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), gomock.Any(), true).
				Return(nil).
				Times(2)

			serviceInstance := newService(
				repository,
				userCommands,
				commands,
				func() []Provider { return nil },
			)

			err := serviceInstance.Broadcast(t.Context(), SendRequest{
				Title:      i18n.DetachedMessage{Key: "title-key"},
				Summary:    i18n.DetachedMessage{Key: "summary-key"},
				OccurredAt: time.Now(),
				Category:   CategoryNginxReloadSucceeded,
			})

			assert.NoError(t, err)
		})
	})
}
