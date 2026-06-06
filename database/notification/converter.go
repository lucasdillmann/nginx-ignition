package notification

import (
	"encoding/json"

	"dillmann.com.br/nginx-ignition/core/notification"
)

func notificationToDomain(model *notificationModel) (*notification.Notification, error) {
	payload := notification.Payload{}
	err := json.Unmarshal([]byte(model.Payload), &payload)
	if err != nil {
		return nil, err
	}

	return &notification.Notification{
		ID:                model.ID,
		UserID:            model.UserID,
		Title:             model.Title,
		Summary:           model.Summary,
		Category:          notification.Category(model.Category),
		Payload:           payload,
		ReadAt:            model.ReadAt,
		CreatedAt:         model.CreatedAt,
		DeliveryCompleted: model.DeliveryCompleted,
	}, nil
}

func notificationToModel(domain *notification.Notification) (*notificationModel, error) {
	payload, err := json.Marshal(domain.Payload)
	if err != nil {
		return nil, err
	}

	return &notificationModel{
		ID:                domain.ID,
		UserID:            domain.UserID,
		Title:             domain.Title,
		Summary:           domain.Summary,
		Category:          string(domain.Category),
		Payload:           string(payload),
		ReadAt:            domain.ReadAt,
		CreatedAt:         domain.CreatedAt,
		DeliveryCompleted: domain.DeliveryCompleted,
	}, nil
}

func configurationToDomain(model *configurationModel) (*notification.Configuration, error) {
	parameters := make(map[string]any)
	err := json.Unmarshal([]byte(model.Parameters), &parameters)
	if err != nil {
		return nil, err
	}

	categories, err := categoriesFromDatabase(model.Categories)
	if err != nil {
		return nil, err
	}

	return &notification.Configuration{
		ID:         model.ID,
		UserID:     model.UserID,
		Name:       model.Name,
		Provider:   model.Provider,
		Enabled:    model.Enabled,
		Parameters: parameters,
		Categories: categories,
	}, nil
}

func configurationToModel(domain *notification.Configuration) (*configurationModel, error) {
	parameters, err := json.Marshal(domain.Parameters)
	if err != nil {
		return nil, err
	}

	categories, err := categoriesToDatabase(domain.Categories)
	if err != nil {
		return nil, err
	}

	return &configurationModel{
		ID:         domain.ID,
		UserID:     domain.UserID,
		Name:       domain.Name,
		Provider:   domain.Provider,
		Enabled:    domain.Enabled,
		Parameters: string(parameters),
		Categories: categories,
	}, nil
}

func categoriesToDatabase(categories *[]notification.Category) (*string, error) {
	if categories == nil {
		return nil, nil
	}

	encoded, err := json.Marshal(categories)
	if err != nil {
		return nil, err
	}

	return new(string(encoded)), nil
}

func categoriesFromDatabase(value *string) (*[]notification.Category, error) {
	if value == nil {
		return nil, nil
	}

	categories := make([]notification.Category, 0)
	err := json.Unmarshal([]byte(*value), &categories)
	if err != nil {
		return nil, err
	}

	return new(categories), nil
}

func submissionToDomain(model *providerSubmissionModel) *notification.ProviderSubmission {
	return &notification.ProviderSubmission{
		ID:              model.ID,
		NotificationID:  model.NotificationID,
		ConfigurationID: model.ConfigurationID,
		Provider:        model.Provider,
		Status:          notification.SubmissionStatus(model.Status),
		AttemptCount:    model.AttemptCount,
		LastError:       model.LastError,
		LastAttemptAt:   model.LastAttemptAt,
		SucceededAt:     model.SucceededAt,
	}
}

func submissionToModel(domain *notification.ProviderSubmission) *providerSubmissionModel {
	return &providerSubmissionModel{
		ID:              domain.ID,
		NotificationID:  domain.NotificationID,
		ConfigurationID: domain.ConfigurationID,
		Provider:        domain.Provider,
		Status:          string(domain.Status),
		AttemptCount:    domain.AttemptCount,
		LastError:       domain.LastError,
		LastAttemptAt:   domain.LastAttemptAt,
		SucceededAt:     domain.SucceededAt,
	}
}

func relatedEntityToDomain(model *relatedEntityModel) notification.StoredRelatedEntity {
	name := ""
	if model.Name != nil {
		name = *model.Name
	}

	return notification.StoredRelatedEntity{
		NotificationID: model.NotificationID,
		Type:           model.EntityType,
		ID:             model.EntityID,
		Name:           name,
	}
}

func relatedEntityToModel(domain notification.StoredRelatedEntity) relatedEntityModel {
	var name *string
	if domain.Name != "" {
		name = new(domain.Name)
	}

	return relatedEntityModel{
		NotificationID: domain.NotificationID,
		EntityType:     domain.Type,
		EntityID:       domain.ID,
		Name:           name,
	}
}
