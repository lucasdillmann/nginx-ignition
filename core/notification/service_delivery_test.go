package notification

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_service_delivery(t *testing.T) {
	t.Run("ProcessPendingDeliveries", func(t *testing.T) {
		t.Run("marks submission skipped when configuration is disabled", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			configurationID := uuid.New()
			submissionID := uuid.New()
			repository := NewMockedRepository(ctrl)

			notification := Notification{
				ID:     notificationID,
				UserID: userID,
			}

			disabledConfiguration := &Configuration{
				ID:      configurationID,
				UserID:  userID,
				Enabled: false,
			}

			repository.EXPECT().
				FindNotificationsWithIncompleteDelivery(t.Context()).
				Return([]Notification{notification}, nil)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{
					{
						ID:              submissionID,
						NotificationID:  notificationID,
						ConfigurationID: configurationID,
						Provider:        "SMTP",
						Status:          SubmissionStatusPending,
					},
				}, nil)
			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(disabledConfiguration, nil)
			repository.EXPECT().
				UpdateProviderSubmission(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, submission *ProviderSubmission) error {
					assert.Equal(t, SubmissionStatusSkipped, submission.Status)
					assert.Nil(t, submission.SucceededAt)
					assert.Nil(t, submission.LastError)
					assert.Equal(t, 0, submission.AttemptCount)
					return nil
				})
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{
					{
						ID:              submissionID,
						ConfigurationID: configurationID,
						Status:          SubmissionStatusSkipped,
					},
				}, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), notificationID, true).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			err := serviceInstance.ProcessPendingDeliveries(t.Context())
			assert.NoError(t, err)
		})

		t.Run("continues processing when one notification fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			firstNotificationID := uuid.New()
			secondNotificationID := uuid.New()
			repository := NewMockedRepository(ctrl)

			repository.EXPECT().
				FindNotificationsWithIncompleteDelivery(t.Context()).
				Return([]Notification{
					{ID: firstNotificationID, UserID: uuid.New()},
					{ID: secondNotificationID, UserID: uuid.New()},
				}, nil)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), firstNotificationID).
				Return(nil, assert.AnError)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), secondNotificationID).
				Return(nil, nil)
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), secondNotificationID).
				Return(nil, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), secondNotificationID, true).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			err := serviceInstance.ProcessPendingDeliveries(t.Context())
			assert.NoError(t, err)
		})

		t.Run("marks submission successful on send", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			configurationID := uuid.New()
			submissionID := uuid.New()
			repository := NewMockedRepository(ctrl)

			repository.EXPECT().
				FindNotificationsWithIncompleteDelivery(t.Context()).
				Return([]Notification{{ID: notificationID, UserID: userID, Title: "Certificate renewed"}}, nil)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{
					{
						ID:              submissionID,
						NotificationID:  notificationID,
						ConfigurationID: configurationID,
						Provider:        "SMTP",
						Status:          SubmissionStatusPending,
					},
				}, nil)
			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(&Configuration{ID: configurationID, UserID: userID, Enabled: true}, nil)
			repository.EXPECT().
				UpdateProviderSubmission(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, submission *ProviderSubmission) error {
					assert.Equal(t, SubmissionStatusSuccess, submission.Status)
					require.NotNil(t, submission.SucceededAt)
					assert.Equal(t, 1, submission.AttemptCount)
					return nil
				})
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{{ID: submissionID, Status: SubmissionStatusSuccess}}, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), notificationID, true).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			err := serviceInstance.ProcessPendingDeliveries(t.Context())
			assert.NoError(t, err)
		})

		t.Run("marks submission failed when provider is not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			submissionID := uuid.New()
			repository := NewMockedRepository(ctrl)

			repository.EXPECT().
				FindNotificationsWithIncompleteDelivery(t.Context()).
				Return([]Notification{{ID: notificationID, UserID: userID}}, nil)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{
					{
						ID:             submissionID,
						NotificationID: notificationID,
						Provider:       "UNKNOWN",
						Status:         SubmissionStatusPending,
					},
				}, nil)
			repository.EXPECT().
				UpdateProviderSubmission(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, submission *ProviderSubmission) error {
					assert.Equal(t, SubmissionStatusFailed, submission.Status)
					require.NotNil(t, submission.LastError)
					assert.Equal(t, "provider not found", *submission.LastError)
					return nil
				})
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{{ID: submissionID, Status: SubmissionStatusFailed}}, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), notificationID, true).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			err := serviceInstance.ProcessPendingDeliveries(t.Context())
			assert.NoError(t, err)
		})

		t.Run("marks submission failed when configuration is not found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			configurationID := uuid.New()
			submissionID := uuid.New()
			repository := NewMockedRepository(ctrl)

			repository.EXPECT().
				FindNotificationsWithIncompleteDelivery(t.Context()).
				Return([]Notification{{ID: notificationID, UserID: userID}}, nil)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{
					{
						ID:              submissionID,
						NotificationID:  notificationID,
						ConfigurationID: configurationID,
						Provider:        "SMTP",
						Status:          SubmissionStatusPending,
					},
				}, nil)
			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(nil, nil)
			repository.EXPECT().
				UpdateProviderSubmission(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, submission *ProviderSubmission) error {
					assert.Equal(t, SubmissionStatusFailed, submission.Status)
					require.NotNil(t, submission.LastError)
					assert.Equal(t, "configuration not found", *submission.LastError)
					return nil
				})
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{{ID: submissionID, Status: SubmissionStatusFailed}}, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), notificationID, true).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			err := serviceInstance.ProcessPendingDeliveries(t.Context())
			assert.NoError(t, err)
		})

		t.Run("retries submission on send error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			configurationID := uuid.New()
			submissionID := uuid.New()
			repository := NewMockedRepository(ctrl)

			repository.EXPECT().
				FindNotificationsWithIncompleteDelivery(t.Context()).
				Return([]Notification{{ID: notificationID, UserID: userID}}, nil)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{
					{
						ID:              submissionID,
						NotificationID:  notificationID,
						ConfigurationID: configurationID,
						Provider:        "SMTP",
						Status:          SubmissionStatusPending,
					},
				}, nil)
			enabledConfiguration := &Configuration{
				ID:      configurationID,
				UserID:  userID,
				Enabled: true,
			}

			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(enabledConfiguration, nil)
			repository.EXPECT().
				UpdateProviderSubmission(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, submission *ProviderSubmission) error {
					assert.Equal(t, SubmissionStatusPending, submission.Status)
					require.NotNil(t, submission.LastError)
					assert.Equal(t, errSendFailed.Error(), *submission.LastError)
					assert.Equal(t, 1, submission.AttemptCount)
					return nil
				})
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{{
					ID:              submissionID,
					ConfigurationID: configurationID,
					Status:          SubmissionStatusPending,
				}}, nil)

			serviceInstance := newService(repository, nil, nil, failingTestProviders(errSendFailed))

			err := serviceInstance.ProcessPendingDeliveries(t.Context())
			assert.NoError(t, err)
		})

		t.Run("marks submission failed at max attempts", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			configurationID := uuid.New()
			submissionID := uuid.New()
			repository := NewMockedRepository(ctrl)

			repository.EXPECT().
				FindNotificationsWithIncompleteDelivery(t.Context()).
				Return([]Notification{{ID: notificationID, UserID: userID}}, nil)
			repository.EXPECT().
				FindPendingSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{
					{
						ID:              submissionID,
						NotificationID:  notificationID,
						ConfigurationID: configurationID,
						Provider:        "SMTP",
						Status:          SubmissionStatusPending,
						AttemptCount:    maxDeliveryAttempts - 1,
					},
				}, nil)
			repository.EXPECT().
				FindConfigurationByIDAndUserID(t.Context(), configurationID, userID).
				Return(&Configuration{ID: configurationID, UserID: userID, Enabled: true}, nil)
			repository.EXPECT().
				UpdateProviderSubmission(t.Context(), gomock.Any()).
				DoAndReturn(func(_ any, submission *ProviderSubmission) error {
					assert.Equal(t, SubmissionStatusFailed, submission.Status)
					assert.Equal(t, maxDeliveryAttempts, submission.AttemptCount)
					return nil
				})
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{{ID: submissionID, Status: SubmissionStatusFailed}}, nil)
			repository.EXPECT().
				SetDeliveryCompleted(t.Context(), notificationID, true).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, failingTestProviders(errSendFailed))

			err := serviceInstance.ProcessPendingDeliveries(t.Context())
			assert.NoError(t, err)
		})
	})
}
