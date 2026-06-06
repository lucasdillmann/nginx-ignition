package notification

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type AvailableProvider struct {
	Name                  *i18n.Message
	ID                    string
	ImportantInstructions []*i18n.Message
	ConfigurationFields   []dynamicfields.DynamicField
}

type Commands interface {
	ListNotifications(
		ctx context.Context,
		userID uuid.UUID,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[Notification], error)
	GetNotification(ctx context.Context, userID, id uuid.UUID) (*Notification, error)
	MarkAsRead(ctx context.Context, userID, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	UnreadCount(ctx context.Context, userID uuid.UUID) (int, error)

	GetLastForUserCategoryAndRelatedEntity(
		ctx context.Context,
		userID uuid.UUID,
		category Category,
		entityType string,
		entityID uuid.UUID,
	) (*Notification, error)

	ListConfigurations(ctx context.Context, userID uuid.UUID) ([]Configuration, error)
	GetConfiguration(ctx context.Context, userID, id uuid.UUID) (*Configuration, error)
	SaveConfiguration(ctx context.Context, userID uuid.UUID, configuration *Configuration) error
	DeleteConfiguration(ctx context.Context, userID, id uuid.UUID) error
	GetAvailableProviders(ctx context.Context) ([]AvailableProvider, error)
	ListCategories(ctx context.Context) ([]CategoryInfo, error)

	Publish(ctx context.Context, userID uuid.UUID, request SendRequest) (*Notification, error)
	Broadcast(ctx context.Context, request SendRequest) error
	ProcessPendingDeliveries(ctx context.Context) error
}
