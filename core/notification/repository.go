package notification

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Repository interface {
	SaveNotification(
		ctx context.Context,
		value *Notification,
		relatedEntities []StoredRelatedEntity,
	) error
	FindNotificationByIDAndUserID(
		ctx context.Context,
		notificationID, userID uuid.UUID,
	) (*Notification, error)
	FindNotificationPage(
		ctx context.Context,
		userID uuid.UUID,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[Notification], error)
	MarkNotificationAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
	MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error
	CountUnreadNotifications(ctx context.Context, userID uuid.UUID) (int, error)
	GetLastForUserCategoryAndRelatedEntity(
		ctx context.Context,
		userID uuid.UUID,
		category Category,
		entityType string,
		entityID uuid.UUID,
	) (*Notification, error)
	SetDeliveryCompleted(ctx context.Context, notificationID uuid.UUID, completed bool) error
	FindRelatedEntitiesByNotificationID(
		ctx context.Context,
		notificationID uuid.UUID,
	) ([]StoredRelatedEntity, error)
	FindRelatedEntitiesByNotificationIDs(
		ctx context.Context,
		notificationIDs []uuid.UUID,
	) (map[uuid.UUID][]StoredRelatedEntity, error)
	FindConfigurationByIDAndUserID(
		ctx context.Context,
		configurationID, userID uuid.UUID,
	) (*Configuration, error)
	FindConfigurationsByUserID(ctx context.Context, userID uuid.UUID) ([]Configuration, error)
	SaveConfiguration(ctx context.Context, value *Configuration) error
	DeleteConfigurationByIDAndUserID(
		ctx context.Context,
		configurationID, userID uuid.UUID,
	) error
	ConfigurationExistsByName(
		ctx context.Context,
		userID uuid.UUID,
		name string,
		excludeID *uuid.UUID,
	) (bool, error)
	FindEnabledConfigurationsByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]Configuration, error)
	SaveProviderSubmissions(ctx context.Context, submissions []ProviderSubmission) error
	FindSubmissionsByNotificationID(
		ctx context.Context,
		notificationID uuid.UUID,
	) ([]ProviderSubmission, error)
	FindSubmissionsByNotificationIDs(
		ctx context.Context,
		notificationIDs []uuid.UUID,
	) (map[uuid.UUID][]ProviderSubmission, error)
	UpdateProviderSubmission(ctx context.Context, value *ProviderSubmission) error
	FindPendingSubmissionsByNotificationID(
		ctx context.Context,
		notificationID uuid.UUID,
	) ([]ProviderSubmission, error)
	FindNotificationsWithIncompleteDelivery(ctx context.Context) ([]Notification, error)
}
