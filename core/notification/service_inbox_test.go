package notification

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_service_inbox(t *testing.T) {
	t.Run("ListNotifications", func(t *testing.T) {
		t.Run("enriches page with related entities and submissions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			entityID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			page := pagination.New(0, 10, 1, []Notification{{ID: notificationID, UserID: userID}})
			repository.EXPECT().
				FindNotificationPage(t.Context(), userID, 10, 0, (*string)(nil)).
				Return(page, nil)
			repository.EXPECT().
				FindRelatedEntitiesByNotificationIDs(t.Context(), []uuid.UUID{notificationID}).
				Return(map[uuid.UUID][]StoredRelatedEntity{
					notificationID: {{Type: "certificate", ID: entityID, Name: "example.com"}},
				}, nil)
			repository.EXPECT().
				FindSubmissionsByNotificationIDs(t.Context(), []uuid.UUID{notificationID}).
				Return(map[uuid.UUID][]ProviderSubmission{
					notificationID: {{Provider: "SMTP", Status: SubmissionStatusSuccess}},
				}, nil)

			result, err := serviceInstance.ListNotifications(t.Context(), userID, 10, 0, nil)

			require.NoError(t, err)
			require.Len(t, result.Contents, 1)
			require.Len(t, result.Contents[0].RelatedEntities, 1)
			assert.Equal(t, entityID, result.Contents[0].RelatedEntities[0].ID)
			require.Len(t, result.Contents[0].Submissions, 1)
			assert.Equal(t, "SMTP", result.Contents[0].Submissions[0].Provider)
		})

		t.Run("returns empty page without enrichment queries", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			page := pagination.New(0, 10, 0, []Notification{})
			repository.EXPECT().
				FindNotificationPage(t.Context(), userID, 10, 0, (*string)(nil)).
				Return(page, nil)

			result, err := serviceInstance.ListNotifications(t.Context(), userID, 10, 0, nil)

			require.NoError(t, err)
			assert.Empty(t, result.Contents)
		})
	})

	t.Run("GetNotification", func(t *testing.T) {
		t.Run("returns nil when notification does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			repository.EXPECT().
				FindNotificationByIDAndUserID(t.Context(), notificationID, userID).
				Return(nil, nil)

			notification, err := serviceInstance.GetNotification(
				t.Context(),
				userID,
				notificationID,
			)

			require.NoError(t, err)
			assert.Nil(t, notification)
		})

		t.Run("loads related entities and submissions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			entityID := uuid.New()
			repository := NewMockedRepository(ctrl)
			serviceInstance := newService(repository, nil, nil, testProviders)

			expected := &Notification{ID: notificationID, UserID: userID}
			repository.EXPECT().
				FindNotificationByIDAndUserID(t.Context(), notificationID, userID).
				Return(expected, nil)
			repository.EXPECT().
				FindRelatedEntitiesByNotificationID(t.Context(), notificationID).
				Return([]StoredRelatedEntity{{Type: "certificate", ID: entityID}}, nil)
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return([]ProviderSubmission{{Provider: "SMTP"}}, nil)

			notification, err := serviceInstance.GetNotification(
				t.Context(),
				userID,
				notificationID,
			)

			require.NoError(t, err)
			require.Len(t, notification.RelatedEntities, 1)
			require.Len(t, notification.Submissions, 1)
		})
	})

	t.Run("MarkAsRead", func(t *testing.T) {
		t.Run("marks notification as read", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				MarkNotificationAsRead(t.Context(), userID, notificationID).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			err := serviceInstance.MarkAsRead(t.Context(), userID, notificationID)
			assert.NoError(t, err)
		})
	})

	t.Run("MarkAllAsRead", func(t *testing.T) {
		t.Run("marks all notifications as read", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				MarkAllNotificationsAsRead(t.Context(), userID).
				Return(nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			err := serviceInstance.MarkAllAsRead(t.Context(), userID)
			assert.NoError(t, err)
		})
	})

	t.Run("UnreadCount", func(t *testing.T) {
		t.Run("returns unread count", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			repository := NewMockedRepository(ctrl)
			repository.EXPECT().
				CountUnreadNotifications(t.Context(), userID).
				Return(4, nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			count, err := serviceInstance.UnreadCount(t.Context(), userID)

			require.NoError(t, err)
			assert.Equal(t, 4, count)
		})
	})

	t.Run("GetLastForUserCategoryAndRelatedEntity", func(t *testing.T) {
		t.Run("returns notification with related entities", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userID := uuid.New()
			entityID := uuid.New()
			notificationID := uuid.New()
			repository := NewMockedRepository(ctrl)

			expected := &Notification{
				ID:       notificationID,
				UserID:   userID,
				Category: CategoryCertificateExpiring,
			}
			relatedEntities := []StoredRelatedEntity{
				{
					NotificationID: notificationID,
					Type:           "certificate",
					ID:             entityID,
					Name:           "example.com",
				},
			}

			repository.EXPECT().
				GetLastForUserCategoryAndRelatedEntity(
					t.Context(),
					userID,
					CategoryCertificateExpiring,
					"certificate",
					entityID,
				).
				Return(expected, nil)
			repository.EXPECT().
				FindRelatedEntitiesByNotificationID(t.Context(), notificationID).
				Return(relatedEntities, nil)
			repository.EXPECT().
				FindSubmissionsByNotificationID(t.Context(), notificationID).
				Return(nil, nil)

			serviceInstance := newService(repository, nil, nil, testProviders)

			notification, err := serviceInstance.GetLastForUserCategoryAndRelatedEntity(
				t.Context(),
				userID,
				CategoryCertificateExpiring,
				"certificate",
				entityID,
			)

			require.NoError(t, err)
			require.NotNil(t, notification)
			require.Len(t, notification.RelatedEntities, 1)
			assert.Equal(t, entityID, notification.RelatedEntities[0].ID)
		})
	})
}
