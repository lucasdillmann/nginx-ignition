package notification

import (
	"context"
	"sort"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

func (s *service) ListConfigurations(
	ctx context.Context,
	userID uuid.UUID,
) ([]Configuration, error) {
	configurations, err := s.repository.FindConfigurationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for index := range configurations {
		configuration := &configurations[index]
		provider := s.findProvider(configuration.Provider)
		if provider != nil {
			dynamicfields.RemoveSensitiveFields(
				&configuration.Parameters,
				provider.ConfigurationFields(ctx),
			)
		}
	}

	return configurations, nil
}

func (s *service) GetConfiguration(
	ctx context.Context,
	userID uuid.UUID,
	id uuid.UUID,
) (*Configuration, error) {
	configuration, err := s.repository.FindConfigurationByIDAndUserID(ctx, id, userID)
	if err != nil || configuration == nil {
		return nil, err
	}

	provider := s.findProvider(configuration.Provider)
	if provider != nil {
		dynamicfields.RemoveSensitiveFields(
			&configuration.Parameters,
			provider.ConfigurationFields(ctx),
		)
	}

	return configuration, nil
}

func (s *service) SaveConfiguration(
	ctx context.Context,
	userID uuid.UUID,
	configuration *Configuration,
) error {
	existing, err := s.repository.FindConfigurationByIDAndUserID(
		ctx,
		configuration.ID,
		userID,
	)
	if err != nil {
		return err
	}

	provider := s.findProvider(configuration.Provider)
	if existing != nil {
		configuration.UserID = existing.UserID
		if provider != nil {
			configuration.Parameters = dynamicfields.MergeSensitiveFields(
				configuration.Parameters,
				existing.Parameters,
				provider.ConfigurationFields(ctx),
			)
		}
	} else {
		configuration.UserID = userID
	}

	if err := newValidator(
		s.repository,
		provider,
	).validate(ctx, userID, configuration); err != nil {
		return err
	}

	return s.repository.SaveConfiguration(ctx, configuration)
}

func (s *service) DeleteConfiguration(ctx context.Context, userID, id uuid.UUID) error {
	return s.repository.DeleteConfigurationByIDAndUserID(ctx, id, userID)
}

func (s *service) GetAvailableProviders(ctx context.Context) ([]AvailableProvider, error) {
	registeredProviders := s.providers()
	sort.Slice(registeredProviders, func(left, right int) bool {
		return registeredProviders[left].Name(ctx).String() <
			registeredProviders[right].Name(ctx).String()
	})

	output := make([]AvailableProvider, len(registeredProviders))
	for index, provider := range registeredProviders {
		output[index] = AvailableProvider{
			ID:                    provider.ID(),
			Name:                  provider.Name(ctx),
			ImportantInstructions: provider.ImportantInstructions(ctx),
			ConfigurationFields:   provider.ConfigurationFields(ctx),
		}
	}

	return output, nil
}

func (s *service) ListCategories(ctx context.Context) ([]CategoryInfo, error) {
	categories := AllCategories()
	output := make([]CategoryInfo, len(categories))

	for index, category := range categories {
		output[index] = CategoryInfo{
			ID:   category,
			Name: CategoryName(ctx, category),
		}
	}

	return output, nil
}
