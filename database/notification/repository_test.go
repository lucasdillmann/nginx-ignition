package notification

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/notification"
	coreuser "dillmann.com.br/nginx-ignition/core/user"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/common/testutils"
	dbuser "dillmann.com.br/nginx-ignition/database/user"
)

func Test_Repository(t *testing.T) {
	testutils.RunWithMockedDatabases(t, runRepositoryTests)
}

func runRepositoryTests(t *testing.T, db *database.Database) {
	repo := New(db)
	userRepo := dbuser.New(db)

	user := &coreuser.User{
		ID:                   uuid.New(),
		Name:                 "Test User",
		Username:             "testuser-" + uuid.New().String(),
		NotificationLanguage: "en",
		PasswordHash:         "hash",
		PasswordSalt:         "salt",
		Permissions: coreuser.Permissions{
			Hosts:        coreuser.ReadWriteAccessLevel,
			Streams:      coreuser.ReadWriteAccessLevel,
			Certificates: coreuser.ReadWriteAccessLevel,
			Logs:         coreuser.ReadOnlyAccessLevel,
			Integrations: coreuser.ReadWriteAccessLevel,
			AccessLists:  coreuser.ReadWriteAccessLevel,
			Settings:     coreuser.ReadWriteAccessLevel,
			Users:        coreuser.ReadWriteAccessLevel,
			NginxServer:  coreuser.ReadWriteAccessLevel,
			ExportData:   coreuser.ReadOnlyAccessLevel,
			VPNs:         coreuser.ReadWriteAccessLevel,
			Caches:       coreuser.ReadWriteAccessLevel,
			TrafficStats: coreuser.ReadOnlyAccessLevel,
		},
		Enabled: true,
	}
	require.NoError(t, userRepo.Save(t.Context(), user))

	t.Run("SaveNotification", func(t *testing.T) {
		t.Run("persists notification with related entities", func(t *testing.T) {
			certificateID := uuid.New()
			value := newNotification(user.ID)
			relatedEntities := []notification.StoredRelatedEntity{
				{
					NotificationID: value.ID,
					Type:           "CERTIFICATE",
					ID:             certificateID,
					Name:           "example.com",
				},
			}

			require.NoError(t, repo.SaveNotification(t.Context(), value, relatedEntities))

			saved, err := repo.FindNotificationByIDAndUserID(t.Context(), value.ID, user.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, value.Title, saved.Title)
			assert.Equal(t, value.Summary, saved.Summary)
			assert.Equal(t, value.Category, saved.Category)
			assert.Equal(t, value.Payload.Tags, saved.Payload.Tags)

			entities, err := repo.FindRelatedEntitiesByNotificationID(t.Context(), value.ID)
			require.NoError(t, err)
			require.Len(t, entities, 1)
			assert.Equal(t, "CERTIFICATE", entities[0].Type)
			assert.Equal(t, certificateID, entities[0].ID)
		})
	})

	t.Run("FindNotificationPage", func(t *testing.T) {
		t.Run("filters notifications by search terms", func(t *testing.T) {
			matching := newNotification(user.ID)
			matching.ID = uuid.New()
			matching.Title = "Certificate renewed"
			matching.Summary = "example.com was renewed"
			require.NoError(t, repo.SaveNotification(t.Context(), matching, nil))

			other := newNotification(user.ID)
			other.ID = uuid.New()
			other.Title = "Nginx reload failed"
			other.Summary = "reload failed"
			require.NoError(t, repo.SaveNotification(t.Context(), other, nil))

			searchTerm := new("certificate")
			page, err := repo.FindNotificationPage(
				t.Context(),
				user.ID,
				10,
				0,
				searchTerm,
			)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, page.TotalItems, 1)

			for _, item := range page.Contents {
				assert.Contains(
					t,
					strings.ToLower(item.Title)+strings.ToLower(item.Summary),
					*searchTerm,
				)
			}
		})
	})

	t.Run("MarkNotificationAsRead", func(t *testing.T) {
		t.Run("sets read_at for owned notification", func(t *testing.T) {
			value := newNotification(user.ID)
			require.NoError(t, repo.SaveNotification(t.Context(), value, nil))

			require.NoError(t, repo.MarkNotificationAsRead(t.Context(), user.ID, value.ID))

			saved, err := repo.FindNotificationByIDAndUserID(t.Context(), value.ID, user.ID)
			require.NoError(t, err)
			require.NotNil(t, saved.ReadAt)
		})
	})

	t.Run("MarkAllNotificationsAsRead", func(t *testing.T) {
		t.Run("marks every unread notification for user", func(t *testing.T) {
			value := newNotification(user.ID)
			value.ID = uuid.New()
			require.NoError(t, repo.SaveNotification(t.Context(), value, nil))

			require.NoError(t, repo.MarkAllNotificationsAsRead(t.Context(), user.ID))

			count, err := repo.CountUnreadNotifications(t.Context(), user.ID)
			require.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	})

	t.Run("GetLastForUserCategoryAndRelatedEntity", func(t *testing.T) {
		t.Run("returns most recent matching notification", func(t *testing.T) {
			certificateID := uuid.New()
			older := newNotification(user.ID)
			older.ID = uuid.New()
			older.Category = notification.CategoryCertificateExpiring
			older.CreatedAt = time.Now().UTC().Add(-2 * time.Hour)
			require.NoError(
				t,
				repo.SaveNotification(t.Context(), older, []notification.StoredRelatedEntity{
					{NotificationID: older.ID, Type: "CERTIFICATE", ID: certificateID},
				}),
			)

			newer := newNotification(user.ID)
			newer.ID = uuid.New()
			newer.Category = notification.CategoryCertificateExpiring
			newer.CreatedAt = time.Now().UTC()
			require.NoError(
				t,
				repo.SaveNotification(t.Context(), newer, []notification.StoredRelatedEntity{
					{NotificationID: newer.ID, Type: "CERTIFICATE", ID: certificateID},
				}),
			)

			last, err := repo.GetLastForUserCategoryAndRelatedEntity(
				t.Context(),
				user.ID,
				notification.CategoryCertificateExpiring,
				"CERTIFICATE",
				certificateID,
			)
			require.NoError(t, err)
			require.NotNil(t, last)
			assert.Equal(t, newer.ID, last.ID)
		})
	})

	t.Run("SaveConfiguration", func(t *testing.T) {
		t.Run("round-trips categories encoding", func(t *testing.T) {
			value := newConfiguration(user.ID)
			value.Categories = new([]notification.Category{
				notification.CategoryCertificateRenewed,
			})
			require.NoError(t, repo.SaveConfiguration(t.Context(), value))

			saved, err := repo.FindConfigurationByIDAndUserID(t.Context(), value.ID, user.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			require.NotNil(t, saved.Categories)
			assert.Equal(t, *value.Categories, *saved.Categories)
		})

		t.Run("null categories means all categories", func(t *testing.T) {
			value := newConfiguration(user.ID)
			value.ID = uuid.New()
			value.Categories = nil
			require.NoError(t, repo.SaveConfiguration(t.Context(), value))

			saved, err := repo.FindConfigurationByIDAndUserID(t.Context(), value.ID, user.ID)
			require.NoError(t, err)
			assert.Nil(t, saved.Categories)
		})

		t.Run("rejects update for another user configuration", func(t *testing.T) {
			otherUser := &coreuser.User{
				ID:                   uuid.New(),
				Name:                 "Other User",
				Username:             "otheruser-" + uuid.New().String(),
				NotificationLanguage: "en",
				PasswordHash:         "hash",
				PasswordSalt:         "salt",
				Permissions:          user.Permissions,
				Enabled:              true,
			}
			require.NoError(t, userRepo.Save(t.Context(), otherUser))

			value := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), value))

			attempt := newConfiguration(otherUser.ID)
			attempt.ID = value.ID
			attempt.Name = "hijacked"

			err := repo.SaveConfiguration(t.Context(), attempt)
			require.ErrorIs(t, err, notification.ErrConfigurationNotFound)

			saved, findErr := repo.FindConfigurationByIDAndUserID(t.Context(), value.ID, user.ID)
			require.NoError(t, findErr)
			require.NotNil(t, saved)
			assert.Equal(t, value.Name, saved.Name)
		})
	})

	t.Run("SaveProviderSubmissions", func(t *testing.T) {
		t.Run("stores pending submissions for delivery", func(t *testing.T) {
			value := newNotification(user.ID)
			require.NoError(t, repo.SaveNotification(t.Context(), value, nil))

			configuration := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), configuration))

			submission := notification.ProviderSubmission{
				ID:              uuid.New(),
				NotificationID:  value.ID,
				ConfigurationID: configuration.ID,
				Provider:        configuration.Provider,
				Status:          notification.SubmissionStatusPending,
				AttemptCount:    0,
			}
			require.NoError(
				t,
				repo.SaveProviderSubmissions(
					t.Context(),
					[]notification.ProviderSubmission{submission},
				),
			)

			pending, err := repo.FindPendingSubmissionsByNotificationID(t.Context(), value.ID)
			require.NoError(t, err)
			require.Len(t, pending, 1)
			assert.Equal(t, notification.SubmissionStatusPending, pending[0].Status)
		})
	})

	t.Run("FindNotificationsWithIncompleteDelivery", func(t *testing.T) {
		t.Run("returns notifications awaiting delivery completion", func(t *testing.T) {
			value := newNotification(user.ID)
			value.ID = uuid.New()
			value.DeliveryCompleted = false
			require.NoError(t, repo.SaveNotification(t.Context(), value, nil))

			items, err := repo.FindNotificationsWithIncompleteDelivery(t.Context())
			require.NoError(t, err)

			found := false
			for _, item := range items {
				if item.ID == value.ID {
					found = true
					break
				}
			}
			assert.True(t, found)
		})
	})

	t.Run("DeleteConfigurationByIDAndUserID", func(t *testing.T) {
		t.Run("deletes configuration with related submissions", func(t *testing.T) {
			value := newNotification(user.ID)
			require.NoError(t, repo.SaveNotification(t.Context(), value, nil))

			configuration := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), configuration))

			submission := notification.ProviderSubmission{
				ID:              uuid.New(),
				NotificationID:  value.ID,
				ConfigurationID: configuration.ID,
				Provider:        configuration.Provider,
				Status:          notification.SubmissionStatusPending,
				AttemptCount:    0,
			}
			require.NoError(
				t,
				repo.SaveProviderSubmissions(
					t.Context(),
					[]notification.ProviderSubmission{submission},
				),
			)

			require.NoError(
				t,
				repo.DeleteConfigurationByIDAndUserID(t.Context(), configuration.ID, user.ID),
			)

			saved, err := repo.FindConfigurationByIDAndUserID(
				t.Context(),
				configuration.ID,
				user.ID,
			)
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})

	t.Run("ConfigurationExistsByName", func(t *testing.T) {
		t.Run("returns true when another configuration uses the name", func(t *testing.T) {
			value := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), value))

			otherID := uuid.New()
			exists, err := repo.ConfigurationExistsByName(
				t.Context(),
				user.ID,
				value.Name,
				&otherID,
			)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when excluding the same configuration", func(t *testing.T) {
			value := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), value))

			exists, err := repo.ConfigurationExistsByName(
				t.Context(),
				user.ID,
				value.Name,
				&value.ID,
			)
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("SetDeliveryCompleted", func(t *testing.T) {
		t.Run("updates delivery flag", func(t *testing.T) {
			value := newNotification(user.ID)
			value.ID = uuid.New()
			require.NoError(t, repo.SaveNotification(t.Context(), value, nil))

			require.NoError(t, repo.SetDeliveryCompleted(t.Context(), value.ID, true))

			items, err := repo.FindNotificationsWithIncompleteDelivery(t.Context())
			require.NoError(t, err)

			for _, item := range items {
				assert.NotEqual(t, value.ID, item.ID)
			}
		})
	})

	t.Run("CountUnreadNotifications", func(t *testing.T) {
		t.Run("counts only unread notifications for user", func(t *testing.T) {
			unread := newNotification(user.ID)
			unread.ID = uuid.New()
			require.NoError(t, repo.SaveNotification(t.Context(), unread, nil))

			read := newNotification(user.ID)
			read.ID = uuid.New()
			require.NoError(t, repo.SaveNotification(t.Context(), read, nil))
			require.NoError(t, repo.MarkNotificationAsRead(t.Context(), user.ID, read.ID))

			count, err := repo.CountUnreadNotifications(t.Context(), user.ID)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, count, 1)
		})
	})

	t.Run("FindEnabledConfigurationsByUserID", func(t *testing.T) {
		t.Run("returns only enabled configurations", func(t *testing.T) {
			enabled := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), enabled))

			disabled := newConfiguration(user.ID)
			disabled.ID = uuid.New()
			disabled.Enabled = false
			require.NoError(t, repo.SaveConfiguration(t.Context(), disabled))

			configurations, err := repo.FindEnabledConfigurationsByUserID(t.Context(), user.ID)
			require.NoError(t, err)

			for _, configuration := range configurations {
				assert.True(t, configuration.Enabled)
			}
		})
	})

	t.Run("UpdateProviderSubmission", func(t *testing.T) {
		t.Run("persists submission status changes", func(t *testing.T) {
			value := newNotification(user.ID)
			require.NoError(t, repo.SaveNotification(t.Context(), value, nil))

			configuration := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), configuration))

			submission := notification.ProviderSubmission{
				ID:              uuid.New(),
				NotificationID:  value.ID,
				ConfigurationID: configuration.ID,
				Provider:        configuration.Provider,
				Status:          notification.SubmissionStatusPending,
			}
			require.NoError(
				t,
				repo.SaveProviderSubmissions(
					t.Context(),
					[]notification.ProviderSubmission{submission},
				),
			)

			submission.Status = notification.SubmissionStatusSuccess
			require.NoError(t, repo.UpdateProviderSubmission(t.Context(), &submission))

			submissions, err := repo.FindSubmissionsByNotificationID(t.Context(), value.ID)
			require.NoError(t, err)
			require.Len(t, submissions, 1)
			assert.Equal(t, notification.SubmissionStatusSuccess, submissions[0].Status)
		})
	})

	t.Run("FindSubmissionsByNotificationIDs", func(t *testing.T) {
		t.Run("returns submissions grouped by notification", func(t *testing.T) {
			first := newNotification(user.ID)
			first.ID = uuid.New()
			second := newNotification(user.ID)
			second.ID = uuid.New()
			require.NoError(t, repo.SaveNotification(t.Context(), first, nil))
			require.NoError(t, repo.SaveNotification(t.Context(), second, nil))

			configuration := newConfiguration(user.ID)
			require.NoError(t, repo.SaveConfiguration(t.Context(), configuration))

			firstSubmission := notification.ProviderSubmission{
				ID:              uuid.New(),
				NotificationID:  first.ID,
				ConfigurationID: configuration.ID,
				Provider:        configuration.Provider,
				Status:          notification.SubmissionStatusPending,
			}
			secondSubmission := notification.ProviderSubmission{
				ID:              uuid.New(),
				NotificationID:  second.ID,
				ConfigurationID: configuration.ID,
				Provider:        configuration.Provider,
				Status:          notification.SubmissionStatusPending,
			}
			require.NoError(
				t,
				repo.SaveProviderSubmissions(
					t.Context(),
					[]notification.ProviderSubmission{
						firstSubmission,
						secondSubmission,
					},
				),
			)

			grouped, err := repo.FindSubmissionsByNotificationIDs(
				t.Context(),
				[]uuid.UUID{first.ID, second.ID},
			)
			require.NoError(t, err)
			assert.Len(t, grouped[first.ID], 1)
			assert.Len(t, grouped[second.ID], 1)
		})
	})

	t.Run("FindRelatedEntitiesByNotificationIDs", func(t *testing.T) {
		t.Run("returns related entities grouped by notification", func(t *testing.T) {
			certificateID := uuid.New()
			value := newNotification(user.ID)
			require.NoError(
				t,
				repo.SaveNotification(
					t.Context(),
					value,
					[]notification.StoredRelatedEntity{
						{
							NotificationID: value.ID,
							Type:           "CERTIFICATE",
							ID:             certificateID,
							Name:           "example.com",
						},
					},
				),
			)

			grouped, err := repo.FindRelatedEntitiesByNotificationIDs(
				t.Context(),
				[]uuid.UUID{value.ID},
			)
			require.NoError(t, err)
			require.Len(t, grouped[value.ID], 1)
			assert.Equal(t, certificateID, grouped[value.ID][0].ID)
		})
	})
}
