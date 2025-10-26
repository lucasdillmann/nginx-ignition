package integration

import (
	"encoding/json"

	"dillmann.com.br/nginx-ignition/core/integration"
)

func toDomain(model *integrationModel) (*integration.Integration, error) {
	parameters := make(map[string]any)
	err := json.Unmarshal([]byte(model.Parameters), &parameters)
	if err != nil {
		return nil, err
	}

	return &integration.Integration{
		ID:         model.ID,
		Driver:     model.Driver,
		Name:       model.Name,
		Enabled:    model.Enabled,
		Parameters: parameters,
	}, nil
}

func toModel(domain *integration.Integration) (*integrationModel, error) {
	parameters, err := json.Marshal(domain.Parameters)
	if err != nil {
		return nil, err
	}

	return &integrationModel{
		ID:         domain.ID,
		Driver:     domain.Driver,
		Name:       domain.Name,
		Enabled:    domain.Enabled,
		Parameters: string(parameters),
	}, nil
}
