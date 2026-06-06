package notification

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/notification"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) notification.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) SaveNotification(
	ctx context.Context,
	value *notification.Notification,
	relatedEntities []notification.StoredRelatedEntity,
) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	model, err := notificationToModel(value)
	if err != nil {
		return err
	}

	_, err = transaction.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}

	if len(relatedEntities) > 0 {
		entityModels := make([]relatedEntityModel, 0, len(relatedEntities))
		for _, entity := range relatedEntities {
			entityModels = append(entityModels, relatedEntityToModel(entity))
		}

		_, err = transaction.NewInsert().Model(&entityModels).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return transaction.Commit()
}

func (r *repository) FindNotificationByIDAndUserID(
	ctx context.Context,
	notificationID, userID uuid.UUID,
) (*notification.Notification, error) {
	var model notificationModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIDFilter, notificationID).
		Where("user_id = ?", userID).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return notificationToDomain(&model)
}

func (r *repository) FindNotificationPage(
	ctx context.Context,
	userID uuid.UUID,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[notification.Notification], error) {
	models := make([]notificationModel, 0)

	applyFilters := func(query *bun.SelectQuery) *bun.SelectQuery {
		query = query.Where("user_id = ?", userID)

		if searchTerms != nil {
			query = query.Where(
				"LOWER(title) LIKE LOWER(?) OR LOWER(summary) LIKE LOWER(?)",
				"%"+*searchTerms+"%",
				"%"+*searchTerms+"%",
			)
		}

		return query
	}

	count, err := applyFilters(
		r.database.Select().Model((*notificationModel)(nil)),
	).Count(ctx)
	if err != nil {
		return nil, err
	}

	err = applyFilters(
		r.database.Select().Model(&models),
	).
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		OrderExpr("?TableAlias.created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	contents := make([]notification.Notification, 0, len(models))
	for index := range models {
		item, convertErr := notificationToDomain(&models[index])
		if convertErr != nil {
			return nil, convertErr
		}
		contents = append(contents, *item)
	}

	return pagination.New(pageNumber, pageSize, count, contents), nil
}

func (r *repository) MarkNotificationAsRead(
	ctx context.Context,
	userID, notificationID uuid.UUID,
) error {
	now := time.Now()

	_, err := r.database.Update().
		Model((*notificationModel)(nil)).
		Set("read_at = ?", now).
		Where(constants.ByIDFilter, notificationID).
		Where("user_id = ?", userID).
		Where("read_at IS NULL").
		Exec(ctx)

	return err
}

func (r *repository) MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()

	_, err := r.database.Update().
		Model((*notificationModel)(nil)).
		Set("read_at = ?", now).
		Where("user_id = ?", userID).
		Where("read_at IS NULL").
		Exec(ctx)

	return err
}

func (r *repository) CountUnreadNotifications(ctx context.Context, userID uuid.UUID) (int, error) {
	return r.database.Select().
		Model((*notificationModel)(nil)).
		Where("user_id = ?", userID).
		Where("read_at IS NULL").
		Count(ctx)
}

func (r *repository) GetLastForUserCategoryAndRelatedEntity(
	ctx context.Context,
	userID uuid.UUID,
	category notification.Category,
	entityType string,
	entityID uuid.UUID,
) (*notification.Notification, error) {
	var model notificationModel

	err := r.database.Select().
		Model(&model).
		Join("INNER JOIN notification_related_entity AS related_entity ON related_entity.notification_id = ?TableAlias.id").
		Where("user_id = ?", userID).
		Where("category = ?", string(category)).
		Where("related_entity.entity_type = ?", entityType).
		Where("related_entity.entity_id = ?", entityID).
		OrderExpr("?TableAlias.created_at DESC").
		Limit(1).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return notificationToDomain(&model)
}

func (r *repository) SetDeliveryCompleted(
	ctx context.Context,
	notificationID uuid.UUID,
	completed bool,
) error {
	_, err := r.database.Update().
		Model((*notificationModel)(nil)).
		Set("delivery_completed = ?", completed).
		Where(constants.ByIDFilter, notificationID).
		Exec(ctx)

	return err
}

func (r *repository) FindRelatedEntitiesByNotificationID(
	ctx context.Context,
	notificationID uuid.UUID,
) ([]notification.StoredRelatedEntity, error) {
	models := make([]relatedEntityModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("notification_id = ?", notificationID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]notification.StoredRelatedEntity, 0, len(models))
	for index := range models {
		result = append(result, relatedEntityToDomain(&models[index]))
	}

	return result, nil
}

func (r *repository) FindRelatedEntitiesByNotificationIDs(
	ctx context.Context,
	notificationIDs []uuid.UUID,
) (map[uuid.UUID][]notification.StoredRelatedEntity, error) {
	result := make(map[uuid.UUID][]notification.StoredRelatedEntity)
	if len(notificationIDs) == 0 {
		return result, nil
	}

	models := make([]relatedEntityModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("notification_id IN (?)", bun.List(notificationIDs)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	for index := range models {
		entity := relatedEntityToDomain(&models[index])
		result[entity.NotificationID] = append(result[entity.NotificationID], entity)
	}

	return result, nil
}

func (r *repository) FindConfigurationByIDAndUserID(
	ctx context.Context,
	configurationID, userID uuid.UUID,
) (*notification.Configuration, error) {
	var model configurationModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIDFilter, configurationID).
		Where("user_id = ?", userID).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return configurationToDomain(&model)
}

func (r *repository) FindConfigurationsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]notification.Configuration, error) {
	models := make([]configurationModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("user_id = ?", userID).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return configurationsToDomain(models)
}

func (r *repository) SaveConfiguration(
	ctx context.Context,
	value *notification.Configuration,
) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	model, err := configurationToModel(value)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().
		Model((*configurationModel)(nil)).
		Where(constants.ByIDFilter, model.ID).
		Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		result, updateErr := transaction.NewUpdate().
			Model(model).
			Where(constants.ByIDFilter, model.ID).
			Where("user_id = ?", value.UserID).
			Exec(ctx)
		if updateErr != nil {
			return updateErr
		}

		affected, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		if affected == 0 {
			return notification.ErrConfigurationNotFound
		}
	} else {
		_, err = transaction.NewInsert().Model(model).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return transaction.Commit()
}

func (r *repository) DeleteConfigurationByIDAndUserID(
	ctx context.Context,
	configurationID, userID uuid.UUID,
) error {
	_, err := r.database.Delete().
		Model((*configurationModel)(nil)).
		Where(constants.ByIDFilter, configurationID).
		Where("user_id = ?", userID).
		Exec(ctx)

	return err
}

func (r *repository) ConfigurationExistsByName(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	excludeID *uuid.UUID,
) (bool, error) {
	query := r.database.Select().
		Model((*configurationModel)(nil)).
		Where("user_id = ?", userID).
		Where("name = ?", name)

	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	return query.Exists(ctx)
}

func (r *repository) FindEnabledConfigurationsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]notification.Configuration, error) {
	models := make([]configurationModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("user_id = ?", userID).
		Where("enabled = ?", true).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return configurationsToDomain(models)
}

func (r *repository) SaveProviderSubmissions(
	ctx context.Context,
	submissions []notification.ProviderSubmission,
) error {
	if len(submissions) == 0 {
		return nil
	}

	models := make([]providerSubmissionModel, 0, len(submissions))
	for index := range submissions {
		models = append(models, *submissionToModel(&submissions[index]))
	}

	_, err := r.database.Insert().Model(&models).Exec(ctx)
	return err
}

func (r *repository) FindSubmissionsByNotificationID(
	ctx context.Context,
	notificationID uuid.UUID,
) ([]notification.ProviderSubmission, error) {
	models := make([]providerSubmissionModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("notification_id = ?", notificationID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]notification.ProviderSubmission, 0, len(models))
	for index := range models {
		result = append(result, *submissionToDomain(&models[index]))
	}

	return result, nil
}

func (r *repository) FindSubmissionsByNotificationIDs(
	ctx context.Context,
	notificationIDs []uuid.UUID,
) (map[uuid.UUID][]notification.ProviderSubmission, error) {
	result := make(map[uuid.UUID][]notification.ProviderSubmission)
	if len(notificationIDs) == 0 {
		return result, nil
	}

	models := make([]providerSubmissionModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("notification_id IN (?)", bun.List(notificationIDs)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	for index := range models {
		submission := submissionToDomain(&models[index])
		result[submission.NotificationID] = append(result[submission.NotificationID], *submission)
	}

	return result, nil
}

func (r *repository) UpdateProviderSubmission(
	ctx context.Context,
	value *notification.ProviderSubmission,
) error {
	model := submissionToModel(value)

	_, err := r.database.Update().
		Model(model).
		WherePK().
		Exec(ctx)

	return err
}

func (r *repository) FindPendingSubmissionsByNotificationID(
	ctx context.Context,
	notificationID uuid.UUID,
) ([]notification.ProviderSubmission, error) {
	models := make([]providerSubmissionModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("notification_id = ?", notificationID).
		Where("status = ?", string(notification.SubmissionStatusPending)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]notification.ProviderSubmission, 0, len(models))
	for index := range models {
		result = append(result, *submissionToDomain(&models[index]))
	}

	return result, nil
}

func (r *repository) FindNotificationsWithIncompleteDelivery(
	ctx context.Context,
) ([]notification.Notification, error) {
	models := make([]notificationModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("delivery_completed = ?", false).
		Order("created_at").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]notification.Notification, 0, len(models))
	for index := range models {
		item, convertErr := notificationToDomain(&models[index])
		if convertErr != nil {
			return nil, convertErr
		}
		result = append(result, *item)
	}

	return result, nil
}

func configurationsToDomain(models []configurationModel) ([]notification.Configuration, error) {
	result := make([]notification.Configuration, 0, len(models))
	for index := range models {
		item, err := configurationToDomain(&models[index])
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}
	return result, nil
}
