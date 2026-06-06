package notification

import (
	"encoding/json"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
	"dillmann.com.br/nginx-ignition/core/notification"
)

func toNotificationDTO(data *notification.Notification) notificationResponse {
	return notificationResponse{
		ID:              data.ID,
		Title:           data.Title,
		Summary:         data.Summary,
		Category:        string(data.Category),
		Read:            data.ReadAt != nil,
		CreatedAt:       data.CreatedAt,
		OccurredAt:      data.Payload.OccurredAt,
		Tags:            data.Payload.Tags,
		Sections:        toContentSectionDTOs(data.Payload.Sections),
		Actions:         toActionDTOs(data.Payload.Actions),
		RelatedEntities: toRelatedEntityDTOs(data.RelatedEntities),
	}
}

func toContentSectionDTOs(
	sections []notification.DeliverableContentSection,
) []contentSectionResponse {
	result := make([]contentSectionResponse, len(sections))
	for index, section := range sections {
		result[index] = contentSectionResponse{
			Title: section.Title,
			Body:  section.Body,
		}
	}

	return result
}

func toActionDTOs(actions []notification.DeliverableAction) []actionResponse {
	result := make([]actionResponse, len(actions))
	for index, action := range actions {
		result[index] = actionResponse{
			Label: action.Label,
			URL:   action.URL,
		}
	}

	return result
}

func toRelatedEntityDTOs(entities []notification.RelatedEntity) []relatedEntityResponse {
	result := make([]relatedEntityResponse, len(entities))
	for index, entity := range entities {
		result[index] = relatedEntityResponse{
			Type: entity.Type,
			ID:   entity.ID,
			Name: entity.Name,
		}
	}

	return result
}

func toCategoryDTO(data notification.CategoryInfo) categoryResponse {
	return categoryResponse{
		ID:   string(data.ID),
		Name: data.Name,
	}
}

func categoriesToJSON(categories *[]notification.Category) (json.RawMessage, error) {
	if categories == nil {
		return json.RawMessage("null"), nil
	}

	values := make([]string, len(*categories))
	for index, category := range *categories {
		values[index] = string(category)
	}

	encoded, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func toConfigurationDTO(data *notification.Configuration) (configurationResponse, error) {
	categories, err := categoriesToJSON(data.Categories)
	if err != nil {
		return configurationResponse{}, err
	}

	return configurationResponse{
		ID:         data.ID,
		Name:       data.Name,
		Provider:   data.Provider,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
		Categories: categories,
	}, nil
}

func parseCategories(raw json.RawMessage) (*[]notification.Category, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	if string(raw) == "null" {
		return nil, nil
	}

	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}

	categories := make([]notification.Category, len(values))
	for index, value := range values {
		categories[index] = notification.Category(value)
	}

	return new(categories), nil
}

func toConfigurationDomain(
	userID, id uuid.UUID,
	request *configurationRequest,
) (*notification.Configuration, error) {
	categories, err := parseCategories(request.Categories)
	if err != nil {
		return nil, err
	}

	return &notification.Configuration{
		ID:         id,
		UserID:     userID,
		Name:       request.Name,
		Provider:   request.Provider,
		Enabled:    request.Enabled,
		Parameters: request.Parameters,
		Categories: categories,
	}, nil
}

func toAvailableProviderDTO(data notification.AvailableProvider) availableProviderResponse {
	return availableProviderResponse{
		ID:                    data.ID,
		Name:                  data.Name,
		ImportantInstructions: data.ImportantInstructions,
		ConfigurationFields:   dynamicfield.ToResponse(data.ConfigurationFields),
	}
}
