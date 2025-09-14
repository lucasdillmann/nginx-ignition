package integration

import (
	"context"
	"sort"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type service struct {
	repository       Repository
	adaptersResolver func() ([]Adapter, error)
}

var defaultSettings = &Integration{
	ID:         "",
	Enabled:    false,
	Parameters: make(map[string]any),
}

func newService(repository Repository, adaptersResolver func() ([]Adapter, error)) *service {
	return &service{
		repository:       repository,
		adaptersResolver: adaptersResolver,
	}
}

func (s *service) list(ctx context.Context) ([]*ListOutput, error) {
	adapters, err := s.adaptersResolver()
	if err != nil {
		return nil, err
	}

	sort.Slice(adapters, func(i, j int) bool {
		return (adapters)[i].Priority() < (adapters)[j].Priority()
	})

	var outputs []*ListOutput
	for _, adapter := range adapters {
		settings, err := s.repository.FindByID(ctx, adapter.ID())
		if err != nil {
			return nil, err
		}

		if settings == nil {
			settings = defaultSettings
		}

		outputs = append(outputs, &ListOutput{
			ID:          adapter.ID(),
			Name:        adapter.Name(),
			Description: adapter.Description(),
			Enabled:     settings.Enabled,
		})
	}

	return outputs, nil
}

func (s *service) getById(ctx context.Context, id string) (*GetByIdOutput, error) {
	adapter := s.findAdapter(id)
	if adapter == nil {
		return nil, nil
	}

	settings, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		settings = defaultSettings
	}

	dynamicFields := adapter.ConfigurationFields()
	dynamic_fields.RemoveSensitiveFields(&settings.Parameters, dynamicFields)

	return &GetByIdOutput{
		ID:                  id,
		Name:                adapter.Name(),
		Description:         adapter.Description(),
		Enabled:             settings.Enabled,
		Parameters:          settings.Parameters,
		ConfigurationFields: dynamicFields,
	}, nil
}

func (s *service) listOptions(
	ctx context.Context,
	integrationId string,
	pageNumber, pageSize int,
	searchTerms *string,
	tcpOnly bool,
) (*pagination.Page[*AdapterOption], error) {
	adapter := s.findAdapter(integrationId)
	if adapter == nil {
		return nil, nil
	}

	settings, err := s.findSettings(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if !settings.Enabled {
		return nil, integrationDisabledError()
	}

	options, err := adapter.GetAvailableOptions(ctx, settings.Parameters, pageNumber, pageSize, searchTerms, tcpOnly)
	if err != nil {
		return nil, err
	}

	sort.Slice(options.Contents, func(i, j int) bool {
		return options.Contents[i].Name < options.Contents[j].Name
	})

	return options, nil
}

func (s *service) getOptionById(ctx context.Context, integrationId, optionId string) (*AdapterOption, error) {
	adapter := s.findAdapter(integrationId)
	if adapter == nil {
		return nil, nil
	}

	settings, err := s.findSettings(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if !settings.Enabled {
		return nil, integrationDisabledError()
	}

	return adapter.GetAvailableOptionById(ctx, settings.Parameters, optionId)
}

func (s *service) configureById(ctx context.Context, id string, enabled bool, parameters map[string]any) error {
	adapter := s.findAdapter(id)
	if adapter == nil {
		return integrationNotFoundError()
	}

	if enabled {
		if err := dynamic_fields.Validate(adapter.ConfigurationFields(), parameters); err != nil {
			return err
		}
	}

	configuration := &Integration{
		ID:         id,
		Enabled:    enabled,
		Parameters: parameters,
	}

	return s.repository.Save(ctx, configuration)
}

func (s *service) getOptionUrl(ctx context.Context, integrationId, optionId string) (*string, error) {
	adapter := s.findAdapter(integrationId)
	if adapter == nil {
		return nil, integrationNotFoundError()
	}

	settings, err := s.findSettings(ctx, integrationId)
	if err != nil {
		return nil, err
	}

	if !settings.Enabled {
		return nil, integrationDisabledError()
	}

	url, err := adapter.GetOptionProxyUrl(ctx, settings.Parameters, optionId)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *service) findAdapter(id string) Adapter {
	adapters, err := s.adaptersResolver()
	if err != nil {
		return nil
	}

	for _, adapter := range adapters {
		if adapter.ID() == id {
			return adapter
		}
	}

	return nil
}

func (s *service) findSettings(ctx context.Context, id string) (*Integration, error) {
	settings, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, integrationNotConfiguredError()
	}

	return settings, nil
}

func integrationDisabledError() error {
	return core_error.New("Integration is disabled", true)
}

func integrationNotConfiguredError() error {
	return core_error.New("Integration is not configured", true)
}

func integrationNotFoundError() error {
	return core_error.New("Integration not found with provided ID", true)
}
