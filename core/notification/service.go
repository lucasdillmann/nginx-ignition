package notification

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/user"
)

type service struct {
	repository   Repository
	userCommands user.Commands
	i18nCommands i18n.Commands
	providers    func() []Provider
}

func newService(
	repository Repository,
	userCommands user.Commands,
	i18nCommands i18n.Commands,
	providers func() []Provider,
) *service {
	return &service{
		repository:   repository,
		userCommands: userCommands,
		i18nCommands: i18nCommands,
		providers:    providers,
	}
}

func (s *service) findProvider(providerID string) Provider {
	for _, provider := range s.providers() {
		if provider.ID() == providerID {
			return provider
		}
	}
	return nil
}

func toRelatedEntities(values []StoredRelatedEntity) []RelatedEntity {
	result := make([]RelatedEntity, len(values))
	for index, value := range values {
		result[index] = RelatedEntity{
			Type: value.Type,
			ID:   value.ID,
			Name: value.Name,
		}
	}
	return result
}
