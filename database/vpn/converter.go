package vpn

import (
	"encoding/json"

	"dillmann.com.br/nginx-ignition/core/vpn"
)

func toDomain(model *vpnModel) (*vpn.VPN, error) {
	parameters := make(map[string]any)
	err := json.Unmarshal([]byte(model.Parameters), &parameters)
	if err != nil {
		return nil, err
	}

	return &vpn.VPN{
		ID:         model.ID,
		Driver:     model.Driver,
		Name:       model.Name,
		Enabled:    model.Enabled,
		Parameters: parameters,
	}, nil
}

func toModel(domain *vpn.VPN) (*vpnModel, error) {
	parameters, err := json.Marshal(domain.Parameters)
	if err != nil {
		return nil, err
	}

	return &vpnModel{
		ID:         domain.ID,
		Driver:     domain.Driver,
		Name:       domain.Name,
		Enabled:    domain.Enabled,
		Parameters: string(parameters),
	}, nil
}
