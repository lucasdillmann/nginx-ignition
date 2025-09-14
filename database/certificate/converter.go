package certificate

import (
	"encoding/json"

	"dillmann.com.br/nginx-ignition/core/certificate"
)

func toDomain(model *certificateModel) (*certificate.Certificate, error) {
	params, err := parseParameters(model.Parameters)
	if err != nil {
		return nil, err
	}

	return &certificate.Certificate{
		ID:                 model.ID,
		DomainNames:        model.DomainNames,
		ProviderID:         model.ProviderID,
		IssuedAt:           model.IssuedAt,
		ValidUntil:         model.ValidUntil,
		ValidFrom:          model.ValidFrom,
		RenewAfter:         model.RenewAfter,
		PrivateKey:         model.PrivateKey,
		PublicKey:          model.PublicKey,
		CertificationChain: model.CertificationChain,
		Parameters:         params,
		Metadata:           model.Metadata,
	}, nil
}

func toModel(domain *certificate.Certificate) (*certificateModel, error) {
	params, err := formatParameters(domain.Parameters)
	if err != nil {
		return nil, err
	}

	return &certificateModel{
		ID:                 domain.ID,
		DomainNames:        domain.DomainNames,
		ProviderID:         domain.ProviderID,
		IssuedAt:           domain.IssuedAt,
		ValidUntil:         domain.ValidUntil,
		ValidFrom:          domain.ValidFrom,
		RenewAfter:         domain.RenewAfter,
		PrivateKey:         domain.PrivateKey,
		PublicKey:          domain.PublicKey,
		CertificationChain: domain.CertificationChain,
		Parameters:         params,
		Metadata:           domain.Metadata,
	}, nil
}

func parseParameters(params string) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal([]byte(params), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func formatParameters(params map[string]any) (string, error) {
	result, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
