package certificate

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
	"dillmann.com.br/nginx-ignition/core/certificate"
)

func toAvailableProviderResponse(
	ctx context.Context,
	input []certificate.AvailableProvider,
) []availableProviderResponse {
	responses := make([]availableProviderResponse, 0, len(input))
	for _, provider := range input {
		responses = append(responses, availableProviderResponse{
			ID:            provider.ID(),
			Name:          provider.Name(ctx),
			Priority:      provider.Priority(),
			DynamicFields: dynamicfield.ToResponse(provider.DynamicFields(ctx)),
		})
	}

	return responses
}

func toIssueCertificateResponse(
	cert *certificate.Certificate,
	err error,
) *issueCertificateResponse {
	var errorReason *string
	if err != nil {
		errorStr := err.Error()
		errorReason = &errorStr
	}

	var certificateID *uuid.UUID
	if cert != nil {
		certificateID = &cert.ID
	}

	return &issueCertificateResponse{
		Success:       err == nil,
		ErrorReason:   errorReason,
		CertificateID: certificateID,
	}
}

func toRenewCertificateResponse(err error) *renewCertificateResponse {
	var errorReason *string
	if err != nil {
		errorStr := err.Error()
		errorReason = &errorStr
	}

	return &renewCertificateResponse{
		Success:     err == nil,
		ErrorReason: errorReason,
	}
}

func toCertificateResponse(input *certificate.Certificate) *certificateResponse {
	return &certificateResponse{
		ID:          input.ID,
		DomainNames: input.DomainNames,
		ProviderID:  input.ProviderID,
		IssuedAt:    input.IssuedAt,
		ValidUntil:  input.ValidUntil,
		ValidFrom:   input.ValidFrom,
		RenewAfter:  input.RenewAfter,
		Parameters:  input.Parameters,
	}
}

func toIssueCertificateRequest(input *issueCertificateRequest) *certificate.IssueRequest {
	return &certificate.IssueRequest{
		ProviderID:  input.ProviderID,
		DomainNames: input.DomainNames,
		Parameters:  input.Parameters,
	}
}
