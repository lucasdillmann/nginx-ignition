package notification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/notification"
)

func Test_Converter(t *testing.T) {
	t.Run("notificationToDomain", func(t *testing.T) {
		t.Run("successfully converts a complete model to domain", func(t *testing.T) {
			readAt := time.Now().UTC().Truncate(time.Second)
			createdAt := readAt.Add(-time.Hour)
			payload := notification.Payload{
				OccurredAt: createdAt,
				Tags:       map[string]string{"domain": "example.com"},
				Sections: []notification.DeliverableContentSection{
					{Body: "Renewal completed successfully."},
				},
				Actions: []notification.DeliverableAction{
					{Label: "View certificate", URL: "/certificates"},
				},
			}
			payloadJSON, err := json.Marshal(payload)
			require.NoError(t, err)

			model := &notificationModel{
				ID:                uuid.New(),
				UserID:            uuid.New(),
				Title:             "Certificate renewed",
				Summary:           "Certificate example.com was renewed",
				Category:          string(notification.CategoryCertificateRenewed),
				Payload:           string(payloadJSON),
				ReadAt:            &readAt,
				CreatedAt:         createdAt,
				DeliveryCompleted: true,
			}

			domain, err := notificationToDomain(model)

			require.NoError(t, err)
			assert.Equal(t, model.ID, domain.ID)
			assert.Equal(t, model.UserID, domain.UserID)
			assert.Equal(t, model.Title, domain.Title)
			assert.Equal(t, model.Summary, domain.Summary)
			assert.Equal(t, notification.CategoryCertificateRenewed, domain.Category)
			assert.Equal(t, payload, domain.Payload)
			assert.Equal(t, model.ReadAt, domain.ReadAt)
			assert.Equal(t, model.CreatedAt, domain.CreatedAt)
			assert.True(t, domain.DeliveryCompleted)
		})

		t.Run("returns error when payload is invalid JSON", func(t *testing.T) {
			model := &notificationModel{
				ID:      uuid.New(),
				UserID:  uuid.New(),
				Payload: "{invalid",
			}

			domain, err := notificationToDomain(model)

			assert.Nil(t, domain)
			assert.Error(t, err)
		})
	})

	t.Run("notificationToModel", func(t *testing.T) {
		t.Run("successfully converts a complete domain to model", func(t *testing.T) {
			domain := newNotification(uuid.New())
			domain.ReadAt = new(time.Now().UTC().Truncate(time.Second))
			domain.DeliveryCompleted = true

			model, err := notificationToModel(domain)

			require.NoError(t, err)
			assert.Equal(t, domain.ID, model.ID)
			assert.Equal(t, domain.UserID, model.UserID)
			assert.Equal(t, domain.Title, model.Title)
			assert.Equal(t, domain.Summary, model.Summary)
			assert.Equal(t, string(domain.Category), model.Category)
			assert.Equal(t, domain.ReadAt, model.ReadAt)
			assert.Equal(t, domain.CreatedAt, model.CreatedAt)
			assert.True(t, model.DeliveryCompleted)

			roundTrip, err := notificationToDomain(model)
			require.NoError(t, err)
			assert.Equal(t, domain.Payload, roundTrip.Payload)
		})
	})

	t.Run("configurationToDomain", func(t *testing.T) {
		t.Run("successfully converts a complete model to domain", func(t *testing.T) {
			parametersJSON, err := json.Marshal(map[string]any{
				"host": "smtp.example.com",
				"port": 587,
			})
			require.NoError(t, err)
			categoriesJSON := `["CERTIFICATE_RENEWED","CERTIFICATE_EXPIRING"]`

			model := &configurationModel{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				Name:       "Email alerts",
				Provider:   "SMTP",
				Enabled:    true,
				Parameters: string(parametersJSON),
				Categories: &categoriesJSON,
			}

			domain, err := configurationToDomain(model)

			require.NoError(t, err)
			assert.Equal(t, model.ID, domain.ID)
			assert.Equal(t, model.UserID, domain.UserID)
			assert.Equal(t, model.Name, domain.Name)
			assert.Equal(t, model.Provider, domain.Provider)
			assert.Equal(t, model.Enabled, domain.Enabled)
			assert.Equal(t, "smtp.example.com", domain.Parameters["host"])
			assert.Equal(t, float64(587), domain.Parameters["port"])
			require.NotNil(t, domain.Categories)
			assert.Equal(t, []notification.Category{
				notification.CategoryCertificateRenewed,
				notification.CategoryCertificateExpiring,
			}, *domain.Categories)
		})

		t.Run("maps nil categories to nil", func(t *testing.T) {
			model := &configurationModel{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				Name:       "All categories",
				Provider:   "SMTP",
				Parameters: `{}`,
			}

			domain, err := configurationToDomain(model)

			require.NoError(t, err)
			assert.Nil(t, domain.Categories)
		})

		t.Run("maps empty categories JSON to empty slice", func(t *testing.T) {
			categoriesJSON := `[]`
			model := &configurationModel{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				Name:       "No categories",
				Provider:   "SMTP",
				Parameters: `{}`,
				Categories: &categoriesJSON,
			}

			domain, err := configurationToDomain(model)

			require.NoError(t, err)
			require.NotNil(t, domain.Categories)
			assert.Empty(t, *domain.Categories)
		})

		t.Run("returns error when parameters are invalid JSON", func(t *testing.T) {
			model := &configurationModel{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				Parameters: "{invalid",
			}

			domain, err := configurationToDomain(model)

			assert.Nil(t, domain)
			assert.Error(t, err)
		})

		t.Run("returns error when categories are invalid JSON", func(t *testing.T) {
			categoriesJSON := `{invalid`
			model := &configurationModel{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				Parameters: `{}`,
				Categories: &categoriesJSON,
			}

			domain, err := configurationToDomain(model)

			assert.Nil(t, domain)
			assert.Error(t, err)
		})
	})

	t.Run("configurationToModel", func(t *testing.T) {
		t.Run("successfully converts a complete domain to model", func(t *testing.T) {
			categories := []notification.Category{
				notification.CategoryNginxReloadFailed,
				notification.CategoryNginxReloadSucceeded,
			}
			domain := &notification.Configuration{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				Name:     "Nginx alerts",
				Provider: "SMTP",
				Enabled:  false,
				Parameters: map[string]any{
					"from": "alerts@example.com",
				},
				Categories: &categories,
			}

			model, err := configurationToModel(domain)

			require.NoError(t, err)
			assert.Equal(t, domain.ID, model.ID)
			assert.Equal(t, domain.UserID, model.UserID)
			assert.Equal(t, domain.Name, model.Name)
			assert.Equal(t, domain.Provider, model.Provider)
			assert.False(t, model.Enabled)
			require.NotNil(t, model.Categories)

			roundTrip, err := configurationToDomain(model)
			require.NoError(t, err)
			assert.Equal(t, domain.Parameters["from"], roundTrip.Parameters["from"])
			assert.Equal(t, domain.Categories, roundTrip.Categories)
		})

		t.Run("maps nil categories to nil", func(t *testing.T) {
			domain := newConfiguration(uuid.New())

			model, err := configurationToModel(domain)

			require.NoError(t, err)
			assert.Nil(t, model.Categories)
		})

		t.Run("maps empty categories slice to empty JSON array", func(t *testing.T) {
			categories := []notification.Category{}
			domain := &notification.Configuration{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				Name:       "Disabled categories",
				Provider:   "SMTP",
				Parameters: map[string]any{},
				Categories: &categories,
			}

			model, err := configurationToModel(domain)

			require.NoError(t, err)
			require.NotNil(t, model.Categories)
			assert.Equal(t, `[]`, *model.Categories)
		})
	})

	t.Run("submissionToDomain", func(t *testing.T) {
		t.Run("successfully converts a complete model to domain", func(t *testing.T) {
			lastAttemptAt := time.Now().UTC().Truncate(time.Second)
			succeededAt := lastAttemptAt.Add(time.Minute)
			lastError := "connection refused"

			model := &providerSubmissionModel{
				ID:              uuid.New(),
				NotificationID:  uuid.New(),
				ConfigurationID: uuid.New(),
				Provider:        "SMTP",
				Status:          string(notification.SubmissionStatusSuccess),
				AttemptCount:    2,
				LastError:       &lastError,
				LastAttemptAt:   &lastAttemptAt,
				SucceededAt:     &succeededAt,
			}

			domain := submissionToDomain(model)

			assert.Equal(t, model.ID, domain.ID)
			assert.Equal(t, model.NotificationID, domain.NotificationID)
			assert.Equal(t, model.ConfigurationID, domain.ConfigurationID)
			assert.Equal(t, model.Provider, domain.Provider)
			assert.Equal(t, notification.SubmissionStatusSuccess, domain.Status)
			assert.Equal(t, model.AttemptCount, domain.AttemptCount)
			assert.Equal(t, model.LastError, domain.LastError)
			assert.Equal(t, model.LastAttemptAt, domain.LastAttemptAt)
			assert.Equal(t, model.SucceededAt, domain.SucceededAt)
		})

		t.Run("maps all submission statuses", func(t *testing.T) {
			statuses := []notification.SubmissionStatus{
				notification.SubmissionStatusPending,
				notification.SubmissionStatusSuccess,
				notification.SubmissionStatusFailed,
				notification.SubmissionStatusSkipped,
			}

			for _, status := range statuses {
				model := &providerSubmissionModel{
					ID:     uuid.New(),
					Status: string(status),
				}

				domain := submissionToDomain(model)

				assert.Equal(t, status, domain.Status)
			}
		})
	})

	t.Run("submissionToModel", func(t *testing.T) {
		t.Run("successfully converts a complete domain to model", func(t *testing.T) {
			lastAttemptAt := time.Now().UTC().Truncate(time.Second)
			domain := &notification.ProviderSubmission{
				ID:              uuid.New(),
				NotificationID:  uuid.New(),
				ConfigurationID: uuid.New(),
				Provider:        "SMTP",
				Status:          notification.SubmissionStatusFailed,
				AttemptCount:    3,
				LastError:       new("timeout"),
				LastAttemptAt:   &lastAttemptAt,
			}

			model := submissionToModel(domain)

			assert.Equal(t, domain.ID, model.ID)
			assert.Equal(t, domain.NotificationID, model.NotificationID)
			assert.Equal(t, domain.ConfigurationID, model.ConfigurationID)
			assert.Equal(t, domain.Provider, model.Provider)
			assert.Equal(t, string(domain.Status), model.Status)
			assert.Equal(t, domain.AttemptCount, model.AttemptCount)
			assert.Equal(t, domain.LastError, model.LastError)
			assert.Equal(t, domain.LastAttemptAt, model.LastAttemptAt)
			assert.Nil(t, model.SucceededAt)
		})

		t.Run("round trips submission status through model", func(t *testing.T) {
			domain := &notification.ProviderSubmission{
				ID:     uuid.New(),
				Status: notification.SubmissionStatusSkipped,
			}

			roundTrip := submissionToDomain(submissionToModel(domain))

			assert.Equal(t, domain.Status, roundTrip.Status)
		})
	})

	t.Run("relatedEntityToDomain", func(t *testing.T) {
		t.Run("successfully converts a complete model to domain", func(t *testing.T) {
			notificationID := uuid.New()
			entityID := uuid.New()
			name := "example.com"

			model := &relatedEntityModel{
				NotificationID: notificationID,
				EntityType:     "certificate",
				EntityID:       entityID,
				Name:           &name,
			}

			domain := relatedEntityToDomain(model)

			assert.Equal(t, notificationID, domain.NotificationID)
			assert.Equal(t, "certificate", domain.Type)
			assert.Equal(t, entityID, domain.ID)
			assert.Equal(t, name, domain.Name)
		})

		t.Run("maps nil name to empty string", func(t *testing.T) {
			model := &relatedEntityModel{
				NotificationID: uuid.New(),
				EntityType:     "host",
				EntityID:       uuid.New(),
			}

			domain := relatedEntityToDomain(model)

			assert.Empty(t, domain.Name)
		})
	})

	t.Run("relatedEntityToModel", func(t *testing.T) {
		t.Run("successfully converts a complete domain to model", func(t *testing.T) {
			notificationID := uuid.New()
			entityID := uuid.New()

			domain := notification.StoredRelatedEntity{
				NotificationID: notificationID,
				Type:           "certificate",
				ID:             entityID,
				Name:           "example.com",
			}

			model := relatedEntityToModel(domain)

			assert.Equal(t, notificationID, model.NotificationID)
			assert.Equal(t, domain.Type, model.EntityType)
			assert.Equal(t, entityID, model.EntityID)
			require.NotNil(t, model.Name)
			assert.Equal(t, domain.Name, *model.Name)
		})

		t.Run("maps empty name to nil", func(t *testing.T) {
			domain := notification.StoredRelatedEntity{
				NotificationID: uuid.New(),
				Type:           "host",
				ID:             uuid.New(),
			}

			model := relatedEntityToModel(domain)

			assert.Nil(t, model.Name)
		})

		t.Run("round trips related entity through model", func(t *testing.T) {
			domain := notification.StoredRelatedEntity{
				NotificationID: uuid.New(),
				Type:           "certificate",
				ID:             uuid.New(),
				Name:           "example.com",
			}

			roundTrip := relatedEntityToDomain(new(relatedEntityToModel(domain)))

			assert.Equal(t, domain, roundTrip)
		})
	})
}
