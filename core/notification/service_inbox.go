package notification

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

func (s *service) ListNotifications(
	ctx context.Context,
	userID uuid.UUID,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[Notification], error) {
	page, err := s.repository.FindNotificationPage(ctx, userID, pageSize, pageNumber, searchTerms)
	if err != nil {
		return nil, err
	}

	if page == nil || len(page.Contents) == 0 {
		return page, nil
	}

	notificationIDs := make([]uuid.UUID, len(page.Contents))
	for index, item := range page.Contents {
		notificationIDs[index] = item.ID
	}

	relatedEntitiesByNotificationID, err := s.repository.FindRelatedEntitiesByNotificationIDs(
		ctx,
		notificationIDs,
	)
	if err != nil {
		return nil, err
	}

	submissionsByNotificationID, err := s.repository.FindSubmissionsByNotificationIDs(
		ctx,
		notificationIDs,
	)
	if err != nil {
		return nil, err
	}

	for index := range page.Contents {
		notification := &page.Contents[index]
		notification.RelatedEntities = toRelatedEntities(
			relatedEntitiesByNotificationID[notification.ID],
		)
		notification.Submissions = submissionsByNotificationID[notification.ID]
	}

	return page, nil
}

func (s *service) GetNotification(
	ctx context.Context,
	userID uuid.UUID,
	id uuid.UUID,
) (*Notification, error) {
	notification, err := s.repository.FindNotificationByIDAndUserID(ctx, id, userID)
	if err != nil || notification == nil {
		return nil, err
	}

	return s.loadNotificationDetails(ctx, notification)
}

func (s *service) MarkAsRead(ctx context.Context, userID, id uuid.UUID) error {
	return s.repository.MarkNotificationAsRead(ctx, userID, id)
}

func (s *service) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.repository.MarkAllNotificationsAsRead(ctx, userID)
}

func (s *service) UnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repository.CountUnreadNotifications(ctx, userID)
}

func (s *service) GetLastForUserCategoryAndRelatedEntity(
	ctx context.Context,
	userID uuid.UUID,
	category Category,
	entityType string,
	entityID uuid.UUID,
) (*Notification, error) {
	notification, err := s.repository.GetLastForUserCategoryAndRelatedEntity(
		ctx,
		userID,
		category,
		entityType,
		entityID,
	)
	if err != nil || notification == nil {
		return nil, err
	}

	return s.loadNotificationDetails(ctx, notification)
}

func (s *service) loadNotificationDetails(
	ctx context.Context,
	notification *Notification,
) (*Notification, error) {
	relatedEntities, err := s.repository.FindRelatedEntitiesByNotificationID(
		ctx,
		notification.ID,
	)
	if err != nil {
		return nil, err
	}

	notification.RelatedEntities = toRelatedEntities(relatedEntities)
	submissions, err := s.repository.FindSubmissionsByNotificationID(ctx, notification.ID)
	if err != nil {
		return nil, err
	}

	notification.Submissions = submissions
	return notification, nil
}
